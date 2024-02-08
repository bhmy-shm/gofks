package memcache

import (
	"fmt"
	"sync"
	"time"
)

type MemCache struct {
	//最大内存
	maxMemorySize int64

	//最大内存字符串表示字段
	maxMemoryStr string

	//当前所用内存
	useMemorySize int64

	//内存存储data
	values map[string]*memCacheValue

	//读写锁(并发读，单独写)
	lock sync.RWMutex

	//清除缓存的时间间隔
	clearTimeInterval time.Duration
}

type memCacheValue struct {
	value      interface{}   //value值
	expireTime time.Time     //过期时间点
	expire     time.Duration //过期时长
	size       int64         //当前对象的size大小
}

func NewMemCache() *MemCache {
	mc := &MemCache{
		values:            make(map[string]*memCacheValue, 0),
		clearTimeInterval: time.Second, //秒级别
	}
	go mc.clearExpireKey() //定期清空缓存策略
	return mc
}

// SetMaxMemory size 1KB 100KB 1MB 2MB 1GB
func (mc *MemCache) SetMaxMemory(size string) bool {
	mc.maxMemorySize, mc.maxMemoryStr = ParseSize(size)

	return false
}

func (mc *MemCache) Set(key string, val interface{}, expire time.Duration) bool {
	//锁
	mc.lock.Lock()
	defer mc.lock.Unlock()

	//赋值
	v := &memCacheValue{
		value:      val,
		expireTime: time.Now().Add(expire),
		size:       GetValueSize(val),
	}

	//更新与添加操作
	mc.del(key)
	mc.add(key, v)
	//限制最大内存
	if mc.useMemorySize > mc.maxMemorySize {
		mc.del(key)
		panic(fmt.Sprintf("max MemorySize=%d is oom", mc.maxMemorySize))
	}

	return false
}

func (mc *MemCache) get(key string) (*memCacheValue, bool) {
	val, ok := mc.values[key]
	return val, ok
}

func (mc *MemCache) del(key string) {
	tmp, ok := mc.get(key)
	if ok && tmp != nil {
		//删掉内存占用
		mc.useMemorySize -= tmp.size
		//删除kv
		delete(mc.values, key)
	}
}

func (mc *MemCache) add(key string, val *memCacheValue) {
	mc.values[key] = val
	mc.useMemorySize += val.size
}

func (mc *MemCache) Get(key string) (interface{}, bool) {
	//锁
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	//获取缓存内容
	mcv, ok := mc.get(key)
	if ok {
		//判定缓存是否过期 (但是这种属于只有读取时才会删除)
		if mcv.expire != 0 && mcv.expireTime.Before(time.Now()) {
			//过期时间早于当前日期，过期
			mc.del(key)
			return nil, false
		}
		return mcv.value, ok
	}
	return nil, false
}

func (mc *MemCache) Del(key string) bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.del(key)
	return true
}

// Exists 判断key是否存在
func (mc *MemCache) Exists(key string) bool {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	_, ok := mc.values[key]
	return ok
}

// Flush 清空
func (mc *MemCache) Flush() bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mc.values = make(map[string]*memCacheValue, 0)
	mc.useMemorySize = 0
	return true
}

// Keys 获取缓存中所有keys的数量
func (mc *MemCache) Keys() int64 {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	return int64(len(mc.values))
}

// clear 轮询所有内容，保证过期第一时间删除
func (mc *MemCache) clearExpireKey() {
	//定时器
	timerTicker := time.NewTicker(mc.clearTimeInterval)
	defer timerTicker.Stop()
	for {
		select {
		case <-timerTicker.C:
			//清理缓存
			for key, item := range mc.values {
				if item.expire != 0 && time.Now().After(item.expireTime) {
					mc.lock.Lock()
					mc.del(key)
					mc.lock.Unlock()
				}
			}
		}
	}
}
