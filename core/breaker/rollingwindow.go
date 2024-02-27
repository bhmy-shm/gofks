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

	rw.updateOffset() //更新当前滑动窗口状态，算法函数。决定当前要操作的桶的偏移量。

	rw.win.add(rw.offset, v)
}

// Reduce 遍历所有桶（可选是否包含当前桶），并对它们执行一个回调函数
// 用来计算桶中的值，和其他的统计信息。
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

// 计算自上次更新桶更新（updateOffset），经过了多少时间间隔。
// 这个返回值，用来判断这些桶是否需要被重制。
func (rw *RollingWindow) span() int {
	offset := int(timex.Since(rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}

	return rw.size
}

// 更新当前的偏移量，并重置过期的桶
func (rw *RollingWindow) updateOffset() {

	//计算上次桶更新以来经历了多少时间间隔，拿到的结果就是要重制的桶的数量
	span := rw.span()
	if span <= 0 {
		//没有桶重置，直接返回。
		return
	}

	offset := rw.offset
	// 遍历所有的过期桶，调用resetBucket进行重置。保证桶内的数据始终是最新的时间间隔内的数据。
	for i := 0; i < span; i++ {
		rw.win.resetBucket((offset + i + 1) % rw.size)
	}

	// 更新滑动窗口时间段，通过偏移量更新。
	rw.offset = (offset + span) % rw.size
	now := timex.Now()

	// 更新最后一个桶的开始时间。对齐最近的时间间隔。
	rw.lastTime = now - (now-rw.lastTime)%rw.interval
}

// IgnoreCurrentBucket 是一个 RollingWindowOption 函数，用于设置 ignoreCurrent 字段，让 Reduce 方法在执行时忽略当前的桶。
func IgnoreCurrentBucket() RollingWindowOption {
	return func(w *RollingWindow) {
		w.ignoreCurrent = true
	}
}
