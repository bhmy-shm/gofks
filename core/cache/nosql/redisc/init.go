package redisc

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var (
	redisPoolOnce sync.Once
	redisPool     *redis.Pool

	redisMqPoolOnce sync.Once
	redisMqPool     *redis.Pool
)

func GetRedisPool(conf *gofkConf.RedisConfig) *redis.Pool {
	if redisPool != nil {
		return redisPool
	}
	redisPoolOnce.Do(func() {
		redisPool = &redis.Pool{
			MaxIdle:   conf.MaxIdle(),   /* 3 最大的空闲连接数*/
			MaxActive: conf.MaxActive(), /* 8 最大的激活连接数*/
			Wait:      conf.Wait(),
			Dial: func() (redis.Conn, error) {
				if len(conf.Password()) > 0 {
					return redis.Dial(conf.Network(), conf.Address(), redis.DialPassword(conf.Password()))
				} else {
					return redis.Dial(conf.Network(), conf.Address())
				}
			},
		}
	})
	return redisPool
}

func GetRedisPoolMq(conf *gofkConf.RedisConfig) *redis.Pool {
	if redisMqPool != nil {
		return redisMqPool
	}
	redisMqPoolOnce.Do(func() {
		redisMqPool = &redis.Pool{
			MaxIdle:   conf.MaxIdle(),   /* 3 最大的空闲连接数*/
			MaxActive: conf.MaxActive(), /* 8 最大的激活连接数*/
			Wait:      conf.Wait(),
			Dial: func() (redis.Conn, error) {
				if len(conf.Password()) > 0 {
					return redis.Dial(conf.Network(), conf.Address(), redis.DialPassword(conf.Password()))
				} else {
					return redis.Dial(conf.Network(), conf.Address())
				}
			},
		}
	})
	return redisPool
}
