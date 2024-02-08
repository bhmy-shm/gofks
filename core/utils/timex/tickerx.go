package timex

import (
	"errors"
	"time"
)

type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

var errTimeout = errors.New("ticker-timeout")

type realTicker struct {
	*time.Ticker
}

func NewTicker(duration time.Duration) Ticker {
	return &realTicker{
		Ticker: time.NewTicker(duration),
	}
}

// Chan 返回只读通道，用于接收定时器触发的时间信号
func (rt *realTicker) Chan() <-chan time.Time {
	return rt.C
}

// ========================================================

type (
	FakeTickerInter interface {
		Ticker
		Done()
		Tick()
		Wait(d time.Duration) error
	}

	fakeTicker struct {
		c    chan time.Time
		done chan struct{}
	}
)

func NewFakeTicker() FakeTickerInter {
	return &fakeTicker{
		c:    make(chan time.Time, 1),
		done: make(chan struct{}, 1),
	}
}

func (ft *fakeTicker) Chan() <-chan time.Time {
	return ft.c
}

func (ft *fakeTicker) Done() {
	ft.done <- struct{}{}
}

func (ft *fakeTicker) Stop() {
	close(ft.c)
}

func (ft *fakeTicker) Tick() {
	ft.c <- time.Now()
}

func (ft *fakeTicker) Wait(d time.Duration) error {
	select {
	case <-time.After(d):
		return errTimeout
	case <-ft.done:
		return nil
	}
}
