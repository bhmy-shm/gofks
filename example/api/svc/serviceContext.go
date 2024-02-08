package svc

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/client"
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
