package metrics

import (
	"github.com/bhmy-shm/gofks/core/barrier"
	"github.com/bhmy-shm/gofks/core/utils/timex"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type (
	//TaskContainer接口定义了一个可以用作底层容器的类型，用于执行周期性任务。

	TaskContainer interface {
		// AddTask adds the task into the container. 如果添加后需要刷新，则返回true
		AddTask(task interface{}) bool
		// Execute 在刷新容器时，处理容器中收集到的任务。
		Execute(tasks interface{})
		// RemoveAll 移除容器中包含的任务，并将它们返回。
		RemoveAll() interface{}
	}
)

// PeriodicalExecutor  是一个周期性执行任务的执行器。
type PeriodicalExecutor struct {
	commander chan interface{}
	inflight  int32

	lock      sync.Mutex
	interval  time.Duration
	container TaskContainer
	newTicker func(duration time.Duration) timex.Ticker

	//可以并发执行多个server的周期任务
	guarded   bool
	wg        *sync.WaitGroup
	wgBarrier *barrier.Barrier
}

// NewPeriodicalExecutor returns a PeriodicalExecutor with given interval and container.
func NewPeriodicalExecutor(interval time.Duration, container TaskContainer) *PeriodicalExecutor {
	executor := &PeriodicalExecutor{
		// buffer 1 to let the caller go quickly 快速添加任务
		commander: make(chan interface{}, 1),
		interval:  interval,
		container: container,
		newTicker: func(d time.Duration) timex.Ticker {
			return timex.NewTicker(d)
		},
	}

	return executor
}

// Add 追加一个任务
func (pe *PeriodicalExecutor) Add(task interface{}) {
	if newPe, ok := pe.addAndCheck(task); ok {
		pe.commander <- newPe
	}

}

// Flush 强制执行任务
func (pe *PeriodicalExecutor) Flush() bool {

	//开启新的任务执行的生命周期
	pe.enterExecution()

	//执行任务并结束生命周期
	return pe.executeTasks(func() interface{} {
		pe.lock.Lock()
		defer pe.lock.Unlock()

		//结束后清空任务
		return pe.container.RemoveAll()
	})
}

// ==================

// 追加任务并校验当前周期内是否还剩余任务，如果无剩余则直接刷新
func (pe *PeriodicalExecutor) addAndCheck(task interface{}) (interface{}, bool) {
	pe.lock.Lock()
	defer func() {
		//TODO 在追加任务时要判断当前周期是否没有任务了，没有了就刷新一下
		if !pe.guarded {
			pe.guarded = true
			// defer to unlock quickly
			defer pe.backgroundFlush()
		}
		pe.lock.Unlock()
	}()

	if pe.container.AddTask(task) {
		atomic.AddInt32(&pe.inflight, 1) //原子性的增加1个任务记录
		return pe.container.RemoveAll(), true
	}
	return nil, false
}

// 如果包含任务，则执行任务，执行任务结束后释放goroutine
func (pe *PeriodicalExecutor) executeTasks(task interface{}) bool {
	defer pe.doneExecution()

	ok := pe.hasTasks(task)
	if ok {
		pe.container.Execute(task)
	}

	return ok
}

// 判断task是否包含任务
func (pe *PeriodicalExecutor) hasTasks(task interface{}) bool {
	if task == nil {
		return false
	}

	val := reflect.ValueOf(task)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() > 0
	default:
		// 默认返回true，需要调用者自己评定
		return true
	}
}

// 添加一个新的协程生命周期，用来执行任务
func (pe *PeriodicalExecutor) enterExecution() {
	pe.wgBarrier.Guard(func() {
		pe.wg.Add(1)
	})
}

// 结束添加的协程生命周期，用来关闭任务
func (pe *PeriodicalExecutor) doneExecution() {
	pe.wg.Done()
}

// 判断是否可以停止周期任务
func (pe *PeriodicalExecutor) quitExecution(last time.Duration) (stop bool) {

	//判断任务执行是否在指定的超时时间范围内，如果没有达到阈值，代表还能继续执行，那就让它再跑一会
	if timex.Since(last) <= pe.interval*10 {
		return
	}

	pe.lock.Lock()
	if atomic.LoadInt32(&pe.inflight) == 0 {
		pe.guarded = false
		stop = true
	}
	pe.lock.Unlock()

	return
}

// 生命周期等待任务执行，或定时器过期刷新周期
func (pe *PeriodicalExecutor) backgroundFlush() {
	barrier.GoSafe(func() {
		// 在退出前刷新
		defer pe.Flush()

		ticker := pe.newTicker(pe.interval)
		defer ticker.Stop()

		var commanded bool
		last := timex.Now()
		for {
			select {

			case vals := <-pe.commander:

				commanded = true
				atomic.AddInt32(&pe.inflight, -1)
				pe.enterExecution()
				pe.executeTasks(vals)
				last = timex.Now()

			case <-ticker.Chan():

				if commanded {
					commanded = false
				} else if pe.Flush() {
					last = timex.Now()
				} else if pe.quitExecution(last) {
					return
				}

			}
		}
	})
}
