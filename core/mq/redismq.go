package mq

import (
	"github.com/bhmy-shm/gofks/core/cache/nosql/redisc"
	gofkConfs "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/gomodule/redigo/redis"
)

type (
	RedisMqOptions struct {
		Topic    map[string]string // 路由键
		Address  string
		Password string
	}

	RedisMqOpf func(options *RedisMqOptions)

	RedisMQ struct {
		opts *RedisMqOptions
		conn *redis.Pool //连接
		conf *gofkConfs.RedisConfig
	}
)

func defaultRedisMq(conf *gofkConfs.RedisConfig) *RedisMQ {
	return &RedisMQ{
		conn: redisc.GetRedisPoolMq(conf),
	}
}

func NewRedisMq(conf *gofkConfs.RedisConfig, opts ...RedisMqOpf) *RedisMQ {
	mq := defaultRedisMq(conf)

	for _, fn := range opts {
		fn(mq.opts)
	}
	return mq
}

func AddRedisMqTopic(topic string) RedisMqOpf {
	return func(options *RedisMqOptions) {
	}
}
