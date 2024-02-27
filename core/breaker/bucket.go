package breaker

/*
	在滑动窗口算法中，桶，是用来收集和维护一段时间内的统计信息的数据结构。监控和计算过去一段时间内请求的成功和失败的次数。
	滑动窗口通常由多个桶组成，每个桶代表窗口时间中的一个子时间段。
	例如，如果您有一个10秒的窗口，并且把它分成了40个桶，那么每个桶就代表着250毫秒的时间段。

	具体包括：
	1. 时间段细分：每个桶对应窗口中的一个固定时间段，负责收集该时间段内的请求数据。
	2. 滑动特性：根据时间的推移，最旧的桶会被新桶替换，保持窗口的总时长不变。
	3. 数据聚合：当需要计算过去一段时间内的统计数据时，所有桶中的数据会被聚合起来，得到计算结果。
	4. 熔断器决策：通过分析成功率 和 失败率，可以决定熔断器当前是否应该放行，开闭状态。
*/

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
