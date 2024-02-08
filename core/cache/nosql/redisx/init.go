package redisx

import (
	"context"
	"crypto/tls"
	"fmt"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/utils/syncx"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

const (
	defaultSlowThreshold = time.Millisecond * 100
	neverExpire          = 0
)

var (
	slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
)

func getRedis(conf *gofkConf.RedisConfig) (RedisInter, error) {
	switch conf.Type() {
	case ClusterType:
		return getCluster(conf)
	case NodeType:
		return getClient(conf)
	default:
		return nil, fmt.Errorf("redis type '%s' is not supported", conf.Type())
	}
}

var (
	redisClientOnce  sync.Once
	redisClusterOnce sync.Once

	redisClient  *redis.Client
	redisCluster *redis.ClusterClient
)

func getClient(conf *gofkConf.RedisConfig) (*redis.Client, error) {

	if redisClient != nil {
		return redisClient, nil
	}

	var (
		err       error
		tlsConfig *tls.Config
	)

	redisClientOnce.Do(func() {
		if conf.Tls() {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		redisClient = redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     conf.Address(),
			Password: conf.Password(), //密码
			DB:       0,               // redis数据库

			//连接池容量及闲置连接数量
			PoolSize:     15, // 连接池数量
			MinIdleConns: 10, //好比最小连接数
			//超时
			DialTimeout:  5 * time.Second, //连接建立超时时间
			ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			//命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

			//安全策略
			TLSConfig: tlsConfig,
		})

		pong, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(fmt.Errorf("connect error:%s", err))
		}

		redisClient.AddHook(durationHook)

		logx.Info(pong)
	})

	return redisClient, err
}

func getCluster(conf *gofkConf.RedisConfig) (*redis.ClusterClient, error) {

	if redisCluster != nil {
		return redisCluster, nil
	}

	var (
		err       error
		tlsConfig *tls.Config
	)

	redisClusterOnce.Do(func() {
		if conf.Tls() {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		redisCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{conf.Address()},
			Password: conf.Password(), //密码

			//连接池容量及闲置连接数量
			PoolSize:     15, // 连接池数量
			MinIdleConns: 10, //好比最小连接数
			//超时
			DialTimeout:  5 * time.Second, //连接建立超时时间
			ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
			PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

			//命令执行失败时的重试策略
			MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
			MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

			//安全策略
			TLSConfig: tlsConfig,
		})

		pong, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(fmt.Errorf("connect error:%s", err))
		}

		redisCluster.AddHook(durationHook)

		logx.Info(pong)
	})

	return redisCluster, err
}
