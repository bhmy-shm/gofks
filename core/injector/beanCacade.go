package injector

import "sync"

type BeanCache struct {
	cache map[string]interface{}
	mutex sync.RWMutex
}

func NewBeanCache() *BeanCache {
	return &BeanCache{
		cache: make(map[string]interface{}),
		mutex: sync.RWMutex{},
	}
}

func (bc *BeanCache) Add(k string, v interface{}) {

}
