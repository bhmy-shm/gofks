package lru

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type GCacheOption func(g *GCache)
type GCacheOptions []GCacheOption

func (opts GCacheOptions) apply(g *GCache) {
	for _, fn := range opts {
		fn(g)
	}
}

// 冗余数据。方便执行删除操作
type cacheData struct {
	key      string
	value    interface{}
	expireAt time.Time
}

func newCacheData(key string, value interface{}, expireAt time.Time) *cacheData {
	return &cacheData{key: key, value: value, expireAt: expireAt}
}

type GCache struct {
	maxSize int //限制最大key的数量，0代表不限制
	elist   *list.List
	edata   map[string]*list.Element
	lock    sync.Mutex
}

func WithMaxSize(size int) GCacheOption {
	return func(g *GCache) {
		if size > 0 {
			g.maxSize = size
		}
	}
}

func NewGCache(opts ...GCacheOption) *GCache {
	cache := &GCache{elist: list.New(), edata: make(map[string]*list.Element)}
	GCacheOptions(opts).apply(cache)
	cache.Clear() //不断删除过期数据
	return cache
}

func (this *GCache) Len() int {
	return len(this.edata)
}

func (this *GCache) Clear() {
	go func() {
		for {
			this.removeExpired()
			time.Sleep(time.Second * 1)
		}
	}()
}

// Get 获取缓存
func (this *GCache) Get(key string) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()

	if v, ok := this.edata[key]; ok {

		//代表过期
		if time.Now().After(v.Value.(*cacheData).expireAt) {
			return nil
		}
		this.elist.MoveToFront(v)
		return v.Value.(*cacheData).value
	}
	return nil
}

const NotExpireTTL = time.Hour * 24 * 365 * 10 //不过期时间

// Set 写入缓存
func (this *GCache) Set(key string, newv interface{}, ttl time.Duration) {
	this.lock.Lock()
	defer this.lock.Unlock()

	var setExpire time.Time
	if ttl == 0 {
		setExpire = time.Now().Add(NotExpireTTL)
	} else {
		setExpire = time.Now().Add(ttl)
	}

	newCache := newCacheData(key, newv, setExpire)

	if v, ok := this.edata[key]; ok {
		v.Value = newCache
		this.elist.MoveToFront(v)
	} else {
		this.edata[key] = this.elist.PushFront(newCache)

		//限制缓存最大数量是否溢出
		if this.maxSize > 0 && len(this.edata) > this.maxSize {
			this.removeOldest()
		}
	}
}

// 遍历所有元素
func (this *GCache) Print() {
	ele := this.elist.Front()
	if ele == nil {
		return
	}
	for {
		fmt.Println(ele.Value.(*cacheData).value)
		ele = ele.Next()
		if ele == nil {
			break
		}
	}
}

// 删除末尾
func (this *GCache) removeOldest() {
	back := this.elist.Back()
	if back == nil {
		return
	}
	this.removeItem(back)
}

func (this *GCache) removeItem(ele *list.Element) {
	key := ele.Value.(*cacheData).key

	//删除map里面的key
	delete(this.edata, key)

	//删除链表里面的节点
	this.elist.Remove(ele)
}

// 轮询删除过期数据
func (this *GCache) removeExpired() {
	for _, v := range this.edata {
		if time.Now().After(v.Value.(*cacheData).expireAt) {
			this.lock.Lock()
			this.removeItem(v)
			this.lock.Unlock()
		}
	}
}
