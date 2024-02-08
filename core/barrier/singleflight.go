package barrier

import "sync"

/*

SingleFlight 是一个用于避免重复执行相同关键字的函数调用的机制。
它的目的是在并发环境下，通过共享调用结果来提高性能，避免重复的计算或查询工作。

详细说明：
	在并发环境下，当多个协程或线程同时调用同一个函数，并且传入相同的参数时,
	如果这个函数是一个耗时的操作（例如数据库查询、网络请求等），那么就会造成重复的工作和资源浪费。
	SingleFlight 机制通过对这种情况的判断和处理，使得只有一个调用会真正执行函数，而其他调用会等待并共享这个调用的结果。

实现思路：
	SingleFlight 的核心思想是通过一个中间层来管理并发调用，并跟踪相同关键字的调用状态。
	当发现有相同关键字的调用时，后续的调用会等待并共享之前调用的结果，而不会再次执行相同的函数。
	这样可以避免重复的工作，并提高并发性能。

使用场景：
	1. 并发场景下的读取缓存操作，多个线程间同时读取相同key的操作，可以将key放到中间层来管理调度。
*/

type (
	//SingleFlight 并发任务抽象
	SingleFlight interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
		DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
	}

	//并发任务实现
	call struct {
		wg  sync.WaitGroup
		val interface{}
		err error
	}

	//并发任务中央管控
	flightGroup struct {
		calls map[string]*call
		lock  sync.Mutex //并发操作互斥保护
	}
)

func NewFlightGroup() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

func (g *flightGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {

	// 添加并发任务
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	// 执行任务
	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error) {

	// 添加并发任务
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	// 执行任务
	g.makeCall(c, key, fn)
	return c.val, true, c.err
}

func (g *flightGroup) createCall(key string) (*call, bool) {
	g.lock.Lock()

	if c, found := g.calls[key]; found {
		g.lock.Unlock() //如果传入的任务已经存在，先接触互斥锁，让其他资源访问临界值
		c.wg.Wait()     //等待任务执行结束
		return c, true
	}

	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (interface{}, error)) {
	defer func() {
		//在执行完任务后，需要互斥安全的删除该并发任务
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()

		//结束线程
		c.wg.Done()
	}()

	//执行任务
	c.val, c.err = fn()
}
