package event

import (
	"encoding/json"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/mq"
	"github.com/redis/go-redis/v9"
	"log"
)

type (
	PackageMsg interface {
		PackEncode() []byte
		PackDecode(data []byte, value interface{}) error

		SetMessage(options *MsgOptions)
		GetMessage() []byte
	}

	MsgAble interface {
		*JsonMsg | *XMLMsg
		PackageMsg
	}

	MsgOption func(options *MsgOptions)

	MsgOptions struct {
		Id uint64
		*JsonParams
		*JsonBody
		*JsonNotify
		*JsonBase
		*JsonError
		*XMLHeader
		*XMLBody
	}

	Message[T MsgAble] struct {
		Able     T
		MqRabbit *mq.RabbitMQ
		MqRedis  *redis.Client
		Opts     *MsgOptions
	}
)

func defaultMessage[T MsgAble]() *Message[T] {
	return &Message[T]{
		Opts: &MsgOptions{
			JsonNotify: defaultJsonNotify(),
			JsonBase:   defaultJsonBase(),
			XMLHeader:  defaultXmlHeader(),
			JsonParams: new(JsonParams),
		},
	}
}

func NewMessage[T MsgAble](able T, opts ...MsgOption) *Message[T] {

	msg := defaultMessage[T]()

	for _, fn := range opts {
		fn(msg.Opts)
	}

	msg.Able = able
	return msg
}

func (m *Message[T]) Pack() {
	m.Able.SetMessage(m.Opts)
}

func (m *Message[T]) UnPackFromBytes(msg []byte) *Message[T] {

	//解析 msg 然后为 opts 赋值
	err := json.Unmarshal(msg, m.Opts)
	if err != nil {
		log.Println(err)
	}

	//然后写入对应的 T
	m.Able.SetMessage(m.Opts)
	return m
}

// ------- Mq 操作 ---------

func (m *Message[T]) RedisMq(conf *gofkConf.RedisConfig) {}

func (m *Message[T]) RabbitMq(conf *gofkConf.MqConfig) {}

// --------- --------- --------- ---------

func (m *Message[T]) GetId() uint64 {
	return m.Opts.Id
}

func (m *Message[T]) GetMethod() string {
	return m.Opts.Method
}

func (m *Message[T]) GetBytes() []byte {
	return m.Able.GetMessage()
}

func (m *Message[T]) GetData() string {
	return string(m.Able.GetMessage())
}

// -------- --------- --------- ---------

func (m *Message[T]) ProgramID() string {
	return m.Opts.ProgramID
}

func (m *Message[T]) Type() string {
	return string(m.Opts.JsonNotify.Type)
}

// Send 向指定的消息队列发送消息
func (m *Message[T]) Send(topic ...string) error {

	if m.MqRedis == nil && m.MqRabbit == nil {
		return errorx.ErrCodeMqALLConnFailed
	}

	//连接指定的MQ，然后发送
	return nil
}
