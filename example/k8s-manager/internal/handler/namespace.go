package handler

import (
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/maps"
	corev1 "k8s.io/api/core/v1"
)

type NsHandler struct {
	NsMap *maps.NsMapStruct `inject:"-"`
}

func (this *NsHandler) OnAdd(obj interface{}) {
	this.NsMap.Add(obj.(*corev1.Namespace))
}

func (this *NsHandler) OnUpdate(oldObj, newObj interface{}) {
	this.NsMap.Update(newObj.(*corev1.Namespace))

}

func (this *NsHandler) OnDelete(obj interface{}) {
	if d, ok := obj.(*corev1.Namespace); ok {
		this.NsMap.Delete(d)
	}
}
