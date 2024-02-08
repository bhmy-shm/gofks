package wscore

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var UpWebsocket = websocket.Upgrader{

	//读缓冲区大小
	ReadBufferSize: 1024 * 1024,

	//写缓冲区大小
	WriteBufferSize: 1024 * 1024,

	//传输数据压缩（需要cpu进行计算）
	//推荐使用：在大量数据和带宽受限场景下使用
	EnableCompression: false,

	//Origin 校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func checkOrigins(r *http.Request) bool {
	// 设置你期望的 Origins
	expectedOrigins := []string{"https://example.com", "https://sub.example.com"}

	// 获取请求的 Origin
	origin := r.Header.Get("Origin")

	// 检查请求的 Origin 是否在期望的列表中
	for _, expectedOrigin := range expectedOrigins {
		if origin == expectedOrigin {
			return true
		}
	}
	return false
}

func internalError(ws *websocket.Conn, err error) error {
	return ws.WriteMessage(
		websocket.TextMessage,
		[]byte("Internal server error: "+err.Error()),
	)
}
