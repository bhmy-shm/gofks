package event

import (
	"sync"
)

var (
	Bus *EventBus
)

func InitBus() {
	//Bus = NewEventBus()
}

type (
	EventBus struct {
		Bus sync.Map
	}

	EventChannel struct {
		channel chan *EventData
		topics  map[string]struct{}
	}

	EventData struct {
		data []byte
	}
)

func NewTopic() *EventChannel {
	return &EventChannel{
		channel: make(chan *EventData),
		topics:  make(map[string]struct{}),
	}
}

func (bus *EventBus) Sub(userKey string, topic *EventChannel) *EventChannel {
	bus.Bus.Store(userKey, topic)
	return bus.GetEventChannel(userKey)
}

func (bus *EventBus) PubByUser(userKey string, params []byte) {
	bus.Bus.Range(func(key, value any) bool {

		if key.(string) == userKey {
			ch := value.(*EventChannel)
			ch.channel <- NewParams(params)
		}

		return true
	})
}

func (bus *EventBus) PubByTopic(topic string, params []byte) {
	bus.Bus.Range(func(key, value any) bool {

		ch := value.(*EventChannel)

		if _, found := ch.topics[topic]; found {
			ch.channel <- NewParams(params)
		}
		return true
	})
}

func (bus *EventBus) AddTopics(userKey string, topics ...string) {
	if len(topics) == 0 {
		return
	}

	bus.Bus.Range(func(key, value any) bool {
		if key.(string) == userKey {
			ch := value.(*EventChannel)
			for _, v := range topics {
				ch.topics[v] = struct{}{}
			}
		}
		return true
	})

}

func (bus *EventBus) RemoveTopics() {}

func (bus *EventBus) GetEventChannel(userKey string) *EventChannel {

	result, ok := bus.Bus.Load(userKey)
	if ok {
		return result.(*EventChannel) // 将 result 转换为对应的值类型
	} else {
		return nil
	}
}

func Push() {

	//生成message

	//将message，发送给指定用户的channel中
}

func NewParams(params []byte) *EventData {
	return &EventData{
		data: params,
	}
}

func (c *EventChannel) Chan() chan *EventData {
	return c.channel
}

func NewEventData(data []byte) *EventData {
	return &EventData{
		data: data,
	}
}

func (c *EventData) GetData() []byte {
	return c.data
}
