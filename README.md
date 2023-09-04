# puer
go协程工具库
目前提供了两种工具类
```
go get github.com/guanyaowen/puer
```

## 1.常驻协程池
使用方法：
```
// 创建一个10核心常驻协程，任务队列容量100的协程池
pool := goworker.NewTaskPool(10, 100)
// 启动协程池，一般是在服务启动时
pool.Run()

// 模拟业务逻辑中的后台异步任务
tasks := make([]goworker.Task, 0, 50)
for i := 0; i < 50; i++ {
	i := i
	// 创建要执行的异步任务
	task := goworker.NewTask(context.Background(), func() error {
		fmt.Println(i)

		if i == 30 {
			return errors.New("error xxx")
		}
		return nil
	})
	// 推送到协程池
	pool.PushTask(task)
	tasks = append(tasks, task)
}

// 等待这批任务执行完，如果是纯异步任务，无需阻塞等待任务执行结果的情况
// 可以不用执行task.Wait()等待任务执行完毕
for _, task := range tasks {
	if err := task.Wait(); err != nil {
		fmt.Println(err)
	}
}

// 关闭协程池，一般是在服务退出时
pool.Stop()
```

## 2.并发协程任务组
使用方法：
```
// 创建一个协程并发任务组
g := &TaskGroup{
	AllowSomeFail: true, // 是否允许部分失败
	Concurrent:    50,   // 调整并发度, 默认值是10
}
start := time.Now()
for i := 0; i < 50; i++ {
	// 启动50个协程
	g.Go(func() error {
		time.Sleep(time.Second)
		return nil
	})
}

// 等待任务组的协程全部执行结束
_ = g.Wait()
fmt.Println("任务耗时：", time.Since(start))
```
