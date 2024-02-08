package syncx

import (
	"sync/atomic"
	"time"
)

// 实现了一个线程安全的 AtomicDuration 类型，它是对 int64 类型的封装，用来表示 time.Duration 以保证在并发环境中的原子性操作。
// time.Duration 本身是一个 int64 类型，表示以纳秒为单位的时间长度。

type AtomicDuration int64

// NewAtomicDuration 返回一个新的 AtomicDuration 实例的指针，该实例的初始值为0
func NewAtomicDuration() *AtomicDuration {
	return new(AtomicDuration)
}

// ForAtomicDuration 返回一个新的 AtomicDuration 实例的指针，并设置初始值为给定的 time.Duration 值。
func ForAtomicDuration(val time.Duration) *AtomicDuration {
	d := NewAtomicDuration()
	d.Set(val)
	return d
}

// CompareAndSwap 执行比较并交换操作。它将当前值与 old 值进行比较，如果当前值等于 old 值，则将其设置为 val。
// 如果成功设置，返回 true；否则返回 false。
func (d *AtomicDuration) CompareAndSwap(old, val time.Duration) bool {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(old), int64(val))
}

// Load 返回当前的 AtomicDuration 值。多线程环境下原子性返回。
func (d *AtomicDuration) Load() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

// Set  将 AtomicDuration 的值设置为 val。这个操作也是原子性的，确保了值的设置不会被并发的其他操作打断。
func (d *AtomicDuration) Set(val time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(val))
}
