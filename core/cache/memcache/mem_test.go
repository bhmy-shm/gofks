package memcache

import (
	"fmt"
	"testing"
	"time"
)

func TestGetValueSize(t *testing.T) {
	GetValueSize(1)
	GetValueSize(100000)
	GetValueSize("ashshah")
	GetValueSize(21.5)
}

func TestMemCache(t *testing.T) {
	cache := NewMemCache()
	cache.SetMaxMemory("100MB")

	cache.Set("int", 1, time.Second*2)
	//for i := 0; i < 5; i++ {
	//	time.Sleep(time.Millisecond * 200)
	//	log.Println(cache.Get("int"))
	//}

	cache.Set("maps", map[string]interface{}{"a": 1}, time.Second*2)
	//for i := 0; i < 5; i++ {
	//	time.Sleep(time.Millisecond * 500)
	//	log.Println(cache.Get("maps"))
	//}
	fmt.Println(cache.Keys())
}
