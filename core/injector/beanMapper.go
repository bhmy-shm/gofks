package injector

import (
	"reflect"
	"sync"
)

// BeanMapper 用于映射对象类型 和 对象实例的关系
type BeanMapper struct {
	//map[reflect.Type]reflect.Value
	mapping *sync.Map
}

func NewBeanMapper() *BeanMapper {
	return &BeanMapper{mapping: new(sync.Map)}
}

func (bm *BeanMapper) add(bean interface{}) {
	t := reflect.TypeOf(bean)

	if t.Kind() != reflect.Ptr && t.Kind() != reflect.Interface {
		panic("require inject is must ptr object, can't(slice,map,channel,func)")
	}

	bm.mapping.Store(t, reflect.ValueOf(bean))
}

func (bm *BeanMapper) get(bean interface{}) reflect.Value {

	var (
		t reflect.Type
		v reflect.Value
	)
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bean)
	}
	if value, found := bm.mapping.Load(t); found {
		return value.(reflect.Value)
	}

	//处理接口继承
	bm.mapping.Range(func(key, value any) bool {
		k := key.(reflect.Type)
		v = value.(reflect.Value)
		if t.Kind() == reflect.Interface && k.Implements(t) {
			bm.mapping.Store(t, v)
			return false
		}
		return true
	})

	if value, ok := bm.mapping.Load(t); ok {
		return value.(reflect.Value)
	}

	return reflect.Value{}
}

func (bf *BeanMapper) Range(f func(key, value any) bool) {
	bf.mapping.Range(f)
}
