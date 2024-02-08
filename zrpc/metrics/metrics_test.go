package metrics

//func runExecutors() {
//
//	container := &MyTaskContainer{
//		name: "testExecutors",
//		pid:  os.Getpid(),
//	}
//
//	executor := executors.NewPeriodicalExecutor(time.Second, container)
//	defer executor.Wait()
//
//	// 添加任务到执行器
//	executor.Add("Task 1")
//	executor.Add("Task 2")
//	executor.Add("Task 3")
//
//	// 等待一段时间，观察任务是否被周期性执行
//	time.Sleep(5 * time.Second)
//}
//
//func TestPeriodExecutors(t *testing.T) {
//	runExecutors()
//}
