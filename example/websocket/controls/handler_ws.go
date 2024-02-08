package controls

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/wscore"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (ws *WsCase) wsHome(ctx *gin.Context) {
	log.Println(ctx.Request.URL)

	if ctx.Request.Method != http.MethodGet {
		http.Error(ctx.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(ctx.Writer, ctx.Request, "pkg/wscore/home.html")
}

func (ws *WsCase) wsServer(ctx *gin.Context) {

	conn, err := wscore.UpWebsocket.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		log.Fatal("websocket 升级失败")
		return
	}

	//判断hub中的连接个数是否超出了最大连接数
	if ws.Hub.SessionLength() >= ws.Hub.Conf().WS.MaxConn {
		errorx.Fatal(conn.Close())
		return
	}

	//生成ws客户端
	ws.Client = wscore.NewClient(ws.Conf.GetWsCore(), conn, ws.WsWire.Ctx.GetRpcRouterHandler())

	//注册客户端到用户管理中心
	ws.Hub.RegisterAdd(ws.Client)

	//开启心跳和读写处理
	ws.Client.Start(ws.Hub)
}
