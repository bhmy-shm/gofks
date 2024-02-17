package service

import v1 "k8s.io/api/apps/v1"

type CommonService struct{}

func NewCommonService() *CommonService {
	return &CommonService{}
}

func (this *CommonService) GetImages(dep v1.Deployment) string {
	return ""
}
