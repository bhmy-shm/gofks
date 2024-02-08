package redisx

import redisBloom "github.com/RedisBloom/redisbloom-go"

const (
	BloomKey      = "bloom"
	BloomCapacity = 10000
	BloomError    = 0.01
)

type BloomFilter interface {
	Reserve(key string, errorRate float64, capacity uint64) error
	Add(key string, item string) (bool, error)
	Exists(key string, item string) (bool, error)
}

type RedisBloom struct {
	client *redisBloom.Client
}

func (r *RedisBloom) Reserve(key string, errorRate float64, capacity uint64) error {
	return r.client.Reserve(key, errorRate, capacity)
}

func (r *RedisBloom) Add(key string, item string) (bool, error) {
	return r.client.Add(key, item)
}

func (r *RedisBloom) Exists(key string, item string) (bool, error) {
	return r.client.Exists(key, item)
}

func NewRedisBloom(addr, name string, pass *string) *RedisBloom {
	client := redisBloom.NewClient(addr, name, pass)
	return &RedisBloom{client: client}
}
