package gofks

import (
	"bytes"
	"github.com/bhmy-shm/gofks/core/cache/buffer"
	"github.com/bhmy-shm/gofks/core/logx"
	"github.com/robfig/cron/v3"
	"log"
	"sync"
)

type (
	CronInter interface {
		Get() *cron.Cron
		Close()
	}
	cronTask struct {
		taskCron *cron.Cron //定时任务
		onceCron sync.Once
	}
)

func NewCronTask() CronInter {
	return &cronTask{
		taskCron: nil,
	}
}

func (ct *cronTask) Get() *cron.Cron {
	ct.onceCron.Do(func() {
		ct.taskCron = cron.New(cron.WithSeconds(), cron.WithLogger(bufferLog()))
	})
	return ct.taskCron
}

func (ct *cronTask) Close() {
	ct.taskCron.Stop()
}

func bufferLog() cron.Logger {
	sw := &buffer.SyncWriter{}
	return cron.PrintfLogger(log.New(sw, "gofks-cron:", log.LstdFlags))
}

func myWriteLog() cron.Logger {
	var buf bytes.Buffer
	sw := logx.NewWriter(&buf)
	return cron.PrintfLogger(log.New(sw, "gofks-cron:", log.LstdFlags))

}
