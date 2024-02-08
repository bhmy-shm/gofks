package injector

import (
	"reflect"
)

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

type BeanFactoryImpl struct {
	beanMapper *BeanMapper
	exprCache  *BeanCache
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{
		beanMapper: NewBeanMapper(),
		exprCache:  NewBeanCache(),
	}
}

func (bf *BeanFactoryImpl) GetBeanMapper() *BeanMapper {
	return bf.beanMapper
}

func (bf *BeanFactoryImpl) Set(list ...interface{}) {
	if list == nil || len(list) == 0 {
		return
	}
	for _, v := range list {
		bf.beanMapper.add(v)
	}
}

func (bf *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	data := bf.beanMapper.get(v)
	if data.IsValid() {
		return data.Interface()
	}
	return nil
}

func (bf *BeanFactoryImpl) Config(confList ...interface{}) {
	for _, cfg := range confList {

		t := reflect.TypeOf(cfg)
		v := reflect.ValueOf(cfg)

		if t.Kind() != reflect.Ptr {
			panic("required ptr object") //必须是指针对象
		}
		if t.Elem().Kind() != reflect.Struct {
			continue
		}

		bf.Set(cfg)                            //构建 cfg 与 mapper的映射
		bf.exprCache.Add(t.Elem().Name(), cfg) //构建 exprCache
		bf.Apply(cfg)                          //处理依赖注入

		//加载传入对象的所有初始化方法，依赖注入方法
		for i := 0; i < t.NumMethod(); i++ {
			method := v.Method(i)
			callRet := method.Call(nil)

			if callRet != nil && len(callRet) == 1 {
				bf.Set(callRet[0].Interface())
			}
		}
	}
}

func (bf *BeanFactoryImpl) Apply(bean interface{}) {
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
	v.Kind().String()

	//处理对象所需要注入的所有字段
	for i := 0; i < v.NumField(); i++ {

		field := v.Type().Field(i)
		tag := field.Tag.Get("inject")

		//判断对象tag是否打上了 inject 依赖注入标签
		can := v.Field(i).CanSet()
		if can && len(tag) > 0 {
			//log.Printf("wire:[%s], 注入:[%s], field:[%s-%s]\n",
			//	v.Type().String(), field.Type.String(), field.Name, field.Type.Kind().String())
			//单例模式
			if value := bf.Get(field.Type); value != nil {
				v.Field(i).Set(reflect.ValueOf(value))
				bf.Apply(value)
			}
		}
	}
}
