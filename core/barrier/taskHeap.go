package barrier

import (
	"container/heap"
	"time"
)

type Task struct {
	Drop        bool
	Duration    time.Duration
	Description string
}

type taskHeap []Task

func (h *taskHeap) Len() int {
	return len(*h)
}

func (h *taskHeap) Less(i, j int) bool {
	return (*h)[i].Duration < (*h)[j].Duration
}

func (h *taskHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *taskHeap) Push(x any) {
	*h = append(*h, x.(Task))
}

func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TopK 返回包含耗时最高的k个任务的切片
func TopK(all []Task, k int) []Task {
	h := new(taskHeap)
	heap.Init(h)

	for _, each := range all {
		if h.Len() < k {
			heap.Push(h, each)
		} else if (*h)[0].Duration < each.Duration {
			heap.Pop(h)
			heap.Push(h, each)
		}
	}

	return *h
}

// GetTopDuration 获取任务列表中最高耗时任务及其耗时的功能
func GetTopDuration(tasks []Task) float32 {
	top := TopK(tasks, 1)
	if len(top) < 1 {
		return 0
	}

	return float32(top[0].Duration) / float32(time.Millisecond)
}
