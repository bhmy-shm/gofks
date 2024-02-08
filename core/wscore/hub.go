package wscore

import (
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"log"
	"sync"
)

type (
	Hub struct {

		//wsCore配置文件
		conf *gofkConfs.WsConfig

		// 用户注册中心
		clients *sync.Map

		// 注册客户端的通道
		register chan *Client

		// 注销客户端的通道
		unregister chan *Client

		// broadcast 广播通道当有消息需要广播给所有客户端时，可以将消息数据写入 broadcast 通道;
		// 然后 Hub 会将该消息发送给所有注册的客户端
		broadcast chan []byte

		// 当前客户端数量
		length int64

		lock sync.Mutex
	}
)

func NewHub(conf *gofkConfs.WsConfig) *Hub {
	return &Hub{
		clients:    new(sync.Map),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		conf:       conf,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register: //注册
			addClient(h.clients, c)
			h.length++

		case c := <-h.unregister: //注销
			deleteClient(h.clients, c)
			close(c.Options.send)
			h.length--

		case message := <-h.broadcast: //所有客户端广播消息
			h.clients.Range(func(key, value any) bool {

				cli := key.(*Client)
				select {
				case cli.Options.send <- message:
					log.Println("[debug] 向客户端广播 message消息", string(message))

				//default 分支是用来处理发送消息失败的情况。
				//当向某个客户端的 send 通道发送消息时，如果通道已满，发送操作就会被阻塞。
				//在这种情况下，default 分支会被执行。
				default:
					close(cli.Options.send)
					deleteClient(h.clients, cli)
					h.length--
				}

				return true
			})
		}
	}
}

func (h *Hub) Clients() *sync.Map {
	return h.clients
}

func (h *Hub) RegisterAdd(cli *Client) {
	h.register <- cli
}

func (h *Hub) UnRegister() chan *Client {
	return h.unregister
}

func (h *Hub) Broadcast() chan []byte {
	return h.broadcast
}

func (h *Hub) SessionLength() int64 {
	return h.length
}

func (h *Hub) Conf() *gofkConfs.WsConfig {
	return h.conf
}

func addClient(p *sync.Map, cli *Client) {
	p.Store(cli, true)
}

func deleteClient(p *sync.Map, cli *Client) {
	p.Delete(cli)
}

func rangeClient(p *sync.Map) int {
	count := 0
	p.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
