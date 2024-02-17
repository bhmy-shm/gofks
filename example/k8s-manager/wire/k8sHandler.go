package wire

import "github.com/bhmy-shm/gofks/example/k8s-manager/internal/handler"

type K8sHandler struct {
}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}

func (this *K8sHandler) DepHandlers() *handler.DepHandler {
	return &handler.DepHandler{}
}

func (this *K8sHandler) PodHandler() *handler.PodHandler {
	return &handler.PodHandler{}
}

func (this *K8sHandler) NsHandler() *handler.NsHandler {
	return &handler.NsHandler{}
}
