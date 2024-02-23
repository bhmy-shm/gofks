package breaker

import (
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"sync"
	"time"
)

type (
	// RollingWindowOption let callers customize the RollingWindow.
	RollingWindowOption func(rollingWindow *RollingWindow)

	// RollingWindow defines a rolling window to calculate the events in buckets with time interval.
	RollingWindow struct {
		lock          sync.RWMutex
		size          int           //窗口大小，即桶的数量
		win           *window       //包含所有桶
		interval      time.Duration //代表每个桶所表的时间间隔
		offset        int           //当前偏移量，当前是哪个桶
		ignoreCurrent bool          //决定是否在统计时，忽略当前桶
		lastTime      time.Duration //上一次桶开始的时间
	}
)

// NewRollingWindow 创建一个新的 RollingWindow 实例。它接受窗口大小和时间间隔作为参数，
// 并可通过 RollingWindowOption 函数来自定义实例。
func NewRollingWindow(size int, interval time.Duration, opts ...RollingWindowOption) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}

	w := &RollingWindow{
		size:     size,
		win:      newWindow(size),
		interval: interval,
		lastTime: timex.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Add  向当前桶中添加一个值
func (rw *RollingWindow) Add(v float64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.win.add(rw.offset, v)
}

// Reduce 遍历所有桶（可选是否包含当前桶），并对它们执行一个回调函数
func (rw *RollingWindow) Reduce(fn func(b *Bucket)) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	// ignore current bucket, because of partial data
	if span == 0 && rw.ignoreCurrent {
		diff = rw.size - 1
	} else {
		diff = rw.size - span
	}
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		rw.win.reduce(offset, diff, fn)
	}
}

// 计算自上次更新以来经过了多少时间间隔。
func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}

	return rw.size
}

// 更新当前的偏移量，并重置过期的桶
func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}

	offset := rw.offset
	// reset expired buckets
	for i := 0; i < span; i++ {
		rw.win.resetBucket((offset + i + 1) % rw.size)
	}

	rw.offset = (offset + span) % rw.size
	now := timex.Now()
	// align to interval time boundary
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}

// Bucket defines the bucket that holds sum and num of additions.
type Bucket struct {
	Sum   float64 //存储桶中所有值的总和。
	Count int64   //存储桶中值的数量
}

// 向桶中添加一个值，并更新计数
func (b *Bucket) add(v float64) {
	b.Sum += v
	b.Count++
}

// 重置桶的状态
func (b *Bucket) reset() {
	b.Sum = 0
	b.Count = 0
}

type window struct {
	buckets []*Bucket // 存储所有桶的切片
	size    int       //窗口大小，桶的数量
}

func newWindow(size int) *window {
	buckets := make([]*Bucket, size)
	for i := 0; i < size; i++ {
		buckets[i] = new(Bucket)
	}
	return &window{
		buckets: buckets,
		size:    size,
	}
}

// 在特定的偏移量处向桶中添加值
func (w *window) add(offset int, v float64) {
	w.buckets[offset%w.size].add(v)
}

// 从特定的起始点开始，对指定数量的桶执行一个函数。
func (w *window) reduce(start, count int, fn func(b *Bucket)) {
	for i := 0; i < count; i++ {
		fn(w.buckets[(start+i)%w.size])
	}
}

// 重置指定偏移量处的桶。
func (w *window) resetBucket(offset int) {
	w.buckets[offset%w.size].reset()
}

// IgnoreCurrentBucket 是一个 RollingWindowOption 函数，用于设置 ignoreCurrent 字段，让 Reduce 方法在执行时忽略当前的桶。
func IgnoreCurrentBucket() RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = true
	}
}
