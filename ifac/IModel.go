package ifac

import (
	"encoding/json"
	"log"
)

type Model interface {
	String() string
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	return Models(b)
}

type Resp struct {
	Code    int
	Message string
	Data    interface{}
}

func (r *Resp) String() string {
	return "Response"
}

func Successful(data interface{}) *Resp {
	return &Resp{
		Code:    20000,
		Message: "执行成功",
		Data:    data,
	}
}

func InternalResp(data interface{}) *Resp {
	return &Resp{
		Code:    50010,
		Message: "处理失败",
		Data:    data,
	}
}
