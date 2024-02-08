package main

import (
	"fmt"
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/wire"
	"github.com/bhmy-shm/gofks/zrpc/interceptros"
	serverInterceptors "github.com/bhmy-shm/gofks/zrpc/interceptros/server"
)

func main() {

	conf := &gofkConf.Config{}

	cascadeWire := wire.CascadeServerWire(conf)

	gofks.NewRpcServer(conf).
		WireApply(cascadeWire).
		Register(
			cascadeWire.CascadeServer,
		).
		Attach(
			serverInterceptors.WithUnaryServerInterceptors(interceptros.LoggerInterceptor),
		).
		Start()

	fmt.Printf("Starting rpc server at %s...\n", conf.GetServer().Listener)
}
