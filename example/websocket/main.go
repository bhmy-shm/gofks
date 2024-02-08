package main

import (
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/websocket/controls"
	"github.com/bhmy-shm/gofks/example/websocket/wire"
)

func main() {

	conf := gofkConf.New()

	gofks.WebSocket(conf, "/v1").
		WireApply(
			wire.NewWsWire(conf),
		).
		Mount(
			controls.WsController(),
		).
		Launch()
}
