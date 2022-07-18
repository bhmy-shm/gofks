package gofk

import "github.com/gin-gonic/gin"

type Fairing interface {
	//执行控制器方法前，如修改请求头信息，判断参数等
	OnRequest(ctx *gin.Context) error

	//执行控制器方法后，修改返回值内容
	OnResponse(result interface{}) (interface{}, error)
}
