package handler

import (
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/maps"
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/service"
	v1 "k8s.io/api/apps/v1"
	"log"
)

type DepHandler struct {
	DepMap     *maps.DeploymentMap        `inject:"-"`
	DepService *service.DeploymentService `inject:"-"`
}

func (this *DepHandler) OnAdd(obj interface{}) {
	this.DepMap.Add(obj.(*v1.Deployment))
	//list, _ := this.DepService.List(obj.(*v1.Deployment).Namespace)
	//wscore.ClientMap.SendAllDepList(list)
}

func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	err := this.DepMap.Update(newObj.(*v1.Deployment))
	if err != nil {
		log.Println(err)
	} else {
		//list, _ := this.DepService.List(newObj.(*v1.Deployment).Namespace)
		//wscore.ClientMap.SendAllDepList(list)
	}
}

func (this *DepHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*v1.Deployment); ok {
		this.DepMap.Delete(d)
		//list, _ := this.DepService.List(obj.(*v1.Deployment).Namespace)
		//wscore.ClientMap.SendAllDepList(list)
	}
}
