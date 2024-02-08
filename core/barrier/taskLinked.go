package barrier

import (
	"sync"
)

type (
	HookOption struct {
		Current TaskHook
		Prev    TaskHook
		After   TaskHook
		Params  []interface{}
	}
	TaskHook func(...interface{})
)

func WithCurrent(current TaskHook, params ...interface{}) HookOption {
	return HookOption{
		Current: current,
		Params:  params,
	}
}

func WithPrev(prev TaskHook, params ...interface{}) HookOption {
	return HookOption{
		Prev:   prev,
		Params: params,
	}
}

func WithAfter(after TaskHook, params ...interface{}) HookOption {
	return HookOption{
		After:  after,
		Params: params,
	}
}

type (
	TopicInter interface {
		AddTopic(string, HookOption) *topicTask
		PrevHandle(string, HookOption) *topicTask
		AfterHandle(string, HookOption) *topicTask
		DoAll()
	}
	topicTask struct {
		tasks sync.Map
		wg    sync.WaitGroup
	}
)

func Topic() TopicInter {
	return &topicTask{}
}

func (t *topicTask) AddTopic(topicId string, opt HookOption) *topicTask {

	executor := newExecutorNode(topicId).addTask(opt)
	t.tasks.Store(executor.getTopic(), executor)
	return t
}

func (t *topicTask) PrevHandle(topicId string, opt HookOption) *topicTask {
	executor, ok := t.tasks.Load(topicId)
	if !ok {
		newExecutor := newExecutorNode(topicId)
		t.tasks.Store(newExecutor.getTopic(), newExecutor)
	} else {
		executor.(*taskExecutor).prevHandle(opt)
	}
	return t
}

func (t *topicTask) AfterHandle(topicId string, opt HookOption) *topicTask {
	executor, ok := t.tasks.Load(topicId)
	if !ok {
		newExecutor := newExecutorNode(topicId)
		t.tasks.Store(newExecutor.getTopic(), newExecutor)
	} else {
		executor.(*taskExecutor).afterHandle(opt)
	}
	return t
}

func (t *topicTask) DoAll() {

	t.tasks.Range(func(key, value any) bool {

		executor := value.(*taskExecutor)

		t.wg.Add(1)

		go func(executor *taskExecutor) {
			defer t.wg.Done()
			executor.do()
		}(executor)

		return true
	})
	t.wg.Wait()
}

type taskExecutor struct {
	topic string
	fn    TaskHook
	param []interface{}
	prev  *taskExecutor
	after *taskExecutor
}

func newExecutorNode(topic string) *taskExecutor {
	return &taskExecutor{
		topic: topic,
	}
}

func (this *taskExecutor) addTask(f HookOption) *taskExecutor {
	newNode := &taskExecutor{
		topic: this.topic,
		fn:    f.Current,
		param: f.Params,
	}

	//如果有下一个节点
	if this.after != nil {
		next := this.after
		this.after = newNode
		newNode.prev = this
		newNode.after = next
		next.prev = newNode
	} else {
		this.after = newNode
		newNode.prev = this
	}
	return newNode
}

func (this *taskExecutor) prevHandle(f HookOption) {
	newNode := &taskExecutor{
		topic: this.topic,
		fn:    f.Prev,
		param: f.Params,
	}
	if this.prev != nil {
		prev := this.prev
		this.prev = newNode
		newNode.prev = prev
		newNode.after = this
		prev.after = newNode
	} else {
		this.prev = newNode
		newNode.after = this
	}
}

func (this *taskExecutor) afterHandle(f HookOption) {
	newExecutor := &taskExecutor{
		topic: this.topic,
		fn:    f.After,
		param: f.Params,
	}

	if this.after != nil {
		next := this.after
		this.after = newExecutor
		newExecutor.prev = this
		newExecutor.after = next
		next.prev = newExecutor
	} else {
		this.after = newExecutor
		newExecutor.prev = this
	}
}

func (this *taskExecutor) do() {

	var (
		current = this
	)
	defer this.clear()

	for current != nil {

		if current.prev != nil {
			if hook := current.prev.fn; hook != nil {
				hook(current.prev.param...)
				current.prev.clear()
			}
		}

		if current.fn != nil {
			if hook := current.fn; hook != nil {
				hook(current.param...)
				current.clear()
			}
		}

		current = current.after
	}
}

func (this *taskExecutor) getTopic() string {
	return this.topic
}

func (this *taskExecutor) clear() {
	this.fn = nil
	this.param = nil
}
