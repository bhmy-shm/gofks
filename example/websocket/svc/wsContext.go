package svc

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/wscore"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/client"
	"github.com/bhmy-shm/gofks/zrpc"
)

type WsContext struct {
	conf       *gofkConf.Config
	AccountRpc client.AccountRouter
}

func NewWsContext(conf *gofkConf.Config) *WsContext {
	return &WsContext{
		conf:       conf,
		AccountRpc: client.NewAccountRouter(zrpc.NewRpcClient(conf.GetRpcClient())),
	}
}

func (ctx *WsContext) GetRpcRouterHandler() wscore.IMsgHandler {
	return ctx.AccountRpc.Handler()
}
