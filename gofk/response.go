package gofk

type Resp struct {
	Code    int
	Message string
	Data    interface{}
}

func (r *Resp) String() string {
	return "Response"
}

func Successful(data interface{}) Resp {
	return Resp{
		Code:    20000,
		Message: "执行成功",
		Data:    data,
	}
}

func InternalResp(data interface{}) Resp {
	return Resp{
		Code:    50010,
		Message: "处理失败",
		Data:    data,
	}
}
