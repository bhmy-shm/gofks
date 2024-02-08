package pkg

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"reflect"
	"strconv"
)

type (
	Encoder interface {
		Encode(interface{}) ([]byte, error) //编码
		Decode([]byte, interface{}) error   //解码
		String() string                     //返回文件名字
	}

	ValueInter[T wrapAble] interface {
		Value() WrapInter[T]
	}

	value[T wrapAble] struct {
		n interface{}
	}
)

func newValue[T wrapAble](data T) *value[T] {
	return &value[T]{n: data}
}

func (v *value[T]) Get() interface{} {
	return v.n
}

func (v *value[T]) Value() WrapInter[T] {
	if s, ok := (v.n).(T); ok {
		return newWrap(s, nil)
	}
	var zero T
	return newWrap(zero, errorx.ErrCodeAssertionValue)
}

func (v *value[T]) Slice() ([]T, error) {
	if a, ok := (v.n).([]T); ok {
		return a, nil
	}
	return nil, errorx.New(errorx.ErrCodeAssertionValue, errorx.WithReason("to Slice failed"))
}

func GetPath[T wrapAble](path ...string) ValueInter[T] {

	if len(path) > 2 {
		return nil
	}

	//如果只传递一个，必须是字段key，且目前只返回一个 json对象
	if len(path) == 1 {
		n := get(path[0])
		if reflect.TypeOf(n).Kind() != reflect.Map {
			return nil
		}

		if tv, ok := n.(T); ok {
			return newValue[T](tv)
		} else {
			errorx.Fatal(errorx.ErrCodeGetPathValue)
		}
	}

	//先取第一个map，如果能够找到，再处理第二个
	data := get(path[0])

	//如果取出来的是个map则继续向下找
	if reflect.TypeOf(data).Kind() == reflect.Map {
		data = data.(FileMap)[path[1]]
		if data == nil {
			return nil
		}
	}

	//如果取出来的是 指定的几个yaml字段类型则返回
	if tv, ok := data.(T); ok {
		return newValue[T](tv)
	} else {
		errorx.Fatal(errorx.ErrCodeGetPathValue)
	}

	return nil
}

func get(p string) interface{} {
	m := GlobalConf

	//从内存中取出配置信息
	val, ok := m[p]
	if !ok {
		return nil
	}
	return val
}

/*
StrVal 获取变量的字符串值
* @浮点型 3.0将会转换成字符串3, "3"
* @非数值或字符类型的变量将会被转换成JSON格式字符串
*/
func StrVal(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		return ""
	}
	return key
}
