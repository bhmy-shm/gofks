package thread

import (
	"log"
	"sort"
	"sync"
	"time"
)

/*
	多消息队列 配合 线程池 构建的 工作池模型
	需要 手动配置 消息队列数量 以及 消息队列中的工作线程数量

	适合于需要阻塞的进程任务，放到队列中执行，但不会确定结果是否成功。（待优化）
*/

const (
	MaxQueue  int = 5
	MinQueue  int = 0
	MaxWorker int = 100
	MinWorker int = 2
)

type MsgHandler struct {
	counter    int //计数器
	queueSize  int //消息队列长度
	workerSize int //线程池长度

	TaskQueue []chan *TaskExecutor //消息管道
	Quit      chan int             //退出队列
	mux       sync.Mutex
	wg        sync.WaitGroup
}

func TaskQueue(size ...int) *MsgHandler {
	if size[0] > MaxQueue || size[0] < MinQueue || len(size) > 1 {
		panic("管道设置不合法")
	}

	if len(size) == 0 {
		size[0] = MinQueue
	}

	return &MsgHandler{
		queueSize: size[0],
		TaskQueue: make([]chan *TaskExecutor, size[0]),
		Quit:      make(chan int),
		counter:   -1,
	}
}

func (this *MsgHandler) queueLength() int {
	return len(this.TaskQueue)
}

func (m *MsgHandler) Close() {
	for _, queue := range m.TaskQueue {
		time.Sleep(time.Millisecond * 500)
		close(queue)
	}
	m.Quit <- 1
	defer close(m.Quit)
	m.TaskQueue = nil
}

//开启阻塞线程，执行工作任务

func (m *MsgHandler) StartWorkerPool(size int) *MsgHandler {
	if size > MaxWorker || size < MinWorker {
		panic("工作池设置不合法")
	}

	for i := 0; i < m.queueLength(); i++ {
		m.TaskQueue[i] = make(chan *TaskExecutor, size) //创建channel用来接收任务
	}

	go m.startWorker()
	m.workerSize = size
	return m
}

func (m *MsgHandler) startWorker() {
	var key = false

	go func() {
		select {
		case q := <-m.Quit:
			if q > 0 {
				key = true
			}
		}
	}()

	for i := 0; i < len(m.TaskQueue); i++ {
		m.wg.Add(1)
		go func(task chan *TaskExecutor) {
			defer m.wg.Done()
			for key == false {
				select {
				case req := <-task:
					//Todo 任务执行的成功或失败没有做处理，没有做保底策略
					req.do()
				}
			}
		}(m.TaskQueue[i])
	}
	m.wg.Wait()
	log.Println("The thread task is finished and be exited")
	return
}

// 轮询

func (m *MsgHandler) SendMsgRobin(executor ...*TaskExecutor) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, exec := range executor {
		msg := m.roundRobin()
		msg <- exec
	}
}

func (m *MsgHandler) roundRobin() chan *TaskExecutor {
	m.counter++

	if m.counter >= len(m.TaskQueue) {
		m.counter = 0
	}
	return m.TaskQueue[m.counter]
}

//根据队列中 已经执行的线程任务量，返回一个线程任务处于中间以下的chan

func (m *MsgHandler) SendMsgLoad(executor ...*TaskExecutor) {
	m.mux.Lock()
	defer m.mux.Unlock()

	for _, exec := range executor {
		msg := m.getMsgLoad()
		if msg == nil {
			//保底策略找一个目前队列线程数量最少的，赛里面
			msg = m.getMinMsgLoad()
			log.Println("保底策略")
		}
		msg <- exec
	}
}

func (m *MsgHandler) getMsgLoad() chan *TaskExecutor {

	//遍历每一个管道
	for i := 0; i < len(m.TaskQueue); i++ {

		mid := m.workerSize / 2

		//获取当前管道执行的线程数量
		workerSize := len(m.TaskQueue[i])

		//如果当前管道未开启，或者工作队列中线程任务是满的，则直接跳过
		if m.TaskQueue[i] == nil || workerSize == m.workerSize {
			continue
		}

		//如果当前线程池一个线程也没有，则直接返回
		if m.TaskQueue[i] != nil && workerSize == 0 {
			return m.TaskQueue[i]
		}

		//如果当前工作的线程数量没有超过中位数，则写入该线程
		if m.TaskQueue[i] != nil && workerSize > 0 {
			if workerSize <= mid {
				return m.TaskQueue[i]
			} else if workerSize > mid {
				continue
			}
		}
	}
	return nil
}

func (m *MsgHandler) getMinMsgLoad() chan *TaskExecutor {
	sort.Slice(m.TaskQueue, func(i, j int) bool {
		return len(m.TaskQueue[i]) < len(m.TaskQueue[j])
	})
	for i := 0; i < len(m.TaskQueue); i++ {
		if len(m.TaskQueue[i]) >= 0 && len(m.TaskQueue[i]) < 5 {
			return m.TaskQueue[i]
		}
	}
	return m.TaskQueue[0]
}
