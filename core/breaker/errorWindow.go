package breaker

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	numHistoryReasons = 5
	timeFormat        = "15:04:05"
)

// errorWindow 跟踪和记录最近发生的错误，
// 这段代码所实现的是一个基础的循环数组日志功能，它记录和格式化错误原因，可以说是实现熔断器模式中的“记录错误”部分

type errorWindow struct {
	reasons [numHistoryReasons]string //reason 数组，用于存储错误原因的字符串
	index   int                       //跟踪当前应该写入数组的位置
	count   int                       //记录已经添加到窗口的错误数量
	lock    sync.Mutex                //并发安全互斥锁
}

// 添加一个新的错误原因到窗口中。它首先锁定窗口，以防止并发写入时的数据竞争。
// 然后，将错误原因和当前时间的字符串表示形式添加到 reasons 数组中，更新 index 和 count，并解锁。
func (ew *errorWindow) add(reason string) {
	ew.lock.Lock()
	ew.reasons[ew.index] = fmt.Sprintf("%s %s", time.Now().Format(timeFormat), reason)
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = MinInt(ew.count+1, numHistoryReasons)
	ew.lock.Unlock()
}

// 用来生成并返回一个包含所有错误原因的字符串。它同样使用锁来保证在读取错误原因时数据的一致性
func (ew *errorWindow) String() string {
	var reasons []string

	ew.lock.Lock()
	// reverse order
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		reasons = append(reasons, ew.reasons[(i+numHistoryReasons)%numHistoryReasons])
	}
	ew.lock.Unlock()

	return strings.Join(reasons, "\n")
}

// MaxInt returns the larger one of a and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the smaller one of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

type promiseWithReason struct {
	promise internalPromise
	errWin  *errorWindow
}

func (p promiseWithReason) Accept() {
	p.promise.Accept()
}

func (p promiseWithReason) Reject(reason string) {
	p.errWin.add(reason)
	p.promise.Reject()
}
