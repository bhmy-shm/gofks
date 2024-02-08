package timex

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func TestRealTickerDoTick(t *testing.T) {

	//启动一个定时器，每秒钟监听信号并输出count
	//这个测试永远不会执行到 ticker.Stop()，因为不会跳出 chan

	ticker := NewTicker(time.Second * 1)
	defer ticker.Stop()

	var count int
	for range ticker.Chan() {
		count++
		log.Println("count:", count)
	}
}

func TestRealTickerDoTickStop(t *testing.T) {

	//启动一个定时器，每秒钟监听信号并输出count
	//假设每分钟只能执行10次，执行完10次就跳出，则执行ticker.Stop 进行回收

	ticker := NewTicker(time.Second * 1)
	defer ticker.Stop()

	var count int
	for range ticker.Chan() {
		count++
		log.Println("count:", count)
		if count == 10 {
			break
		}
	}
}

func TestFakeTickerTimeout(t *testing.T) {
	ticker := NewFakeTicker()
	defer ticker.Stop()

	err := ticker.Wait(time.Millisecond)
	if err != nil {
		log.Println("wait err:", err)
	}
}

func TestFakeTicker(t *testing.T) {
	const total = 10
	ticker := NewFakeTicker() //生成一个手动定时器
	defer ticker.Stop()       //释放定时器资源

	var count int32
	go func() {
		for range ticker.Chan() { //监听定时器信号

			//每监听一次并增加count，直到等于total，关闭定时器
			if atomic.AddInt32(&count, 1) == total {
				ticker.Done()
			}
		}
	}()

	println("wait before:", time.Now().Unix())
	err := ticker.Wait(time.Second * 2) //要求在2秒钟内完成
	if err != nil {
		println("wait after:", err.Error(), time.Now().Unix())
		return
	}
	go func() {
		//手动增加定时器信号，这里增加了5次
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 1)
			ticker.Tick()
		}
	}()

	println("end count:", count, time.Now().Unix())

}
