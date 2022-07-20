package injector

import (
	"reflect"
)

//存储传入的 bean 依赖注入的 type - value

type BeanMapper map[reflect.Type]reflect.Value

func (this BeanMapper) add(bean interface{}) {
	//拿到传入数据的 type 类型
	t := reflect.TypeOf(bean)

	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	this[t] = reflect.ValueOf(bean)
}

//传入的值都是 reflect.Type
func (this BeanMapper) get(bean interface{}) reflect.Value {
	var t reflect.Type //拿到一个类型状态

	//将这个 bean 值转换成 reflect.Type 值
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bean)
	}

	//判断这个map中是否存在相同的 reflect.Type
	if v, ok := this[t]; ok {
		return v //返回这个类型对应的 value
	}

	//如果在map中没有找到相同的
	for k, v := range this {
		//是否实现接口类型，如果当前的这个t 是一个接口，并且map中的对象，也实现了这个接口,则一样进行返回
		if t.Kind() == reflect.Interface && k.Implements(t) {
			return v
		}
	}
	return reflect.Value{}
}
