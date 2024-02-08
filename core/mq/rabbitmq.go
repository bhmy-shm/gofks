package mq

import (
	"fmt"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/streadway/amqp"
	"log"
)

const (
	Direct string = "direct"
	Topic  string = "topic"
)

type RabbitMQ struct {
	conn      *amqp.Connection // 连接
	channel   *amqp.Channel    // 通道
	QueueName string           // 队列名
	Exchange  string           // 交换机
	RouteKey  []string         // 路由键
	Address   string           // 连接信息
}

func NewRabbitMQ(queueName, exchange, address string, key []string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		RouteKey:  key,
		Address:   address,
	}

	var err error
	// 创建连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Address)
	errorx.Fatal(err, "创建连接错误")

	//连接amqp通道
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	errorx.Fatal(err, "获取channel失败")
	return rabbitmq
}

// DirectBind 使用路由键创建队列
func (r *RabbitMQ) DirectBind() error {
	//创建交换器
	err := r.channel.ExchangeDeclare(r.Exchange, Direct, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("exchange error = %v", err)
	}

	//根据创建成功的交换器生成队列
	q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue Declare error = %v", err)
	}

	//绑定路由键
	err = r.channel.QueueBind(q.Name, "direct", r.Exchange, false, nil)
	if err != nil {
		return fmt.Errorf("queue Bind error: %v", err)
	}
	return nil
}

// TopicBind 生成topic交换器
func (r *RabbitMQ) TopicBind() error {
	err := r.channel.ExchangeDeclare(r.Exchange, Topic,
		true, false, false, false, nil)
	errorx.Fatal(err, "生成topic交换器出现问题")

	//生成队列
	q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
	errorx.Fatal(err, "生成topic交换器队列出现问题")

	//绑定所有的路由键值
	for _, v := range r.RouteKey {
		err = r.channel.QueueBind(q.Name, v, r.Exchange, false, nil)
		errorx.Fatal(err, "topic绑定队列失败")
	}

	return nil
}

// AddTopicRouteKey 增加路由键
func (r *RabbitMQ) AddTopicRouteKey(routes ...string) error {
	q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
	errorx.Fatal(err, "绑定队列失败")

	//添加新的routeKey 到队列中
	//绑定队列到 exchange 中
	for _, rr := range routes {
		err = r.channel.QueueBind(q.Name, rr, r.Exchange, false, nil)
		errorx.Fatal(err, "绑定队列新的 RouteKey 出现失败")
	}
	return nil
}

// SendMessage 通过指定的 routeKey 发送消息
func (r *RabbitMQ) SendMessage(msg amqp.Publishing) error {
	if len(r.RouteKey) == 0 {
		//如果没有 RouteKey 代表是 direct 的方式发送消息
		return r.channel.Publish(r.Exchange, Direct, false, false, msg)
	}

	//如果有具体的路由键，按照RouteKey发送
	for _, route := range r.RouteKey {
		err := r.channel.Publish(r.Exchange, route, false, false, msg)
		errorx.Fatal(err, "发送消息失败")
	}
	return nil
}

// Consumer 创建消费者
func (r *RabbitMQ) Consumer(consumer string, callback func(<-chan amqp.Delivery)) error {
	msg, err := r.channel.Consume(r.QueueName, consumer, false, false, false, false, nil)
	if err != nil {
		log.Println("consumer err=", err)
	}
	callback(msg)
	return nil
}

// Close 关闭MQ
func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
