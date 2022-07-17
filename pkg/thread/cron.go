package thread

import (
	"github.com/robfig/cron/v3"
	"sync"
)

var (
	taskCron *cron.Cron     //定时任务
	onceCron sync.Once
)

func GetCronTask() *cron.Cron {
	onceCron.Do(func() {
		taskCron = cron.New(cron.WithSeconds())
	})
	return taskCron
}
