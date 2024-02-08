package redisc

import (
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/gomodule/redigo/redis"
)

func Incrementing(key string) int64 {

	var (
		id  int64
		err error
	)

	conn := redisPool.Get()
	defer conn.Close()

	id, err = redis.Int64(conn.Do("INCR", key))
	if err != nil {
		logx.Error(err)
		return -1
	}

	return id
}
