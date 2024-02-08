package wscore

import (
	"github.com/bhmy-shm/gofks/core/errorx"
	"log"
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	//升级webSocket
	conn, err := UpWebsocket.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//解析请求中的Get参数

	//判断hub中的连接个数是否超出了最大连接数
	if hub.SessionLength() >= hub.Conf().WS.MaxConn {
		errorx.Fatal(conn.Close())
		return
	}

	//生成一个客户端，并注册到用户中心
	cli := NewClient(hub.Conf(), conn, NewMsgHandler())
	hub.RegisterAdd(cli)

	//开启心跳和读写处理
	cli.Start(hub)
}

func serveWsGin() {}
