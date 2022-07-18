package gofk

import (
	"fmt"
	"github.com/bhmy-shm/gofks/Injector"
	expr2 "github.com/bhmy-shm/gofks/expr"
	"github.com/bhmy-shm/gofks/pkg/config"
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"github.com/bhmy-shm/gofks/pkg/thread"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
)

const (
	ConfigPath = "./application.yaml"
)

var Empty = &struct{}{}
var innerRouter *GoftTree // inner tree node . backup httpmethod and path
var innerRouter_once sync.Once

func getInnerRouter() *GoftTree {
	innerRouter_once.Do(func() {
		innerRouter = NewGoftTree()
	})
	return innerRouter
}

type Gofk struct {
	engine   *gin.Engine
	group    *gin.RouterGroup       //路由分组
	file     *config.File           //配置文件
	exprData map[string]interface{} //表达式
}

func Ignite() *Gofk {
	g := &Gofk{engine: gin.New(),
		exprData: map[string]interface{}{},
	}

	g.engine.Use(errorx.ErrorHandler())
	return g
}

func (g *Gofk) Watcher() *Gofk {

	f, err := config.LoadFile()
	errorx.Error(err, "读取监听配置文件失败")

	g.file = f
	g.file.YamlMerge() //yaml的方式加载conf 到f对象，以及内存当中

	//协程监听更新config文件
	go config.ReadWatcher(g.file)
	return g
}

func (g *Gofk) Launch() {
	var (
		port = 8080
		err  error
	)
	//判断是否存在配置文件并转换成map进行记录
	if g.file != nil {
		if g.file.GetConf() != nil {
			//如果已经存在配置文件记录，则通过配置文件拿到port端口号并启动服务
			port, err = config.GetPath("Server", "port").Int()
			errorx.Error(err)
		}
	}

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

func (g *Gofk) Attach(fs ...Fairing) *Gofk {
	for _, f := range fs {
		Injector.BeanFactory.Set(f)
	}

	onceFairingHandler().AddFairing(fs...)
	return g
}

//处理依赖注入

func (g *Gofk) Beans(beans ...Bean) *Gofk {
	//取出Bean的名称，然后加入到 exprData 里面
	for _, bean := range beans {
		g.exprData[bean.Name()] = bean
		Injector.BeanFactory.Set(bean)
	}
	return g
}

func (g *Gofk) Mount(group string, classes ...IClass) *Gofk {
	g.group = g.engine.Group(group)
	for _, class := range classes {
		class.Build(g) //挂载路由
		g.Beans(class) //处理路由中的任务（表达式任务）
	}
	return g
}

func (g *Gofk) Cron(cron string, expr interface{}) *Gofk {
	var err error

	switch expr.(type) {
	case func():
		f := expr.(func())
		_, err = thread.GetCronTask().AddFunc(cron, f)
	case expr2.Expr:
		exp := expr.(expr2.Expr) //这里的 exp 就是传入的 表达式
		_, err = thread.GetCronTask().AddFunc(cron, func() {
			_, expErr := expr2.ExecExpr(exp, g.exprData) //处理表达式
			if expErr != nil {
				log.Println(expErr)
			}
		})
	default:
		log.Fatalln("计划任务 Func 配置有误")
	}

	if err != nil {
		log.Println(err)
	}
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
