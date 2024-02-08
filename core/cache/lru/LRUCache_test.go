package lru

import (
	"container/list"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	r := gin.New()
	//r.Use(middle.CacheIpLimiter(10, 1))
	r.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"message": "ok"})
	})
	r.Run(":8081")
}

func TestList(t *testing.T) {
	ll := list.New()

	ll.PushFront("no1")
	n2 := ll.PushFront("no2")
	ll.InsertBefore("no3", n2)
	ele := ll.Front()
	if ele == nil {
		log.Fatalln("nil element")
	}
	for {
		fmt.Println(ele.Value)
		if ele.Next() == nil {
			break
		}
		ele = ele.Next()
	}
}

func TestGCache(t *testing.T) {
	cache := NewGCache(WithMaxSize(3))
	cache.Set("name", "sunhaiming", time.Second*3)
	cache.Set("age", "19", 0)
	cache.Set("sex", "男", 0)
	//cache.Set("abc", "abc")
	//cache.Get("age")
	//cache.Set("sex", "女")
	//for {
	//	fmt.Printf("name=%v age=%v sex=%v\n",
	//		cache.Get("name"),
	//		cache.Get("age"),
	//		cache.Get("sex"))
	//	time.Sleep(time.Second)
	//}

	for {
		fmt.Println(cache.Len())
		time.Sleep(time.Second * 1)
	}
}
