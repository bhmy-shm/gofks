package middle

import (
	"github.com/bhmy-shm/gofks/core/cache/lru"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"time"
)

/*
 	LRU令牌桶基本要素：
	1.桶的容量
	2.当前令牌数量
	3.互斥锁
*/

type Bucket struct {
	cap      int64
	tokens   int64
	rate     int64 //每秒加入令牌数
	lastTime int64 //最后一次加入令牌数的时间
	lock     sync.Mutex
}

func NewBucket(cap int64, rate int64) *Bucket {
	if cap <= 0 || rate <= 0 {
		panic("error cap")
	}
	bucket := &Bucket{cap: cap, tokens: cap, rate: rate}
	return bucket
}

func (this *Bucket) IsAccept() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	now := time.Now().Unix()
	this.tokens = this.tokens + (now-this.lastTime)*this.rate

	if this.tokens > this.cap {
		this.tokens = this.cap
	}
	this.lastTime = now

	if this.tokens > 0 {
		this.tokens--
		return true
	}
	return false
}

func Limiter(cap, rate int64) gin.HandlerFunc {
	limiter := NewBucket(cap, rate)

	return func(context *gin.Context) {
		if limiter.IsAccept() {
			context.Next()
		} else {
			context.AbortWithStatusJSON(429, "请稍后再进行访问")
		}
	}
}

func ParamLimiter(cap, rate int64, key string) gin.HandlerFunc {
	limiter := NewBucket(cap, rate)

	return func(context *gin.Context) {
		if context.Query(key) != "" {
			//如果找到这个key 则进行限流
			if limiter.IsAccept() {
				context.Next()
			} else {
				context.AbortWithStatusJSON(429, "请稍后再进行访问")
			}
		} else {
			context.Next()
		}
	}
}

var (
	IpCache *lru.GCache
)

func init() {
	IpCache = lru.NewGCache(lru.WithMaxSize(10000))
}

func CacheIpLimiter(cap, rate int64) gin.HandlerFunc {
	return func(context *gin.Context) {
		ip := context.Request.RemoteAddr
		var limiter *Bucket

		if v := IpCache.Get(ip); v != nil {
			limiter = v.(*Bucket)
		} else {
			limiter = NewBucket(cap, rate)
			log.Print("from cache")
			IpCache.Set(ip, limiter, time.Second*5)
		}

		if limiter.IsAccept() {
			context.Next()
		} else {
			context.AbortWithStatusJSON(429, gin.H{"message": "too many requests"})
		}
	}
}
