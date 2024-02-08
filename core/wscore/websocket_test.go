package wscore

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGinWebSocket(t *testing.T) {

	r := gin.New()

	r.GET("/echo", func(ctx *gin.Context) {

		//升级websocket协议
		conn, err := UpWebsocket.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			http.NotFound(ctx.Writer, ctx.Request)
			log.Fatal("升级失败")
			return
		}

		//循环监听发送数据
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte("hello world"))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 2)
		}
	})
	err := r.Run(":8082")
	if err != nil {
		log.Println("gin running is failed：", err)
	}
}
