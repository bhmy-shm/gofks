package wscore

import (
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisc"
	"github.com/bhmy-shm/gofks/core/event"
	"time"
)

func KeepaliveNull() []byte {
	return []byte{}
}

func KeepaliveMessage() []byte {

	msg := event.NewMessage(
		event.AbleJsonMsg(),
		event.WithId(uint64(redisc.Incrementing(event.GlobalMidIncrementKey))),
		event.WithMethod("keepalive-ping"),
		event.WithJsonNotify(&event.JsonNotify{
			Type:      event.WebSocketKeepalive,
			Timestamp: time.Now().Unix(),
		}),
	)

	msg.Pack()

	return msg.Able.PackEncode()
}
