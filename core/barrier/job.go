package barrier

import "sync"

/*
	轻便的 goroutine 封装，快速编写类对象的业务代码，sync.WaitGroup的优雅实现
	需要在进程中阻塞等待封装内的业务执行结束才能跳出
*/

type JobFunc func() interface{}

type job struct {
	jobs []JobFunc
	data chan interface{}
	wg   *sync.WaitGroup
}

func NewJob() *job {
	return &job{
		data: make(chan interface{}),
		wg:   &sync.WaitGroup{},
	}
}

func (j *job) Set(jobs ...JobFunc) {
	if j == nil {
		return
	}
	j.jobs = append(j.jobs, jobs...)
}

func (j *job) do() {
	if j == nil {
		return
	}
	for _, fn := range j.jobs {
		j.wg.Add(1)
		go func(f JobFunc) {
			defer j.wg.Done()
			j.data <- f
		}(fn)
	}
}

func (j *job) Range(f func(err interface{})) {
	if j == nil {
		return
	}
	j.do()
	go func() {
		defer close(j.data)
		j.wg.Wait()
	}()

	for v := range j.data {
		f(v)
	}
}
