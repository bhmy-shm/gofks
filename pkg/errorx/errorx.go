package errorx

import (
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	HTTP_STATUS = "GOFT_STATUS"
)

var (
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
	ErrClosed               = errors.New("pool is closed")
)

func Error(err error, format ...string) {
	if err == nil {
		return
	} else {
		errMsg := err.Error()
		if len(format) > 0 {
			errMsg += format[0]
		}
		panic(errMsg)
	}
}

func Throw(err string, code int, context *gin.Context) {
	context.Set(HTTP_STATUS, code)
	panic(err)
}
