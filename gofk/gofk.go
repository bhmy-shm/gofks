package gofk

import (
	"fmt"
	"github.com/bhmy-shm/gofks/ifac"
	"github.com/bhmy-shm/gofks/pkg/config"
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"github.com/bhmy-shm/gofks/pkg/thread"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

const (
	ConfigPath = "./application.yaml"
)

type Gofk struct {
	engine      *gin.Engine
	group       *gin.RouterGroup //路由分组
	file        *config.File     //配置文件
	beanFactory *BeanFactory     //注解、依赖注入
}

func Ignite() *Gofk {
	g := &Gofk{engine: gin.New(), beanFactory: NewBeanFactory()}

	g.engine.Use(errorx.ErrorHandler())
	g.beanFactory.setBean(config.InitSysConfig())
	return g
}

func (g *Gofk) Watcher() {

	f, err := config.LoadFile()
	errorx.Error(err, "读取监听配置文件失败")

	g.file = f
	g.file.YamlMerge() //yaml的方式加载conf 到f对象，以及内存当中

	//协程监听更新config文件
	go config.ReadWatcher(g.file)
}

func (g *Gofk) Launch() {
	var (
		port int
		err  error
	)
	//判断是否存在配置文件并转换成map进行记录
	if g.file.GetConf() != nil {
		//如果已经存在配置文件记录，则通过配置文件拿到port端口号并启动服务
		port, err = config.GetPath("Server", "port").Int()
		errorx.Error(err)
	}

	//如果没有生成配置记录则走默认端口号
	port = 8080

	//启动定时任务
	thread.GetCronTask().Start()

	//Run gin
	g.engine.Run(fmt.Sprintf(":%d", port))
}

func (g *Gofk) Handle(httpMethod, relativePath string, handlers interface{}) *Gofk {
	//传入的 handlers 业务函数，会在Convert里面判断具体类型，最终返回的是 gin.HandleFunc
	if h := Convert(handlers); h != nil {
		g.engine.Handle(httpMethod, relativePath, h)
	}
	return g
}

func (g *Gofk) Attach(f ifac.Fairing) *Gofk {
	g.engine.Use(func(context *gin.Context) {
		err := f.OnRequest(context)
		if err != nil {
			context.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		} else {
			context.Next()
		}
	})
	return g
}

//处理依赖注入

func (g *Gofk) Beans(beans ...interface{}) *Gofk {
	g.beanFactory.setBean(beans...)
	return g
}

func (g *Gofk) Mount(group string, classes ...IClass) *Gofk {
	g.group = g.engine.Group(group)
	for _, class := range classes {
		class.Build(g) //挂载路由
		//g.setProp(class) //挂载db数据库
		g.beanFactory.inject(class)
	}
	return g
}

func (g *Gofk) Cron(expr string, f func()) *Gofk {
	_, err := thread.GetCronTask().AddFunc(expr, f)
	if err != nil {
		log.Println("add cron is failed", err)
	}

	//启动定时任务放在 lunch 函数中执行
	return g
}

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
