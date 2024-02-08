# Gofks Web Api

一个基于 gofks 的 api程序，主要包含3个模块：
1. controls 路由插件
   - 主要的 gin http 封装 

2. svc context上下文
   - 实现 rpc-client，rpc-websocket-router 的路由调用
   
3. wire 依赖注入 
   - 主要注入 svc context 以及其他符合写法的 inject object

配置文件固定名称：application.yaml

## Automatic started

自动生成项目，需要手动编写一个 .api文件，具体的api语法可以查看 [go-zero文档](https://go-zero.dev/docs/tasks/dsl/api)

编写 v1.api文件
```api
syntax = "v1"

info (
	title:   "gofks-web-api"
	desc:    "v1版本文档"
	version: "v1.0.0"
)

type Response {
	Code   int32       `json:"code"`
	Msg    string      `json:"msg"`
	Reason string      `json:"reason"`
	Data   interface{} `json:"data"`
}

type LoginReq {
	Account     string `json:"account"`
	Password    string `json:"password"`
	Code        string `json:"code"`
	LoginMethod string `json:"loginMethod,optional"` //目前包括 dataWeb safeWeb dataApp 三种登录方式
}

//登录
@server (
	prefix: /v1
)
service Account {
	@doc "account user login"
	@handler accountLogin
	post /login (LoginReq) returns (Response)
}

```

基于api生成项目
```shell
 gofkctl api generate v1.api --dir ./ --style=gofks
```


## Getting started

也可以手动创建一个api项目，更加灵活一些，按照example示例实现就可以。

### 项目目录结构如下：
```markdown
api
|-- README.md
|-- application.yaml
|-- controls
|   |-- account
|   |   |-- account_org.go
|   |   |-- account_user.go
|   |   `-- control.go
|   |-- middleware
|   |   `-- cros.go
|   `-- types
|       `-- account.go
|-- main.go  
|-- svc
|   `-- serviceContext.go
`-- wire
    |-- testService.go
    `-- wire.go
```

### main.go 启动文件

```go
package main

import (
	"github.com/bhmy-shm/gofks/example/api/controls/account"
	"github.com/bhmy-shm/gofks/example/api/wire"
	"github.com/bhmy-shm/gofks/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
)

func main() {

	conf := gofkConf.New() //直接根据默认的 application.yaml 生成conf对象

	gofks.Ignite("/v1"). //设置根路由
		LoadWatch(conf). //配置文件监听
		WireApply(
			wire.NewServiceWrite(conf), //依赖注入
		).
		Mount(account.UserController()). //挂载http路由
		Launch() //启动api程序
}

```

### 实现 NewServiceWrite 依赖注入

```go
package wire

import (
	"github.com/bhmy-shm/gofks/example/api/svc"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
)

type ServiceWire struct {
	Ctx *svc.ServiceContext
	T   *TestService
}

func NewServiceWrite(c *gofkConf.Config) *ServiceWire {
	return &ServiceWire{
		Ctx: svc.NewServiceContext(c),
	}
}

func (s *ServiceWire) ServiceCtx() *svc.ServiceContext {
	return s.Ctx
}

func (s *ServiceWire) Test() *TestService {
	s.T = NewTestService("svc-test")
	return s.T
}

```

```go
package svc

import (
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/client"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/zrpc"
)

type ServiceContext struct {
	AccountRpc    client.AccountClient //rpc webApi 路由
	AccountRouter client.AccountRouter //rpc websocket 路由
}

func NewServiceContext(c *gofkConf.Config) *ServiceContext {

	svc := &ServiceContext{}
	rpcClient := zrpc.NewRpcClient(c.GetRpcClient())

	if c.GetRpcClient().IsLoad() {
		svc.AccountRpc = client.NewUserClient(rpcClient)
	}

	if c.GetServer().EnableWs() {
		svc.AccountRouter = client.NewAccountRouter(rpcClient)
	}
	return svc
}

```

### 挂载 Controller 插件
```go
package account

import (
	"github.com/bhmy-shm/gofks/example/api/wire"
	"github.com/bhmy-shm/gofks/gofks"
)

type AccountCase struct {
	*wire.ServiceWire `inject:"-"`
}

func UserController() *AccountCase {
	return &AccountCase{}
}

func (s *AccountCase) Build(gofk *gofks.Gofk) {

	account := gofk.Group("account")

	user := account.Group("user")
	user.Handle("POST", "/userDetail", s.UserDetail)
	user.Handle("POST", "/userList", s.UserList)
	user.Handle("POST", "/userAdd", s.UserAdd)

	org := account.Group("org")
	org.Handle("POST", "/orgDetail", s.OrgDetail)
}

func (s *AccountCase) Name() string {
	return "userCase"
}

func (s *AccountCase) Wire() *wire.ServiceWire {
	return s.ServiceWire
}

// ===========

```