package metrics

import (
	"github.com/bhmy-shm/gofks/core/barrier"
	"log"
	"os"
	"sync"
	"time"
)

/*
	定义统计模块，用于收集和展示任务执行的统计信息。代码中主要包含了以下几个部分：

	1. 定义了Writer接口，定义了一个Write方法，用于将统计报告写入到指定的地方。
	2. 定义了StatReport结构体，表示一个统计报告的条目，包含了任务名称、时间戳、进程ID、每秒请求数、丢弃数、平均耗时、中位数耗时、90%耗时、99%耗时和99.9%耗时等信息。
	3. 定义了Metrics结构体，用于记录和展示统计信息。其中包含了一个定时执行器executor和一个统计信息容器container。

	提供了一些公共方法，如DisableLog用于禁用日志、SetReportWriter用于设置统计报告的写入器、NewMetrics用于创建一个新的Metrics实例、Add用于添加任务、AddDrop用于添加丢弃的任务、SetName用于设置名称等。

	1. 定义了metricsContainer结构体，用于保存统计信息的容器。其中包含了名称、进程ID、任务列表、执行时间和丢弃数等字段。
	2. 实现了metricsContainer的AddTask方法和Execute方法，用于添加任务和执行统计。
	3. 实现了metricsContainer的RemoveAll方法，用于获取并清空统计信息。
	4. 实现了log方法，用于打印统计报告和写入统计报告。

	提供了一个简单的统计模块，可以用于记录和展示任务执行的统计信息。
*/

var (
	logInterval  = time.Minute
	logEnabled   = true
	writerLock   sync.Mutex
	reportWriter Writer = nil
)

// StatReport is a stat report entry.
type (
	Writer interface {
		Write(report *StatReport) error
	}

	StatReport struct {
		Name          string  `json:"name" description:"任务名称"`
		Timestamp     int64   `json:"tm" description:"统计报告的时间戳，表示报告生成的时间"`
		Pid           int     `json:"pid" description:"进程的ID，用于标识生成报告的进程"`
		ReqsPerSecond float32 `json:"qps" description:"每秒请求数，表示任务的平均每秒请求数量。"`
		Drops         int     `json:"drops"` //Drops：丢弃数，表示任务被丢弃的次数。
		Average       float32 `json:"avg"`   //Average：平均耗时，表示任务的平均执行时间。
		Median        float32 `json:"med"`   //Median：中位数耗时，表示任务执行时间的中间值，即将任务按照执行时间排序后，位于中间位置的任务的执行时间。
		Top90th       float32 `json:"t90"`   //Top90th：90% 耗时，表示任务执行时间的第 90 百分位数，即将任务按照执行时间排序后，位于前 90% 的任务的执行时间。
		Top99th       float32 `json:"t99"`   //Top99th：99% 耗时，表示任务执行时间的第 99 百分位数，即将任务按照执行时间排序后，位于前 99% 的任务的执行时间。
		Top99p9th     float32 `json:"t99p9"` //Top99p9th：99.9% 耗时，表示任务执行时间的第 99.9 百分位数，即将任务按照执行时间排序后，位于前 99.9% 的任务的执行时间。
	}

	// A Metrics is used to log and report stat reports.
	Metrics struct {
		executor  *PeriodicalExecutor
		container *metricsContainer
	}
)

// NewMetrics returns a Metrics.
func NewMetrics(name string) *Metrics {
	container := &metricsContainer{
		name: name,
		pid:  os.Getpid(),
	}

	return &Metrics{
		executor:  NewPeriodicalExecutor(logInterval, container),
		container: container,
	}
}

func (m *Metrics) Add(task barrier.Task) {
	m.executor.Add(task)
}

type (
	tasksDurationPair struct {
		tasks    []barrier.Task
		duration time.Duration
		drops    int
	}
	metricsContainer struct {
		name     string
		pid      int
		tasks    []barrier.Task
		duration time.Duration
		drops    int
	}
)

func (c *metricsContainer) AddTask(v any) bool {
	if task, ok := v.(barrier.Task); ok {
		if task.Drop {
			c.drops++
		} else {
			c.tasks = append(c.tasks, task)
			c.duration += task.Duration
		}
	}
	return false
}

func (c *metricsContainer) Execute(v any) {

	pair := v.(tasksDurationPair)
	tasks := pair.tasks
	duration := pair.duration
	drops := pair.drops
	size := len(tasks)
	report := &StatReport{
		Name:          c.name,
		Timestamp:     time.Now().Unix(),
		Pid:           c.pid,
		ReqsPerSecond: float32(size) / float32(logInterval/time.Second),
		Drops:         drops,
	}

	if size > 0 {
		report.Average = float32(duration/time.Millisecond) / float32(size)

		fiftyPercent := size >> 1
		if fiftyPercent > 0 {
			top50pTasks := barrier.TopK(tasks, fiftyPercent)
			medianTask := top50pTasks[0]
			report.Median = float32(medianTask.Duration) / float32(time.Millisecond)
			tenPercent := fiftyPercent / 5
			if tenPercent > 0 {
				top10pTasks := barrier.TopK(top50pTasks, tenPercent)
				task90th := top10pTasks[0]
				report.Top90th = float32(task90th.Duration) / float32(time.Millisecond)
				onePercent := tenPercent / 10
				if onePercent > 0 {
					top1pTasks := barrier.TopK(top10pTasks, onePercent)
					task99th := top1pTasks[0]
					report.Top99th = float32(task99th.Duration) / float32(time.Millisecond)
					pointOnePercent := onePercent / 10
					if pointOnePercent > 0 {
						topPointOneTasks := barrier.TopK(top1pTasks, pointOnePercent)
						task99Point9th := topPointOneTasks[0]
						report.Top99p9th = float32(task99Point9th.Duration) / float32(time.Millisecond)
					} else {
						report.Top99p9th = barrier.GetTopDuration(top1pTasks)
					}
				} else {
					mostDuration := barrier.GetTopDuration(top10pTasks)
					report.Top99th = mostDuration
					report.Top99p9th = mostDuration
				}
			} else {
				mostDuration := barrier.GetTopDuration(top50pTasks)
				report.Top90th = mostDuration
				report.Top99th = mostDuration
				report.Top99p9th = mostDuration
			}
		} else {
			mostDuration := barrier.GetTopDuration(tasks)
			report.Median = mostDuration
			report.Top90th = mostDuration
			report.Top99th = mostDuration
			report.Top99p9th = mostDuration
		}
	}

	logPrint(report)
}

func (c *metricsContainer) RemoveAll() interface{} {
	tasks := c.tasks
	duration := c.duration
	drops := c.drops
	c.tasks = nil
	c.duration = 0
	c.drops = 0

	return tasksDurationPair{
		tasks:    tasks,
		duration: duration,
		drops:    drops,
	}
}

func writeReport(report *StatReport) {
	writerLock.Lock()
	defer writerLock.Unlock()

	if reportWriter != nil {
		if err := reportWriter.Write(report); err != nil {
			//logx.Error(err)
		}
	}
}

func logPrint(report *StatReport) {
	writeReport(report)
	if logEnabled {
		log.Printf("(%s) - qps: %.1f/s, drops: %d, avg time: %.1fms, med: %.1fms, "+
			"90th: %.1fms, 99th: %.1fms, 99.9th: %.1fms",
			report.Name, report.ReqsPerSecond, report.Drops, report.Average, report.Median,
			report.Top90th, report.Top99th, report.Top99p9th)
	}
}
