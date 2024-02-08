package gofks

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/zrpc"
	"google.golang.org/grpc"
)

func NewRpcServer(conf *gofkConf.Config) *Gofk {

	gofks := &Gofk{}

	gofks.LoadWatch(conf)

	gofks.rpcServer = zrpc.NewRpcServer(conf)

	return gofks
}

func (gk *Gofk) Register(servers ...zrpc.RpcRegisterInter) *Gofk {

	if gk.rpcServer == nil {
		errorx.Fatal(errorx.ErrCodeParamsErr)
		return nil
	}

	gk.rpcServer.Register(servers...)

	return gk
}

func (gk *Gofk) Attach(opts ...grpc.ServerOption) *Gofk {
	if gk.rpcServer == nil {
		errorx.Fatal(errorx.ErrCodeParamsErr)
		return nil
	}

	gk.rpcServer.AddGrpcOptions(opts...)
	return gk
}

func (gk *Gofk) Start() *Gofk {
	if gk.rpcServer == nil {
		errorx.Fatal(errorx.ErrRpcServerNotFound)
		return nil
	}

	gk.applyAll()

	gk.rpcServer.Start()
	return gk
}

func (gk *Gofk) Stop() {
	if gk.rpcServer == nil {
		errorx.Fatal(errorx.ErrRpcServerNotFound)
	}
	gk.rpcServer.Stop()
}
