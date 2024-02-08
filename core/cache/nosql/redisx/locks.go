package redisx

import (
	"context"
	"errors"
	gofkConf "github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/utils/snowflake"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

/*
多实例redis锁
 如果单一的redis实例挂了，所有请求会因为拿不到锁而失败，为了提高容错性，可以使用多个分布式在不同机器上的redis实例。
 只要拿到其中大部分节点的锁就能够加锁成功。
*/

type (
	locksOptions struct {
		conf             *gofkConf.RedisConfig
		lockTTL          time.Duration   //锁过期时间
		waitCASInterval  time.Duration   //自选锁持续时间
		watchDog         chan struct{}   //看门狗结束channel
		resetTTLInterval time.Duration   // 看门狗续期时间
		clients          []*redis.Client // Redis客户端
		successClients   []*redis.Client // 加锁成功的客户端
	}
	LocksOptionFunc func(option *locksOptions)
	lockers         struct {
		key        string //被锁定的资源
		lockRandom string //锁定的随机值
		opts       *locksOptions
	}
)

func defaultLocks(clients []*redis.Client, conf *gofkConf.RedisConfig) *lockers {

	ttl := time.Millisecond * 2000

	return &lockers{
		opts: &locksOptions{
			clients:          clients,
			conf:             conf,
			watchDog:         make(chan struct{}),
			lockTTL:          ttl,
			waitCASInterval:  time.Millisecond * 200,
			resetTTLInterval: ttl / 3,
		},
	}
}

func NewLockers(clients []*redis.Client, conf *gofkConf.RedisConfig, opts ...LocksOptionFunc) LockerInter {

	lock := defaultLocks(clients, conf)

	for _, fn := range opts {
		fn(lock.opts)
	}

	return lock
}

func (l *lockers) SetKey(key string) LockerInter {
	l.key = key
	return l
}

func (l *lockers) Lock() error {

	//生成 random 随机id
	l.restId()

	wg := sync.WaitGroup{}

	//成功获取锁的客户端实例
	successClients := make(chan *redis.Client, len(l.getClient()))

	//遍历每一个redis实例
	for _, client := range l.getClient() {

		wg.Add(1)

		go func(client *redis.Client) {

			defer wg.Done()

			if ok := l.tryLock(client); ok {
				successClients <- client
				l.opts.successClients = append(l.opts.successClients, client)
			}
		}(client)
	}

	wg.Wait()
	defer close(successClients)

	//如果成功加锁的客户端 < 总客户端 的 1半+1,则表示加锁失败
	if len(successClients) < len(l.getClient())/2+1 {
		// 就算加锁失败，也要把已经获得的锁给释放掉
		for client := range successClients {
			go func(client *redis.Client) {
				ctx, cancel := context.WithTimeout(context.Background(), l.ttl())
				defer cancel()

				//关闭锁
				err := client.Eval(ctx, UnLockScript, []string{l.key}, l.lockRandom).Err()
				if err != nil {
					log.Println(err)
				}
			}(client)
		}
		l.opts.successClients = make([]*redis.Client, 0)
		return errors.New("set Locks on redisList failed")
	}

	//如果确定了加锁成功则开启看门狗
	go l.startWatchDog()

	return nil
}

func (l *lockers) tryLock(client *redis.Client) bool {
	success, err := client.SetNX(context.Background(), l.key, l.lockRandom, l.ttl()).Result()
	if err != nil {
		return false
	}
	// 加锁失败
	if !success {
		return false
	}

	return success
}

func (l *lockers) startWatchDog() {
	ticker := timex.NewTicker(l.opts.resetTTLInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.Chan():
			// 延长锁的过期时间
			for _, client := range l.opts.successClients {
				go func(client *redis.Client) {
					ctx, cancel := context.WithTimeout(context.Background(), l.opts.resetTTLInterval)
					defer cancel()
					client.Expire(ctx, l.key, l.ttl())
				}(client)
			}
		case <-l.opts.watchDog:
			// 已经解锁
			return
		}
	}
}

func (l *lockers) UnLock() (bool, error) {

	//解锁
	for _, client := range l.opts.successClients {
		go func(client *redis.Client) {
			client.Eval(context.Background(), UnLockScript, []string{l.key}, l.lockRandom)
		}(client)
	}

	// 关闭看门狗
	close(l.opts.watchDog)
	return true, nil
}

func (l *lockers) WaitCAS() error {
	//多实例加解锁比较繁琐，避免自旋操作
	return errors.New("locks can't waitCAS !~")
}

// 重置请求ID，生成一个新的请求ID
func (l *lockers) restId() {
	l.lockRandom = uuid.Must(uuid.New(), nil).String()
}

// 清除请求ID，将其设置为空。
func (l *lockers) clearId() {
	l.lockRandom = ""
}

// 雪花算法生成id, must NodeID int64
func (l *lockers) restSFId(nodeId int64) {
	l.lockRandom = snowflake.SnowflakeUUid(nodeId)
}

func (l *lockers) getClient() []*redis.Client {
	return l.opts.clients
}

func (l *lockers) ttl() time.Duration {
	return time.Duration(l.opts.lockTTL.Nanoseconds()) / (1000 * 1000)
}
