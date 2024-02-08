package wire

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/server"
	"github.com/bhmy-shm/gofks/example/rpc/cascade-rpc/internal/svc"
)

type CascadeWire struct {
	Conf          *gofkConf.Config      `inject:"-"`
	CascadeCtx    *svc.CascadeContext   `inject:"-"`
	CascadeServer *server.CascadeServer `inject:"-"`
}

func CascadeServerWire(c *gofkConf.Config) *CascadeWire {
	return &CascadeWire{
		Conf: c,
	}
}

func (this *CascadeWire) WireCascadeServer() {
	this.CascadeCtx = svc.NewCascadeContext(this.Conf)
	this.CascadeServer = server.NewCascadeServer(this.CascadeCtx)
}
