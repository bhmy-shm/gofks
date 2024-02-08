package barrier

import "sync"

// 简单互斥锁实现，用来保护执行函数的唯一性

type Barrier struct {
	lock sync.Mutex
}

// Guard 通过互斥锁保护会被共享的资源
func (b *Barrier) Guard(fn func()) {
	Guard(&b.lock, fn)
}

// Guard 通过互斥锁保护会被共享的资源
func Guard(lock sync.Locker, fn func()) {
	lock.Lock()
	defer lock.Unlock()
	fn()
}
