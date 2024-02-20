package main

import (
	"fmt"
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/wire"
	"github.com/bhmy-shm/gofks/zrpc/interceptros"
	serverInterceptors "github.com/bhmy-shm/gofks/zrpc/interceptros/server"
)

func main() {

	conf := gofkConf.Load()

	serviceWire := wire.NewServiceWire(conf)

	gofks.NewRpcServer(conf).
		WireApply(serviceWire).
		Register(serviceWire.Server).
		Attach(
			serverInterceptors.WithUnaryServerInterceptors(interceptros.LoggerInterceptor),
		).
		Start()

	fmt.Printf("Starting rpc server at %s...\n", conf.GetServer().Listener())
}
