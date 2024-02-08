package redisx

import (
	"context"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/bhmy-shm/gofks/core/utils/snowflake"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	UnLockScript string = `
	if redis.call('GET', KEYS[1]) == ARGV[1] then
		return redis.call('DEL', KEYS[1])
	else 
		return 0 
	end
`
	KeepLockScript string = `
	if redis.call('GET', KEYS[1]) == ARGV[1] then
		return redis.call('PEXPIRE', KEYS[1], ARGV[2]) 
	else
		return 0 
	end
`
	okResult string = "OK"
)

type LockerInter interface {
	SetKey(string) LockerInter
	Lock() error
	UnLock() (bool, error)
	WaitCAS() error
}

type (
	lockOption struct {
		lockTTL         time.Duration //锁过期时间
		waitCASInterval time.Duration //自选锁持续时间
		conf            *gofkConf.RedisConfig
		redis           *redis.Client
	}
	LockOptionFunc func(option *lockOption)

	nodeLocker struct {
		key        string        //被锁定的资源
		lockRandom string        //锁定的随机值
		watchDog   chan struct{} //看门狗结束channel
		opts       *lockOption
	}
)

func defaultLocker(conf *gofkConf.RedisConfig) *nodeLocker {

	return &nodeLocker{
		opts: &lockOption{
			conf:            conf,
			lockTTL:         time.Millisecond * 1000,
			waitCASInterval: time.Millisecond * 200,
			redis:           redisClient,
		},
		watchDog: make(chan struct{}),
	}
}

func NewLock(conf *gofkConf.RedisConfig, opts ...LockOptionFunc) LockerInter {

	lk := defaultLocker(conf)

	for _, fn := range opts {
		fn(lk.opts)
	}

	return lk
}

// SetKey 设置锁键
func (l *nodeLocker) SetKey(key string) LockerInter {
	l.key = key
	return l
}

// Lock 抢锁
func (l *nodeLocker) Lock() error {
	if len(l.key) == 0 {
		return errorx.New(errorx.ErrCodeRedisKeyIsEmpty,
			errorx.WithReason("redisC TryLock failed"))
	}

	//尝试获取锁
	success, err := l.tryLock()
	if err != nil {
		return err
	}

	//没成功，则开始自旋获取锁
	if !success {
		err = l.WaitCAS()
		if err != nil {
			return err
		}
	}

	return err
}

// tryLock 真正干活的函数(抢锁)
func (l *nodeLocker) tryLock() (bool, error) {

	//新的抢锁请求ID
	l.restId()

	//抢锁 PX 单位为毫秒，键的过期时间； NX 当键不存在时才能进行设置。
	success, err := l.getClient().SetNX(context.Background(),
		l.key, l.lockRandom, l.ttl()).Result()
	if err != nil || !success {
		l.clearId()
		logx.Error("first tryLock failed")
		return false, err
	}

	//抢锁成功
	logx.Info("TryLock ok", l.key, l.lockRandom, "PX", l.opts.lockTTL.Nanoseconds()/(1000*1000))

	//启动看门狗
	go l.startWatchDog()

	return true, nil
}

// UnLock 释放锁
func (l *nodeLocker) UnLock() (bool, error) {

	if len(l.key) == 0 {
		return true, nil
	}

	if len(l.lockRandom) == 0 {
		return true, nil
	}

	//关闭锁
	err := l.getClient().Eval(context.Background(), UnLockScript, []string{l.key}, l.lockRandom).Err()
	if err != nil {
		return false, err
	}

	//关闭看门狗
	close(l.watchDog)

	return err == nil, err
}

// WaitCAS 自旋锁
func (l *nodeLocker) WaitCAS() error {

	var waitAcs time.Duration

	//如果加锁失败则自旋200ms,如果没有获取成功则返回上锁失败
	timeout := time.After(l.opts.waitCASInterval)
	for {
		select {
		case <-timeout:
			return errorx.ErrCodeRedisWaitCASTimeout //自旋到达超时时长退出
		default:

			//持续尝试获取锁
			ok, err := l.tryLock()
			if err == nil && ok {
				return nil
			}

			if !ok {
				curWait := time.Millisecond * 5
				time.Sleep(curWait)

				//如果持续等待的时间超过总时长则退出
				waitAcs = waitAcs + curWait
				log.Println("waitCAS lockTime:", waitAcs)
				if waitAcs >= l.opts.waitCASInterval {
					return errorx.ErrCodeRedisWaitCASTimeout
				}
			}
		}
	}
}

// WatchLock 看门狗
func (l *nodeLocker) startWatchDog() {
	ticker := time.NewTicker(l.opts.lockTTL / 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			//看门狗: 延长锁的过期使用时间
			dog := func() (bool, error) {

				//ok, err := l.getClient().Expire(context.Background(), l.key, l.ttl()).Result()
				//if err != nil || !ok {
				//	return false, err
				//}

				//脚本写法延时
				scriptKeys := []string{l.key}
				scriptARGV := []interface{}{l.lockRandom, l.ttl()}
				err := l.getClient().Eval(context.Background(), KeepLockScript, scriptKeys, scriptARGV).Err()
				if err != nil {
					return false, err
				}

				return true, err
			}

			if ok, err := dog(); !ok || err != nil {
				return
			}

		case <-l.watchDog: //当锁被Unlock解锁时，中止看门狗的 select
			return
		}
	}
}

// 重置请求ID，生成一个新的请求ID
func (l *nodeLocker) restId() {
	l.lockRandom = uuid.Must(uuid.New(), nil).String()
}

// 清除请求ID，将其设置为空。
func (l *nodeLocker) clearId() {
	l.lockRandom = ""
}

// 雪花算法生成id, must NodeID int64
func (l *nodeLocker) restSFId(nodeId int64) {
	l.lockRandom = snowflake.SnowflakeUUid(nodeId)
}

func (l *nodeLocker) getClient() *redis.Client {
	return l.opts.redis
}

func (l *nodeLocker) ttl() time.Duration {
	return time.Duration(l.opts.lockTTL.Nanoseconds()) / (1000 * 1000)
}

func WithLockTTL(lt time.Duration) LockOptionFunc {
	return func(option *lockOption) {
		option.lockTTL = lt
	}
}

func WithWaitCAS(wt time.Duration) LockOptionFunc {
	return func(option *lockOption) {
		option.waitCASInterval = wt
	}
}

func WithRedisClient(cli *redis.Client) LockOptionFunc {
	return func(option *lockOption) {
		option.redis = cli
	}
}

func WithRedisConf(conf *gofkConf.RedisConfig) LockOptionFunc {
	return func(option *lockOption) {
		option.conf = conf
	}
}
