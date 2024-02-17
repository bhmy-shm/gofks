package svc

import "github.com/bhmy-shm/gofks/example/k8s-manager/internal/service"

type ServiceContext struct {
	PodSvr  *service.PodService        `inject:"-"`
	DepSvr  *service.DeploymentService `inject:"-"`
	CommSvr *service.CommonService     `inject:"-"`
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		CommSvr: service.NewCommonService(),
		PodSvr:  service.Pod(),
		DepSvr:  service.Deployment(),
	}
}

func (s *ServiceContext) CommonService() *service.CommonService {
	return s.CommSvr
}

func (s *ServiceContext) DeploymentService() *service.DeploymentService {
	return s.DepSvr
}

func (s *ServiceContext) PodService() *service.PodService {
	return s.PodSvr
}
