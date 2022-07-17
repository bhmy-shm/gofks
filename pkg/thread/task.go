package thread

import "sync"

//Task任务，topic确保唯一，包含两种回调功能

type TaskFunc func(params ...interface{}) //业务函数
type TaskPrev func()                      //回调前
type TaskAfter func()                     //回调后

var topicMap sync.Map

type TaskExecutor struct {
	topic string
	fn    TaskFunc
	prev  TaskPrev
	after TaskAfter
	param interface{}
}

func Topic(topic string) *TaskExecutor {
	if topic == "" {
		panic("Task-topic 不能为空")
	}
	if _, ok := topicMap.Load(topic); ok {
		panic("topic已存在")
	}
	return &TaskExecutor{
		topic: topic,
	}
}

func (this *TaskExecutor) AddTask(f TaskFunc, p interface{}) *TaskExecutor {
	this.fn = f
	this.param = p
	return this
}

func (this *TaskExecutor) PrevHandle(prev TaskPrev) *TaskExecutor {
	this.prev = prev
	return this
}

func (this *TaskExecutor) AfterHandle(after TaskAfter) *TaskExecutor {
	this.after = after
	return this
}

func (this *TaskExecutor) do() {
	if this.prev != nil {
		this.prev()
	}

	this.fn(this.param)

	if this.after != nil {
		this.after()
	}
}

func (this *TaskExecutor) getTopic() string {
	return this.topic
}
