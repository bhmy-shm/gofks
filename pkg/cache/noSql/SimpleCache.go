package noSql

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/bhmy-shm/gofks/pkg"
	"log"
	"reflect"
	"time"
)

type Operation interface {
	Set(key string, value interface{}, attrs ...*OperationAttr) *pkg.InterfaceResult
	Get(key string) *pkg.InterfaceResult
	MGet(keys ...string) *pkg.InterfaceResult
}

type DBGetterFunc func() interface{}

const (
	JSON = "json"
	GOB  = "gob"
)

type SimpleCache struct {
	Operation  Operation     //string类型
	Expire     time.Duration //超时时间
	DBGetter   DBGetterFunc  //缓存命中数据库
	Serializer string        //序列化方式

	Policy CachePolicy //正则匹配缓存穿透监测
}

func newCache(operation Operation, expire time.Duration, serializer string, policy CachePolicy) *SimpleCache {
	policy.SetOperation(operation)

	return &SimpleCache{
		Operation:  operation,
		Expire:     expire,
		Serializer: serializer,
		Policy:     policy,
	}
}

//设置redis操作类型
func (this *SimpleCache) SetOperation(o Operation) *SimpleCache {
	this.Operation = o
	return this
}

//设置redis穿透策略时，可选择的正则匹配规则
func (this *SimpleCache) Tactics(tics string) *SimpleCache {
	this.Policy.SetTactics(tics)
	return this
}

//写入缓存
func (this *SimpleCache) SetCache(key string, value interface{}) {

	//缓存一定是需要超时时间的，所以加上超时时间
	this.Operation.Set(key, value, WithExpire(this.Expire)).Unwrap()
}

//从缓存中读取数据
func (this *SimpleCache) GetCache(key string) (ret interface{}) {
	//1.缓存穿透的ID检查策略
	if this.Policy != nil {
		if err := this.Policy.Before(key); err != nil {
			log.Fatalln(err)
		}
	}

	switch this.Serializer {
	case JSON:
		ret = this.Operation.Get(key).UnwrapFunc(this.JsonSerializer())
	case GOB:
		ret = this.Operation.Get(key).UnwrapFunc(this.GobSerializer())
	default:
		ret = nil
	}

	retof := reflect.TypeOf(ret)

	if retof.Kind() == reflect.Ptr {
		retof = retof.Elem()
	}
	if retof.Kind() == reflect.Struct {
		vv := ret.(pkg.Value).Str()
		ret = vv
	}

	//如果是空缓存则写入：key，空value
	//if ret.(string) == "" && this.Policy != nil {
	//	this.Policy.IfNil(key, "")
	//} else {
	//	this.SetCache(key, ret) //反之将存在的数据写入redis数据库
	//}
	return ret
}

func (this *SimpleCache) GetCacheForObject(key string, obj interface{}) interface{} {
	ret := this.GetCache(key)
	if ret == nil {
		return nil
	}
	if this.Serializer == JSON {
		err := json.Unmarshal([]byte(ret.(string)), obj)
		if err != nil {
			return nil
		}
	} else if this.Serializer == GOB {
		var buf = &bytes.Buffer{}
		buf.WriteString(ret.(string))
		dec := gob.NewDecoder(buf)
		if dec.Decode(obj) != nil {
			return nil
		}
	}
	return nil
}

func (this *SimpleCache) JsonSerializer() DBGetterFunc {
	return func() interface{} {
		obj := this.DBGetter() //拿到DBGetterFunc 获取的 newsmodel
		if obj == nil {
			return nil
		}
		b, err := json.Marshal(obj)
		if err != nil {
			return nil
		}
		return string(b)
	}
}

func (this *SimpleCache) GobSerializer() DBGetterFunc {
	return func() interface{} {
		obj := this.DBGetter()
		if obj == nil {
			return nil
		}
		var buf = &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		if err := enc.Encode(obj); err != nil {
			return nil
		}
		return buf.String()
	}
}
