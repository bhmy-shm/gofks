package pkg

import (
	"bytes"
	"errors"
)

type value struct {
	n interface{}
	m map[interface{}]interface{}
}

func newValue(data interface{}) Value {
	if data == nil {
		return nil
	}
	res := value{}

	switch data.(type) {
	case map[interface{}]interface{}:
		res.m = data.(map[interface{}]interface{})
	default:
		res.n = data
	}
	return res
}

func UnwrapValue(data interface{}) *value {
	if data == nil {
		return nil
	}
	return &value{n: data}
}

func (v value) Bool() (bool, error) {
	if s, ok := (v.n).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

func (v value) Int() (int, error) {
	if s, ok := (v.n).(int); ok {
		return s, nil
	}
	return 0, errors.New("type assertion to int failed")
}

func (v value) String() (string, error) {
	if s, ok := (v.n).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to String failed")
}

func (v value) Float64() (float64, error) {
	if s, ok := (v.n).(float64); ok {
		return s, nil
	}
	return 0.0, errors.New("type assertion to Float64 failed")
}

func (v value) Slice() ([]interface{}, error) {
	if a, ok := (v.n).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}

func (v value) Interface() interface{} {
	return v.n
}

func (v value) StringSlice() []string {
	arr, err := v.Slice()
	if err != nil {
		return nil
	}
	retArr := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		if arr[i] == nil {
			retArr[i] = ""
			continue
		}
		s, ok := arr[i].(string)
		if !ok {
			return nil
		}
		retArr[i] = s
	}
	return retArr
}

func (v value) Str() string {
	if s, ok := (v.n).(string); ok {
		return s
	}
	return ""
}

func (v value) Json() (string, error) {
	if v.m == nil {
		return "", errors.New("Value does not contain a map type to map JSON ")
	}

	count := 0
	l := len(v.m)
	buf := bytes.Buffer{}

	buf.WriteByte('{')
	for k, vv := range v.m {
		count++
		buf.WriteByte('"')
		buf.WriteString(k.(string))
		buf.WriteByte('"')

		buf.WriteByte(':')

		buf.WriteByte('"')
		buf.WriteString(Strval(vv))
		buf.WriteByte('"')

		if count < l {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.String(), nil
}

func (v value) SetN(n interface{}) value {
	v.n = n
	return v
}
