package wire

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/server"
	"github.com/bhmy-shm/gofks/example/rpc/account-rpc/internal/svc"
	"github.com/bhmy-shm/gofks/zrpc"
)

type ServiceWire struct {
	ctx    *svc.ServiceContext
	Server *server.Server
}

func NewServiceWire(c *gofkConf.Config) *ServiceWire {
	rpc := &ServiceWire{
		ctx: svc.NewServiceContext(c),
	}
	rpc.Server = server.NewServer(rpc.ctx)
	return rpc
}

func (s *ServiceWire) ServiceContext() *svc.ServiceContext {
	return s.ctx
}

func (s *ServiceWire) RpcServer() zrpc.RpcRegisterInter {
	return s.Server
}
