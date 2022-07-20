package gofk

import "github.com/gin-gonic/gin"

type resp struct {
	Code    int
	Message string
	Data    interface{}
}

func Successful(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, resp{
		Code:    20000,
		Message: "执行成功",
		Data:    data,
	})
}

func InternalResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, resp{
		Code:    50010,
		Message: "处理失败",
		Data:    data,
	})
}
