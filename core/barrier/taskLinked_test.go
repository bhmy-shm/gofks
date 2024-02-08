package barrier

import (
	"fmt"
	"testing"
)

func current(params ...interface{}) {
	fmt.Println("current", params)
}

func prev(params ...interface{}) {
	fmt.Println("prev", params)
}

func after(params ...interface{}) {
	fmt.Println("after", params)
}

func Test_topic_task(t *testing.T) {

	task := Topic()

	task.AddTopic("no1", WithCurrent(current, 1, 1, 1)).
		PrevHandle("no1", WithPrev(prev, 2, 2, 2)).
		AfterHandle("no1", WithAfter(after, 3, 3, 3))

	task.AddTopic("no2", WithCurrent(current, 10, 10, 10)).
		PrevHandle("no2", WithPrev(prev, 20, 20, 20)).
		AfterHandle("no2", WithAfter(after, 30, 30, 30))

	task.DoAll()

}
