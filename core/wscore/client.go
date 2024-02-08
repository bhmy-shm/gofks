package wscore

import (
	"bytes"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/event"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/utils/snowflake"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// writeWait: 写入等待时间; 控制写入操作的最大等待时间，防止写入操作长时间阻塞
	writeWait = 10 * time.Second

	// maxMessageSize: 最大消息大小; 如果接收到的消息大小超过了这个限制，就会被丢弃;防止对等端发送过大的消息，从而保护服务器的资源
	maxMessageSize = 8192

	// Pong 等待时间; 用于控制 Pong 消息的等待时间，以便检测和处理连接的活跃性。必须大于ping 的时间。
	pongWait = 10 * time.Second

	// Ping 间隔时间; 定期向对等端发送 Ping 消息以保持连接的活跃性。必须小于 Pong 等待时间。
	pingPeriod = (pongWait * 9) / 10

	// closeGracePeriod: 强制关闭连接的等待时间; 用于控制断开连接之前的等待时间，以便允许正在进行的操作完成。
	closeGracePeriod = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type (
	ClientOptions struct {

		//客户端唯一id
		uuid string

		//客户端向服务端发送消息时存储待发送的消息。
		//当客户端想要发送消息时，它将消息写入 send 通道，而不是直接通过连接发送消息。
		//可以在发送消息时进行缓冲，提高性能和可靠性。
		send chan []byte

		//关闭webSocket通道
		done chan struct{}
	}

	ClientOption func(options *ClientOptions)

	Client struct {
		Options *ClientOptions

		//request 请求上下文
		Ctx IRequest

		//配置文件
		conf *gofkConfs.WsConfig
	}
)

func defaultClient(conf *gofkConfs.WsConfig, conn *websocket.Conn) *Client {

	var uuid string

	if conf.WS.NodeId > 0 {
		uuid = snowflake.SnowflakeUUid(conf.WS.NodeId)
	} else {
		uuid = snowflake.GenerateUniqueID(8)
	}

	return &Client{
		Options: &ClientOptions{
			uuid: uuid,
			send: make(chan []byte, conf.WS.SendBytes),
			done: make(chan struct{}),
		},
		conf: conf,
	}
}

func NewClient(conf *gofkConfs.WsConfig, conn *websocket.Conn, handle IMsgHandler, opts ...ClientOption) *Client {

	cli := defaultClient(conf, conn)

	for _, fn := range opts {
		fn(cli.Options)
	}

	cli.Ctx = NewRequest(
		WithWsConnect(conn),
		WithHandler(handle),
	)
	return cli
}

func (c *Client) GetUUid() string {
	return c.Options.uuid
}

func (c *Client) Start(hub *Hub) {
	go c.ping()

	go c.readPump(hub)

	go c.writePump()
}

func (c *Client) Stop(hub *Hub) {

	// 用户中心注销
	hub.unregister <- c

	// 关闭socket 连接
	c.Ctx.Conn().Close()

	// 通知关闭channel
	c.Options.done <- struct{}{}

	close(c.Options.done)
	close(c.Options.send)
}

// -------------------- 核心读写业务 ----------------------

// Ping 客户端定时向服务端发送心跳保持心跳，直到客户端主动关闭
func (c *Client) ping() {

	ticker := timex.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.Chan():
			if err := c.Ctx.Conn().WriteControl(websocket.PingMessage, KeepaliveMessage(), time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-c.Options.done:
			return
		}
	}
}

// ReadPump 客户端可以持续从服务端接收消息，并对接收到的消息进行处理
func (c *Client) readPump(hub *Hub) {
	defer c.Stop(hub)

	//设置读取消息的最大大小
	c.Ctx.Conn().SetReadLimit(maxMessageSize)

	//设置读取操作的截止时间
	c.Ctx.Conn().SetReadDeadline(time.Now().Add(pongWait))

	//设置处理接收到的 Pong 消息的函数
	c.Ctx.Conn().SetPongHandler(func(pingMsg string) error {
		logx.Info(pingMsg)
		return c.Ctx.Conn().SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		messageType, message, err := c.Ctx.Conn().ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch messageType {
		case websocket.TextMessage:

			//todo 处理读到的客户端协议, 并将响应的结果广播给所有在线的客户端
			c.Ctx.SetMessage(parseReadMsg(message))

			res := c.Ctx.Router().DoMsgHandler(c.Ctx)

			hub.broadcast <- res

		default:
			log.Println("client Read not TextMessage")
		}
	}
}

// WritePump 将消息从消息中心（hub）发送到 WebSocket 连接。
func (c *Client) writePump() {

	for {
		select {
		case message, ok := <-c.Options.send:

			//设置写入超时时间
			errorx.Fatal(c.Ctx.Conn().SetWriteDeadline(time.Now().Add(writeWait)))

			//如果send channel被关闭了，则客户端响应Close给服务端
			if !ok {
				err := c.Ctx.Conn().WriteMessage(websocket.CloseMessage, KeepaliveNull())
				if err != nil {
					log.Println("Client WriteMessage", err)
					errorx.Fatal(err)
				}

				return
			}

			//正常读取send channel 消息，则响应给客户端
			//c.conn.NextWriter() 可以用于写入消息到 WebSocket 连接, 但需要手动关闭写入器。
			//c.conn.WriteMessage() 直接的写入操作, 会自动处理消息的发送和关闭写入器的操作。
			w, err := c.Ctx.Conn().NextWriter(websocket.TextMessage)
			if err != nil {
				logx.Error("writer message failed:", err)
				return
			}

			w.Write(message)

			//获取当前待发送消息的数量，如果还有则继续写入
			n := len(c.Options.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Options.send)
			}

			if err = w.Close(); err != nil {
				return
			}
		}
	}
}

func parseReadMsg(message []byte) *event.Message[*event.JsonMsg] {

	//去掉制表符号
	message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

	//生成message解析工具
	jsonMsg := event.NewMessage[*event.JsonMsg](event.AbleJsonMsg())

	//使用message解析工具拆包（严格按照message协议格式才能正确解析），获取到message数据
	return jsonMsg.UnPackFromBytes(message)
}
