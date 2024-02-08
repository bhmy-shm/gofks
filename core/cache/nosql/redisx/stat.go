package redisx

import (
	"github.com/bhmy-shm/gofks/core/logx"
	"sync/atomic"
	"time"
)

const statInterval = time.Minute

// A Stat is used to stat the cache.
type Stat struct {
	name string
	// export the fields to let the unit tests working,
	// reside in internal package, doesn't matter.
	Total   uint64
	Hit     uint64
	Miss    uint64
	DbFails uint64
}

// NewStat returns a Stat.
func NewStat(name string) *Stat {
	ret := &Stat{
		name: name,
	}
	go ret.statLoop()

	return ret
}

// IncrementTotal 用于递增总请求数。每次调用该方法，总请求数会增加1
func (s *Stat) IncrementTotal() {
	atomic.AddUint64(&s.Total, 1)
}

// IncrementHit 用于递增命中数。每次调用该方法，命中数会增加1。
func (s *Stat) IncrementHit() {
	atomic.AddUint64(&s.Hit, 1)
}

// IncrementMiss 用于递增未命中数。每次调用该方法，未命中数会增加1。
func (s *Stat) IncrementMiss() {
	atomic.AddUint64(&s.Miss, 1)
}

// IncrementDbFails 递增缓存数据库失败数。每次调用该方法，数据库失败数会增加1。
func (s *Stat) IncrementDbFails() {
	atomic.AddUint64(&s.DbFails, 1)
}

func (s *Stat) statLoop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for range ticker.C {
		total := atomic.SwapUint64(&s.Total, 0)
		if total == 0 {
			continue
		}

		hit := atomic.SwapUint64(&s.Hit, 0)
		percent := 100 * float32(hit) / float32(total)
		miss := atomic.SwapUint64(&s.Miss, 0)
		dbf := atomic.SwapUint64(&s.DbFails, 0)
		logx.Statf("db-Cache(%s) - qpm: %d, hit_ratio: %.1f%%, hit: %d, miss: %d, db_fails: %d",
			s.name, total, percent, hit, miss, dbf)
	}
}
