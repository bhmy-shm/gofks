package mathx

import (
	"math/rand"
	"sync"
	"time"
)

/*
	解决问题，避免缓存过期的问题：

1. 使用不稳定的过期时间：通过使用Unstable结构体生成的随机值作为缓存项的过期时间，可以使过期时间不稳定化，避免大量缓存项在同一时间过期。

2. 偏差值：通过在NewUnstable函数中传入一个偏差值来控制随机值的范围。较大的偏差值将导致生成的随机值更加分散，进一步减少缓存项在同一时间过期的可能性。

3. 使用互斥锁：为了保证并发安全，在AroundDuration和AroundInt方法中使用了互斥锁来保护随机值的生成。这样可以防止多个goroutine同时访问并修改随机数生成器，确保生成的随机值是唯一的

*/

// A Unstable 生成围绕给定偏差的均值的随机值
type Unstable struct {
	deviation float64     //偏差参数，表示生成的随机值相对于基准值的偏移程度
	r         *rand.Rand  //rand 随机数生成器
	lock      *sync.Mutex //互斥保护，多线程环境下生成器的安全访问
}

// NewUnstable returns a Unstable.
func NewUnstable(deviation float64) Unstable {
	if deviation < 0 {
		deviation = 0
	}
	if deviation > 1 {
		deviation = 1
	}
	return Unstable{
		deviation: deviation,
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:      new(sync.Mutex),
	}
}

// AroundDuration 生成一个围绕基准时间 base 的随机时间间隔。
func (u Unstable) AroundDuration(base time.Duration) time.Duration {
	u.lock.Lock()
	val := time.Duration((1 + u.deviation - 2*u.deviation*u.r.Float64()) * float64(base))
	u.lock.Unlock()
	return val
}

// AroundInt 生成一个围绕基准整数 base 的随机整数
func (u Unstable) AroundInt(base int64) int64 {
	u.lock.Lock()
	val := int64((1 + u.deviation - 2*u.deviation*u.r.Float64()) * float64(base))
	u.lock.Unlock()
	return val
}
