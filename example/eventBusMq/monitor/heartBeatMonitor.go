package monitor

import (
	"context"
	"fmt"
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisc"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/event"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	HeartBeatKey = "Auto-Increment-HeartBeat-key"
)

func defaultHeartBeat[T event.MsgAble]() *event.Message[T] {
	return &event.Message[T]{
		Able: nil,
		Opts: &event.MsgOptions{
			Id: uint64(redisc.Incrementing(HeartBeatKey)),
			JsonNotify: &event.JsonNotify{
				ProgramID:     uuid.New().String(),
				Type:          event.MonitorHeartBeat,
				KeepaliveTime: time.Second * 10,
			},
		},
	}
}

// NewHearBeat 创建一个心跳类型消息
func NewHearBeat[T event.MsgAble](opts ...event.MsgOption) *event.Message[T] {

	h := defaultHeartBeat[T]()

	for _, fn := range opts {
		fn(h.Opts)
	}

	return h
}

func RegisterMonitor(cli *redis.Client, msg event.EventMonitorItem) error {

	if msg.Type() != string(event.MonitorHeartBeat) {
		return errorx.ErrCodeMsgTypeMismatch
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//从缓存中判断注册的 name（key）hash 是否存在
	ok, err := cli.Exists(ctx, msg.ProgramID()).Result()
	if err != nil {
		return err
	}

	//ok == 1 表示key存在。
	if ok == 1 {
		//如果存在则判断是否包含当前传递过来的类型，如果类型也存在则返回冲突
		maps, err := cli.HGetAll(ctx, msg.ProgramID()).Result()
		if err != nil {
			return err
		}

		for k, v := range maps {
			if msg.Type() == k {
				return fmt.Errorf("注册的消息类型已经存在【%s:%s】", k, v)
			}
		}
	}

	var setMap = make(map[string]interface{})
	setMap[msg.Type()] = 0 //某种消息=1，代表发送心跳
	if err = cli.HSet(ctx, msg.ProgramID(), setMap).Err(); err != nil {
		return err
	}

	return nil
}
