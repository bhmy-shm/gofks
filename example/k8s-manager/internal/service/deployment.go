package service

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/maps"
	"github.com/bhmy-shm/gofks/example/k8s-manager/model"
)

type DeploymentService struct {
	Common *CommonService      `inject:"-"`
	DepMap *maps.DeploymentMap `inject:"-"` //数据层
}

func Deployment() *DeploymentService {
	return &DeploymentService{}
}

func (this *DeploymentService) List(ns string) ([]*model.Deployment, error) {

	depList, err := this.DepMap.ListByNs(ns)
	if err != nil {
		return nil, errorx.Wrap(err, "Dep ListByNs is failed")
	}
	var res []*model.Deployment
	for _, item := range depList {
		res = append(res, &model.Deployment{
			NameSpace: item.Namespace,
			Name:      item.Name,
			Replicas:  [3]int32{item.Status.Replicas, item.Status.AvailableReplicas, item.Status.UnavailableReplicas},
			Images:    this.Common.GetImages(*item),
			IsComplete: func() bool {
				return item.Status.Replicas == item.Status.AvailableReplicas
			}(),
			Message: func() string {
				for _, c := range item.Status.Conditions {
					if c.Type == "Available" && c.Status != "True" {
						return c.Message
					}
				}
				return ""
			}(),
		})
	}
	return res, nil
}
