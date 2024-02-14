package controls

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/example/websocket/wire"
)

type WsCase struct {
	*wire.WsWire `inject:"-"`
}

func WsController() *WsCase {
	return &WsCase{}
}

func (this *WsCase) Build(gofk *gofks.Gofk) {
	ws := gofk.Group("/ws")
	ws.GET("testHome", this.wsHome)
	ws.Handle("GET", "server", this.wsServer)
}

func (this *WsCase) Name() string {
	return "userCase"
}

func (this *WsCase) Wire() *wire.WsWire {
	return this.WsWire
}
