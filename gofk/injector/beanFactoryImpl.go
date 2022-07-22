package injector

import (
	"reflect"
)

//依赖注入

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

type BeanFactoryImpl struct {
	beanMapper BeanMapper
	ExprMap    map[string]interface{}
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper),
		ExprMap: make(map[string]interface{})}
	//ExprMap: make(map[string]interface{}),
}

func (b *BeanFactoryImpl) Set(list ...interface{}) {
	if list == nil || len(list) == 0 {
		return
	}
	for _, v := range list {
		b.beanMapper.add(v)
	}
}

func (b *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	getv := b.beanMapper.get(v)

	//如果v不是零值，则返回true。
	//如果v是零值，则返回false。 如果IsValid返回false，则除String外的所有其他方法都为false。
	if getv.IsValid() {
		return getv.Interface()
	}
	return nil
}

func (b *BeanFactoryImpl) Config(cfgs ...interface{}) {
	for _, cfg := range cfgs {

		//拿到传入的对象类型
		t := reflect.TypeOf(cfg)

		//如果是对象指针则拿出对象
		if t.Kind() != reflect.Ptr {
			panic("required ptr object")
		}
		//强制要求是对象
		if t.Elem().Kind() != reflect.Struct {
			continue
		}

		//把对象 加入依赖注入map池
		b.Set(cfg)

		//处理依赖注入，处理注入对象的每一个字段，
		b.Apply(cfg)

		//取反射对象的所有方法
		v := reflect.ValueOf(cfg)
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)

			//调用方法
			callRet := method.Call(nil)

			//如果调用结果不是空的，则写入方法的调用结果到依赖注入列表
			if callRet != nil && len(callRet) == 1 {
				b.Set(callRet[0].Interface())
			}
		}
	}
}

func (b *BeanFactoryImpl) Apply(bean interface{}) {
	//拿到的是一个对象

	if bean == nil {
		return
	}
	v := reflect.ValueOf(bean)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	//循环遍历对象的每一个字段
	for i := 0; i < v.NumField(); i++ {

		//拿到字段
		field := v.Type().Field(i)

		//处理字段的 inject 标签
		if v.Field(i).CanSet() && field.Tag.Get("inject") != "" {

			if field.Tag.Get("inject") != "-" {
				//ret := expr.BeanExpr(field.Tag.Get("inject"), this.ExprMap)
				//if ret != nil && !ret.IsEmpty() {
				//	retValue := ret[0]
				//	if retValue != nil {
				//		v.Field(i).Set(reflect.ValueOf(retValue))
				//		this.Apply(retValue)
				//	}
				//}
			} else { //单例模式
				//传入字段类型，字段类型是否存在依赖注入列表
				if getV := b.Get(field.Type); getV != nil {
					v.Field(i).Set(reflect.ValueOf(getV))
					b.Apply(getV)
				}
			}
		}
		// 处理其它标签
	}
}

func (b *BeanFactoryImpl) GetBeanMapper() BeanMapper {
	return b.beanMapper
}
