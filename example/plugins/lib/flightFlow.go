package lib

import (
	"fmt"
	"github.com/bhmy-shm/gofks/core/register"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

//全国流量插件

type Flight struct {
	etcdConf *clientv3.Client
	instance *register.ServiceInstance
}

func NewFlight() *Flight {
	wea := new(Flight)
	return wea
}

func (s *Flight) Start() error {
	count := 0
	for i := 0; i < 100; i++ {
		count += 3
		fmt.Println("flight flow start 计划任务", count, time.Now())
		time.Sleep(time.Second * 5)
	}
	return nil
}

func (s *Flight) Exit() {
	return
}

func (s *Flight) Enable() bool {
	return true
}

func (s *Flight) Name() string {
	return "航站地调系统接入"
}
