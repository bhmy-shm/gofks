package wire

import (
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/maps"
)

type K8sMaps struct{}

func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}

func (this *K8sMaps) InitDepMap() *maps.DeploymentMap {
	return &maps.DeploymentMap{}
}

func (this *K8sMaps) InitPodMap() *maps.PodMap {
	return &maps.PodMap{}
}

func (this *K8sMaps) InitNsMap() *maps.NsMapStruct {
	return &maps.NsMapStruct{}
}
