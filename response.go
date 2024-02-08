package gofks

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Code    uint64
	Reason  string
	Message string
	Data    interface{}
}

func Successful(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Resp{
		Code:    20000,
		Message: "执行成功",
		Data:    data,
	})
}

func InternalResp(ctx *gin.Context, err *errorx.Error) {
	ctx.JSON(http.StatusInternalServerError, Resp{
		Code:    err.Code,
		Reason:  err.Reason,
		Message: err.Message,
		Data:    err.Metadata,
	})
}

func ExceptionResp(ctx *gin.Context, resp errorx.ErrCode) {
	err := errorx.New(resp)
	ctx.JSON(http.StatusInternalServerError, Resp{
		Code:    err.Code,
		Reason:  err.Reason,
		Message: err.Message,
		Data:    err.Metadata,
	})
}
