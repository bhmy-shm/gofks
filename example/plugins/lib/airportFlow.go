package lib

import (
	"fmt"
	"github.com/bhmy-shm/gofks/core/register"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//全国流量插件

type Flow struct {
	etcdConf *clientv3.Client
	instance *register.ServiceInstance
}

func NewFlow() *Flow {
	wea := new(Flow)
	return wea
}

func (s *Flow) Start() error {
	count := 0
	for i := 0; i < 100; i++ {
		count += 3
		fmt.Println("airport flow start 计划任务", count, time.Now())
		time.Sleep(time.Second * 3)
	}
	return nil
}

func (s *Flow) Exit() {
	return
}

func (s *Flow) Enable() bool {
	return true
}

func (s *Flow) Name() string {
	return "全国流量接入"
}
