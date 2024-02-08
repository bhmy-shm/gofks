package gofks

import (
	"fmt"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/injector"
	"github.com/bhmy-shm/gofks/middle"
	"github.com/bhmy-shm/gofks/zrpc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/pprof"
)

type (
	GofkCase interface {
		Build(*Gofk)
		Name() string
	}

	Gofk struct {
		engine   *gin.Engine
		group    *gin.RouterGroup
		exprData *injector.BeanCache
		conf     *gofkConf.Config //配置文件

		cron      CronInter       //计划任务
		rpcServer *zrpc.RpcServer //rpc服务端
	}
)

// Ignite 加载web路由
func Ignite(rootPath string, middles ...gin.HandlerFunc) *Gofk {
	g := &Gofk{
		engine:   gin.New(),
		exprData: injector.NewBeanCache(),
		cron:     NewCronTask(),
	}

	//
	g.group = g.engine.Group(rootPath)

	middles = append(middles, middle.ErrorHandler())
	g.engine.Use(middles...)

	return g
}

// Group 设置web子路由
func (g *Gofk) Group(path string) *gin.RouterGroup {
	if g.group == nil {
		g.group = g.engine.Group(path)
		return g.group
	}
	return g.group.Group(path)
}

// GetGroup 获取当前路由
func (g *Gofk) GetGroup() *gin.RouterGroup {
	return g.group
}

// Mount 挂载 web、plugin 服务插件
func (g *Gofk) Mount(classes ...GofkCase) *Gofk {
	for _, class := range classes {
		class.Build(g) //挂载路由
		g.beans(class) //处理路由中的任务（表达式任务）
	}
	return g
}

// Launch 启动 web 服务
func (g *Gofk) Launch() {

	var (
		address string
		err     error
	)

	if g.conf != nil {
		address = g.conf.GetServer().Listener()
	} else {
		address = ":8085"
	}

	//启动定时任务
	g.cron.Get().Start()

	g.applyAll()

	//启动服务
	if err = g.engine.Run(address); err != nil {
		errorx.Fatal(err, "gofk Lauch Running failed")
	}
}

// WireApply 配置依赖，注入注册
func (g *Gofk) WireApply(beans ...interface{}) *Gofk {
	injector.BeanFactory.Config(beans...)
	g.applyAll()
	return g
}

// 加载依赖注入
func (g *Gofk) applyAll() {

	injector.BeanFactory.GetBeanMapper().Range(func(key, value any) bool {

		t := key.(reflect.Type)
		v := value.(reflect.Value)

		if t.Elem().Kind() == reflect.Struct {
			injector.BeanFactory.Apply(v.Interface())
		}
		return true
	})
}

// 服务缓存注册
func (g *Gofk) beans(beans ...GofkCase) *Gofk {

	for _, bean := range beans {

		g.exprData.Add(bean.Name(), bean)

		injector.BeanFactory.Set(bean)
	}
	return g
}

// LoadWatch 加载配置文件与监听
func (g *Gofk) LoadWatch(conf *gofkConf.Config) *Gofk {

	f, err := gofkConf.LoadFile()
	if err != nil {
		log.Fatalln("check the application.yaml file exists:", err)
	}

	//yaml的方式加载conf 到f对象，以及内存当中
	if ok := f.YamlMerge(conf); ok {
		g.conf = conf
	}

	//协程监听更新config文件
	go gofkConf.ReadWatcher(f, conf)

	return g
}

// WebSocket 简单封装
func WebSocket(conf *gofkConf.Config, root string) *Gofk {
	return Ignite(root).LoadWatch(conf)
}

// Cron 计划任务
func (g *Gofk) Cron(cron string, expr interface{}) *Gofk {
	var err error

	switch expr.(type) {
	case func():
		f := expr.(func())
		_, err = g.cron.Get().AddFunc(cron, f)
	//case expr2.Expr:
	//	exp := expr.(expr2.Expr) //这里的 exp 就是传入的 表达式
	//	_, err = g.cron.Get().AddFunc(cron, func() {
	//		_, expErr := expr2.ExecExpr(exp, g.exprData) //处理表达式
	//		if expErr != nil {
	//			log.Println(expErr)
	//		}
	//	})
	default:
		log.Fatalln("计划任务 Func 配置有误")
	}

	if err != nil {
		log.Println(err)
	}
	return g
}

// PProf 程序火箭图
func (g *Gofk) PProf(path string) {

	go func() {
		http.ListenAndServe(path, nil)
	}()

	go func() {
		f, err := os.Create("./cpuprofile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}()
}
