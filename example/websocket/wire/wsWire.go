package wire

import (
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisc"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/wscore"
	"github.com/bhmy-shm/gofks/example/websocket/svc"
	"github.com/gomodule/redigo/redis"
)

type WsWire struct {
	Conf   *gofkConf.Config
	Ctx    *svc.WsContext
	Hub    *wscore.Hub
	Client *wscore.Client
	Redis  *redis.Pool
}

func NewWsWire(conf *gofkConf.Config) *WsWire {
	return &WsWire{
		Conf: conf,
	}
}

func (ws *WsWire) NewContext() *svc.WsContext {
	ws.Ctx = svc.NewWsContext(ws.Conf)
	ws.Ctx.AccountRpc.RegisterAccountMethods()
	return ws.Ctx
}

func (ws *WsWire) NewHub() *wscore.Hub {
	ws.Hub = wscore.NewHub(ws.Conf.GetWsCore())
	go ws.Hub.Run()
	return ws.Hub
}

func (ws *WsWire) NewRedis() *redis.Pool {
	ws.Redis = redisc.GetRedisPool(ws.Conf.RedisConfig)
	return ws.Redis
}
