package gofk

import (
	"fmt"
	"github.com/bhmy-shm/gofks/Injector"
	expr2 "github.com/bhmy-shm/gofks/expr"
	"github.com/bhmy-shm/gofks/middle"
	"github.com/bhmy-shm/gofks/pkg"
	"github.com/bhmy-shm/gofks/pkg/thread"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/pprof"
)

type Gofk struct {
	engine   *gin.Engine
	group    *gin.RouterGroup
	exprData map[string]interface{}
	file     *pkg.File //配置文件
}

func Ignite() *Gofk {
	g := &Gofk{engine: gin.New(),
		exprData: map[string]interface{}{},
	}

	g.engine.Use(middle.ErrorHandler())
	return g
}

func (g *Gofk) Root(path string) *Gofk {
	g.group = g.engine.Group(path)
	return g
}

func (g *Gofk) Group(groupPath string) *gin.RouterGroup {

	if g.group != nil {
		return g.group.Group(groupPath)
	}

	//每次需要新增加一个路由
	return g.engine.Group(groupPath)
}

func (g *Gofk) Attach(handlerFunc ...gin.HandlerFunc) *Gofk {
	if g.group == nil {
		g.group = g.engine.Group("")
	}
	g.group.Use(handlerFunc...)
	return g
}

func (g *Gofk) Beans(beans ...Bean) *Gofk {
	//取出Bean的名称，然后加入到 exprData 里面
	for _, bean := range beans {
		g.exprData[bean.Injector()] = bean
		injector.BeanFactory.Set(bean)
	}
	return g
}

func (g *Gofk) Mount(classes ...IClass) *Gofk {
	for _, class := range classes {
		class.Build(g) //挂载路由
		g.Beans(class) //处理路由中的任务（表达式任务）
	}
	return g
}

func (g *Gofk) Config(beans ...interface{}) *Gofk {
	injector.BeanFactory.Config(beans...)
	return g
}

func (g *Gofk) Watcher() *Gofk {

	f, err := pkg.LoadFile()
	pkg.Error(err, "读取监听配置文件失败")

	g.file = f
	g.file.YamlMerge() //yaml的方式加载conf 到f对象，以及内存当中

	//协程监听更新config文件
	go pkg.ReadWatcher(g.file)
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

func (g *Gofk) applyAll() {
	for t, v := range injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			injector.BeanFactory.Apply(v.Interface())
		}
	}
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
			port, err = pkg.GetPath("Server", "port").Int()
			pkg.Error(err)
		}
	}

	//启动定时任务
	thread.GetCronTask().Start()

	//加载依赖注入
	g.applyAll()

	//启动服务
	g.engine.Run(fmt.Sprintf(":%d", port))
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
