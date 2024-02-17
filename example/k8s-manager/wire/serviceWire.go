package wire

import "github.com/bhmy-shm/gofks/example/k8s-manager/svc"

//注入顺序必须固定：
//wire.NewK8sHandler()
//wire.NewK8sConfig()
//wire.NewK8sMaps()
//wire.NewServiceSvc()

type ServiceWire struct {
	ServiceContext *svc.ServiceContext
	K8sMaps        *K8sMaps
	K8sConfig      *K8sConfig
	K8sHandler     *K8sHandler
}

func NewServiceWire() *ServiceWire {
	return &ServiceWire{
		ServiceContext: svc.NewServiceContext(),
	}
}

//func (w *ServiceWire) Handler() *K8sHandler {
//	w.K8sHandler = NewK8sHandler()
//	return w.K8sHandler
//}
//
//func (w *ServiceWire) Conf() *K8sConfig {
//	w.K8sConfig = NewK8sConfig()
//	return w.K8sConfig
//}
//
//func (w *ServiceWire) Maps() *K8sMaps {
//	w.K8sMaps = NewK8sMaps()
//	return w.K8sMaps
//}

func (w *ServiceWire) Context() *svc.ServiceContext {
	return w.ServiceContext
}
