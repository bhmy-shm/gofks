package wire

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/api/svc"
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
