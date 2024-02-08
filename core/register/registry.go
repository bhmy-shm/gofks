package register

import (
	"context"
	"encoding/json"
)

// Registrar 服务注册
type Registrar interface {
	// Register 注册
	Register(ctx context.Context, service *ServiceInstance) error
	// Deregister 反注册
	Deregister(ctx context.Context, service *ServiceInstance) error
}

// Discovery 服务发现
type Discovery interface {
	// GetService 根据服务名称返回内存中的服务实例。
	GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	// Watch 根据服务名称创建一个监视程序。
	Watch(ctx context.Context, serviceName string) (Watcher, error)
}

// Watcher 服务监听
type Watcher interface {
	/*
		Next在以下两种情况下返回服务:
		1.第一次观看且服务实例列表不为空。
		2.发现的任何服务实例更改。
		如果不满足以上两个条件，它将阻塞，直到超过上下文截止日期或取消
	*/
	Next() ([]*ServiceInstance, error)

	Stop() error
}

// ServiceInstance 服务的实例
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func marshal(s *ServiceInstance) (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (s *ServiceInstance, err error) {
	err = json.Unmarshal(data, &s)
	return
}
