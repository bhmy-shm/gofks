package config

import (
	"reflect"
	"strconv"
)

func GetPath(path ...string) Value {

	if len(path) > 2 {
		return nil
	}

	//如果只传递一个，必须是字段key，且目前只返回一个 json对象
	if len(path) == 1 {
		n := get(path[0])
		if reflect.TypeOf(n).Kind() != reflect.Map {
			return nil
		}
		return newValue(n)
	}

	//先取第一个map，如果能够找到，再处理第二个
	data := get(path[0])

	//如果取出来的是个map则继续向下找
	if reflect.TypeOf(data).Kind() == reflect.Map {
		data = data.(map[interface{}]interface{})[path[1]]
		if data == nil {
			return nil
		}
	}

	//如果取出来的是 指定的几个yaml字段类型则返回
	return newValue(data)
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

/**
 * @Strval 获取变量的字符串值
 * @浮点型 3.0将会转换成字符串3, "3"
 * @非数值或字符类型的变量将会被转换成JSON格式字符串
 */
func Strval(value interface{}) string {
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
