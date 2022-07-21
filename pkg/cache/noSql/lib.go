package noSql

import (
	"sync"
	"time"
)

var NewsCachePool *sync.Pool

func init() {
	NewsCachePool = &sync.Pool{
		New: func() interface{} {
			return newCache(nil,
				time.Second*150, //指定超时时间
				JSON,            //指定序列化方式是json
				NewCrossPolicy("^\\d{1,5}$", time.Second*30), //todo 正则扩展
			)
		},
	}
}

/***
 * 从连接池中取出一个 redis，默认的缓存超时事件150秒，默认的缓存类型 string
 * 如果需要指定缓存类型，需要通过 SetOperation() 方法进行设置
 */
func GetCache() *SimpleCache {
	return NewsCachePool.Get().(*SimpleCache)
}

//释放链接
func ReleaseCache(cache *SimpleCache) {
	NewsCachePool.Put(cache)
}
