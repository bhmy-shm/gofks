package maps

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sync"
)

type PodMap struct {
	data sync.Map // namespaces : corev1.Pod
}

func (this *PodMap) ListByNs(ns string) []*corev1.Pod {
	if list, ok := this.data.Load(ns); ok {
		return list.([]*corev1.Pod)
	}
	return nil
}

func (this *PodMap) Get(ns string, podName string) *corev1.Pod {
	if list, ok := this.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			if pod.Name == podName {
				return pod
			}
		}
	}
	return nil
}

func (this *PodMap) Add(pod *corev1.Pod) {
	if list, ok := this.data.Load(pod.Namespace); ok {
		list = append(list.([]*corev1.Pod), pod)
		this.data.Store(pod.Namespace, list)
	} else {
		this.data.Store(pod.Namespace, []*corev1.Pod{pod})
	}
}

func (this *PodMap) Update(pod *corev1.Pod) error {
	if list, ok := this.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				list.([]*corev1.Pod)[i] = pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found", pod.Name)
}

func (this *PodMap) Delete(pod *corev1.Pod) {
	if list, ok := this.data.Load(pod.Namespace); ok {
		for i, range_pod := range list.([]*corev1.Pod) {
			if range_pod.Name == pod.Name {
				newList := append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				this.data.Store(pod.Namespace, newList)
				break
			}
		}
	}
}

// ListByLabels 根据标签获取 POD列表
func (this *PodMap) ListByLabels(ns string, labels []map[string]string) ([]*corev1.Pod, error) {
	ret := make([]*corev1.Pod, 0)
	if list, ok := this.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			for _, label := range labels {
				if reflect.DeepEqual(pod.Labels, label) { //标签完全匹配
					ret = append(ret, pod)
				}
			}
		}
		return ret, nil
	}
	return nil, fmt.Errorf("pods not found ")
}

func (this *PodMap) DEBUGListByNS(ns string) []*corev1.Pod {
	ret := make([]*corev1.Pod, 0)
	if list, ok := this.data.Load(ns); ok {
		for _, pod := range list.([]*corev1.Pod) {
			ret = append(ret, pod)
		}

	}
	return ret
}
