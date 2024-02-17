package maps

import (
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"sync"
)

type DeploymentMap struct {
	data sync.Map // namespaces : v1.Deployments
}

func (this *DeploymentMap) Add(dep *v1.Deployment) {

	//判断当前命名空间下是否能够找到 deployments 数据
	if list, ok := this.data.Load(dep.Namespace); ok {
		list = append(list.([]*v1.Deployment), dep)
		this.data.Store(dep.Namespace, list)
	} else {
		//如果没有找到直接写入
		this.data.Store(dep.Namespace, []*v1.Deployment{dep})
	}
}

func (this *DeploymentMap) Update(dep *v1.Deployment) error {
	if list, ok := this.data.Load(dep.Namespace); ok {
		//遍历所有得数据，找到相同得进行替换
		for i, rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name {
				list.([]*v1.Deployment)[i] = dep
				break
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found", dep.Name)
}

func (this *DeploymentMap) Delete(dep *v1.Deployment) {
	if list, ok := this.data.Load(dep.Namespace); ok {
		for i, rangeDep := range list.([]*v1.Deployment) {
			if rangeDep.Name == dep.Name {
				newList := append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				this.data.Store(dep.Namespace, newList)
				break
			}
		}
	}
}

func (this *DeploymentMap) ListByNs(ns string) ([]*v1.Deployment, error) {
	//通过namespace 命名空间获取数据
	if list, ok := this.data.Load(ns); ok {
		return list.([]*v1.Deployment), nil
	}
	return nil, fmt.Errorf("record not found")
}

func (this *DeploymentMap) GetDeployment(ns string, depName string) (*v1.Deployment, error) {
	if list, ok := this.data.Load(ns); ok {
		for _, item := range list.([]*v1.Deployment) {
			if item.Name == depName {
				return item, nil
			}
		}
	}
	return nil, fmt.Errorf("record not found")
}
