package redisx

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bhmy-shm/gofks/core/barrier"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/utils/mathx"
	"strings"
	"time"
)

// 为了避免大量缓存项在同一时间过期，使过期时间不稳定
// 将不稳定的过期时间设置为 [0.95, 1.05] * 秒数
const (
	expiryDeviation = 0.05
)

type cacheNode struct {
	cli            *Redis
	stat           *Stat                //redis命中统计
	expire         time.Duration        //缓存key过期时间
	unstableExpiry mathx.Unstable       //随机数
	barrier        barrier.SingleFlight //并发执行读取缓存任务
}

func NewSqlCache(conf *gofkConf.RedisConfig) CacheSession {
	return &cacheNode{
		cli:            NewNode(conf),
		stat:           NewStat("redis"),
		barrier:        barrier.NewFlightGroup(),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
	}
}

// TakeCtx 不带过期时间的缓存查询策略
// @request queryDB 缓存未命中后查询数据库
// @request SetCtx 查库成功后写入缓存但没有过期时间, 如果需要过期时间调用 TakeExpireCtx
func (c *cacheNode) TakeCtx(ctx context.Context, value interface{}, key string, queryDB func(val interface{}) error) error {
	return c.doTake(ctx, key, value, queryDB, func(val interface{}) error {
		return c.SetCtx(ctx, key, value) //执行完查询操作后写入缓存
	})
}

// TakeExpireCtx 带过期时间的缓存查询策略
// @request queryDB 缓存未命中后查询数据库
// @request expire 携带过期时间写入数据库，过期时间采用随机值避免大量的key 同时过期
func (c *cacheNode) TakeExpireCtx(ctx context.Context, value interface{}, key string, queryDB func(val interface{}) error) error {

	expire := c.aroundDuration(c.expire)

	return c.doTake(ctx, key, value, func(v interface{}) error {
		return queryDB(v) //执行查询操作
	}, func(v interface{}) error {
		return c.SetExpireCtx(ctx, key, v, expire) //写入缓存并设置随机过期时间
	})
}

// SetCtx 正常写入缓存
func (c *cacheNode) SetCtx(ctx context.Context, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.cli.SetCtx(ctx, key, string(data), neverExpire) //neverExpire 永不过期
}

// SetExpireCtx 写入缓存（超时时间）
func (c *cacheNode) SetExpireCtx(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.cli.SetCtx(ctx, key, string(data), int(expire.Seconds()))
}

// DelCtx 删除缓存
func (c *cacheNode) DelCtx(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	logger := logx.WithContext(ctx)
	if len(keys) > 1 && c.cli.conf.Type() == ClusterType {
		for _, key := range keys {
			if _, err := c.cli.DelCtx(ctx, key); err != nil {
				logger.Errorf("failed to clear cache with key: %q, error: %v", key, err)
				//TODO 源码c.asyncRetryDelCache(key)
				return errorx.WrapErr(err, errorx.ErrCodeRedisCacheDelFailed)
			}
		}
	} else if _, err := c.cli.DelCtx(ctx, keys...); err != nil {
		logger.Errorf("failed to clear cache with keys: %q, error: %v", strings.Join(keys, ","), err)
		//c.asyncRetryDelCache(keys...)
		return errorx.WrapErr(err, errorx.ErrCodeRedisCacheDelFailed)
	}

	return nil
}

// GetCtx 查询缓存
func (c *cacheNode) GetCtx(ctx context.Context, key string, value interface{}) error {
	return c.doGetCache(ctx, key, value)
}

// doTake 执行缓存读取任务
func (c *cacheNode) doTake(ctx context.Context,
	key string, value interface{},
	queryDB func(val interface{}) error,
	cacheAfter func(val interface{}) error) error {

	//执行读取任务
	valResult, fresh, err := c.barrier.DoEx(key, func() (interface{}, error) {

		if err := c.doGetCache(ctx, key, value); err != nil {

			//没有命中数据
			//if err == errorx.ErrCodeNotFound {
			//	return nil, err
			//}

			//执行db-Query
			if err = queryDB(value); errors.Is(err, errorx.ErrCodeNotFound) {
				return nil, err
			} else if err != nil {
				c.stat.IncrementDbFails()
				return nil, err
			}

			//db-Query 查询成功后，写入缓存当前查询的结果
			if err = cacheAfter(value); err != nil {
				return nil, err
			}
		}

		return json.Marshal(value)
	})
	if err != nil {
		return err
	}
	if fresh {
		return nil
	}

	c.stat.IncrementTotal()
	c.stat.IncrementHit()
	return json.Unmarshal(valResult.([]byte), value)
}

func (c *cacheNode) doGetCache(ctx context.Context, key string, v interface{}) error {

	c.stat.IncrementTotal() //总缓存查询数

	//Get获取缓存数据
	data, err := c.cli.GetCtx(ctx, key)
	if err != nil {
		c.stat.IncrementMiss() //未命中
		return err
	}

	if len(data) == 0 {
		c.stat.IncrementMiss() //未命中
		return errorx.ErrCodeNotFound
	}

	c.stat.IncrementHit() //缓存命中

	//缓存命中后的策略
	return c.processCache(ctx, key, data, v)
}

func (c *cacheNode) processCache(ctx context.Context, key, data string, v interface{}) error {

	var (
		err    error
		logger = logx.WithContext(ctx)
	)

	//将get出来的结果解析到指定结构体对象中
	if err = json.Unmarshal([]byte(data), v); err != nil {
		return err
	}

	//查询成功后删除缓存（查询删缓存策略）
	if _, err = c.cli.DelCtx(ctx, key); err != nil {
		logger.Errorf("delete invalid cache, "+
			"node: %s, key: %s, value: %s, error: %v", c.cli.conf.Address(), key, data, err)
		return errorx.ErrCodeRedisCacheDelFailed
	}

	return err
}

func (c *cacheNode) aroundDuration(duration time.Duration) time.Duration {
	return c.unstableExpiry.AroundDuration(duration)
}
