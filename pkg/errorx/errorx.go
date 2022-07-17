package errorx

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
	ErrClosed               = errors.New("pool is closed")
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				context.AbortWithStatusJSON(400, gin.H{"error": err})
			}
		}()
		context.Next()
	}
}

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
