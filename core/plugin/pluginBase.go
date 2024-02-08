package plugin

import (
	"context"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/register"
	"github.com/bhmy-shm/gofks/core/utils/snowflake"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	PluginBase interface {
		ID() string
		Name() string
		Version() string
		Metadata() map[string]string

		Add(items ...PluginItem)
		Run()
		Stop() error
		GetOpts() *pluginOptions
	}

	pluginAccess struct {
		ctx    context.Context
		cancel func()

		opts    *pluginOptions
		servers []PluginItem //注册插件
		lock    sync.Mutex
		conf    *gofkConfs.PluginConfig
	}
)

func defaultPluginAccess() *pluginAccess {

	//从配置文件获取插件ID
	nodeId := gofkConf.GetPath[int]("plugin", "nodeId").Value().UnwrapOr(1)

	return &pluginAccess{
		opts: &pluginOptions{
			id:               snowflake.SnowflakeUUid(int64(nodeId)),
			metadata:         make(map[string]string),
			ctx:              context.Background(),
			sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
			monitorEnable:    false,
			registerEnable:   false,
			registrarTimeout: 10 * time.Second,
			stopTimeout:      10 * time.Second,
		},
	}
}

func NewBase(conf *gofkConfs.PluginConfig, opts ...OptionFunc) PluginBase {

	access := defaultPluginAccess()
	access.conf = conf

	for _, fn := range opts {
		fn(access.opts)
	}

	access.ctx, access.cancel = context.WithCancel(access.opts.ctx)
	return access
}

// ID returns app instance id.
func (p *pluginAccess) ID() string { return p.opts.id }

// Name returns service name.
func (p *pluginAccess) Name() string { return p.opts.name }

// Version returns app version.
func (p *pluginAccess) Version() string { return p.opts.version }

// Metadata returns service metadata.
func (p *pluginAccess) Metadata() map[string]string { return p.opts.metadata }

func (p *pluginAccess) Add(items ...PluginItem) {
	p.servers = items
}

func (p *pluginAccess) Run() {

	var err error

	//设置唯一的服务基础信息
	p.lock.Lock()
	p.opts.instance = p.buildInstance()
	p.lock.Unlock()

	//生成该服务的上下文信息
	pCtx := WithContext(p.ctx, p)
	heartbeatCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	//是否需要开启插件(事件驱动)监控
	if p.opts.monitorEnable {
		p.monitorStart(heartbeatCtx, p.Name())
	}

	//运行服务
	wg := sync.WaitGroup{}
	for _, server := range p.servers {
		serverCopy := server
		wg.Add(1)
		go func(item PluginItem) {
			defer wg.Done()
			if err = item.Start(); err != nil {
				logx.Error("item server Start failed:", err)
				return
			}
		}(serverCopy)
	}

	//是否需要开启服务注册
	if p.opts.registerEnable {
		if err = p.registerStart(p.opts.instance); err != nil {
			logx.Error("开启服务时注册etcd 出现异常：", err)
		}
	}

	//监听信号退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, p.opts.sigs...)
	go func() {
		select {
		case <-pCtx.Done():
			defer os.Exit(1)
			if err = p.Stop(); err != nil {
				logx.Error("监听到程序超时退出失败")
			} else {
				logx.Info("监听到程序超时退出成功")
			}
		case <-c:
			defer os.Exit(0)
			if err = p.Stop(); err != nil {
				logx.Error("监听到信号退出程序失败")
			} else {
				logx.Info("监听到信号退出程序成功")
			}
		}
	}()

	//阻塞等待协程执行结束
	wg.Wait()
}

func (p *pluginAccess) Stop() error {
	p.lock.Lock()
	instance := p.opts.instance
	p.lock.Unlock()

	if p.opts.registerEnable {
		if err := p.registerEnd(instance); err != nil {
			logx.Error("结束程序时 关闭etcd 出现异常：", err)
		}
	}

	if p.cancel != nil {
		p.cancel()
	}

	for _, server := range p.servers {
		logx.Info("准备结束程序...", server.Name())
		server.Exit()
	}

	log.Println(p.Name(), p.Version(), "服务结束成功!")
	return nil
}

func (p *pluginAccess) GetOpts() *pluginOptions {
	return p.opts
}

// --------- 内部方法 -------

func (p *pluginAccess) buildInstance() *register.ServiceInstance {
	return &register.ServiceInstance{
		ID:       p.opts.id,
		Name:     p.opts.name,
		Version:  p.opts.version,
		Metadata: p.opts.metadata,
	}
}

func (p *pluginAccess) monitorStart(heartbeatCtx context.Context, systemName string) {

	t := timex.NewTicker(time.Duration(p.conf.MonitorInterval()) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.Chan():
			//TODO 消息推送，mq.PubHeartbeatMsg(systemName)， systemName是优化点
		case <-heartbeatCtx.Done():
			return
		}
	}
}

func (p *pluginAccess) registerStart(instance *register.ServiceInstance) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.opts.registrarTimeout)
	defer cancel()

	return p.opts.registrar.Register(ctx, instance)
}

func (p *pluginAccess) registerEnd(instance *register.ServiceInstance) error {

	ctx, cancel := context.WithTimeout(WithContext(p.ctx, p), p.opts.registrarTimeout)
	defer cancel()

	return p.opts.registrar.Deregister(ctx, instance)
}
