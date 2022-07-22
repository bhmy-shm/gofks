package middle

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func ContextTimeout(t time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("设置全局超时时间", t.String())
		ctx.Next()
	}
}
