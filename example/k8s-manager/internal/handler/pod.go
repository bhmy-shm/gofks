package handler

import (
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/maps"
	v1 "k8s.io/api/core/v1"
	"log"
)

type PodHandler struct {
	PodMap *maps.PodMap `inject:"-"`
}

func (this *PodHandler) OnAdd(obj interface{}) {
	this.PodMap.Add(obj.(*v1.Pod))
}

func (this *PodHandler) OnUpdate(oldObj, newObj interface{}) {
	err := this.PodMap.Update(newObj.(*v1.Pod))
	if err != nil {
		log.Println(err)
	}
}

func (this *PodHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*v1.Pod); ok {
		this.PodMap.Delete(d)
	}
}
