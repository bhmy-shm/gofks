package redisx

import (
	"context"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	// ClusterType means redis cluster.
	ClusterType = "cluster"
	// NodeType means redis node.
	NodeType = "node"
)

type (
	RedisInter interface {
		redis.Cmdable
	}

	Redis struct {
		conf *gofkConf.RedisConfig
	}
)

func NewNode(conf *gofkConf.RedisConfig) *Redis {
	r := &Redis{conf: conf}
	return r
}

// SetCtx s
func (s *Redis) SetCtx(ctx context.Context, key, value string, seconds int) error {
	node, err := s.Node()
	if err != nil {
		return err
	}
	return node.Set(ctx, key, value, time.Duration(seconds)*time.Second).Err()
}

// GetCtx s
func (s *Redis) GetCtx(ctx context.Context, key string) (val string, err error) {

	conn, err := getRedis(s.conf)
	if err != nil {
		logx.Error("GetRedis Ctx failed：", err)
		return
	}

	if val, err = conn.Get(ctx, key).Result(); err == redis.Nil {
		err = nil
	}

	return
}

// DelCtx deletes keys.
func (s *Redis) DelCtx(ctx context.Context, keys ...string) (val int, err error) {
	conn, err := getRedis(s.conf)
	if err != nil {
		logx.Error("GetRedis Ctx failed：", err)
		return
	}

	v, err := conn.Del(ctx, keys...).Result()
	if err != nil {
		logx.Error("redis Ddel is failed")
		return
	}

	val = int(v)
	return val, nil
}

// TtlCtx is the implementation of redis ttl command.
func (s *Redis) TtlCtx(ctx context.Context, key string) (val int, err error) {
	//err = s.brk.DoWithAcceptable(func() error {
	conn, err := getRedis(s.conf)
	if err != nil {
		logx.Error("ttlCtx failed:", err)
		return
	}

	duration, err := conn.TTL(ctx, key).Result()
	if err != nil {
		logx.Error("ttlCtx failed:", err)
		return
	}

	val = int(duration / time.Second)
	return
	//}, acceptable)
}

// -------------------------------------------------------------

func (s *Redis) Node() (RedisInter, error) {
	conn, err := getRedis(s.conf)
	if err != nil {
		logx.Error("get Redis Node is failed:", err)
		return nil, err
	}
	return conn, err
}

func (s *Redis) Client() *redis.Client {
	client, err := getClient(s.conf)
	if err != nil {
		logx.Error("get Redis Client failed:", err)
		return nil
	}
	return client
}

func (s *Redis) ClusterClient() *redis.ClusterClient {
	cluster, err := getCluster(s.conf)
	if err != nil {
		logx.Error("get Redis Cluster Client failed:", err)
		return nil
	}
	return cluster
}
