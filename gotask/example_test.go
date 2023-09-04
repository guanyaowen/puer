package gotask

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func add(num *int64) {
	atomic.AddInt64(num, 1)
}

func TestExampleTaskGroup_Go(t *testing.T) {
	g := &TaskGroup{}
	var num int64

	start := time.Now()
	for i := 0; i < 100; i++ {
		g.Go(func() error {
			add(&num)
			return nil
		})
	}

	_ = g.Wait()
	fmt.Println("任务耗时：", time.Since(start))
	fmt.Println("并发累加结果：", num)
}

func sleep() {
	time.Sleep(time.Second)
}

func Test_Concurrent(t *testing.T) {
	// 调整并发度
	// 默认值是10
	g := &TaskGroup{
		Concurrent: 50,
	}
	start := time.Now()
	for i := 0; i < 50; i++ {
		g.Go(func() error {
			sleep()
			return nil
		})
	}

	_ = g.Wait()
	fmt.Println("任务耗时：", time.Since(start))
}

func Test_AllowSomeFail(t *testing.T) {
	// 允许部分任务失败
	g := TaskGroup{
		AllowSomeFail: true,
	}

	var num int64
	g.Go(func() error {
		add(&num)
		return errors.New("error1")
	})

	g.Go(func() error {
		add(&num)
		return errors.New("error2")
	})

	g.Go(func() error {
		add(&num)
		return errors.New("error3")
	})

	g.Go(func() error {
		add(&num)
		return errors.New("error4")
	})

	err := g.Wait()
	fmt.Println("首次捕捉到的err：", err)
	fmt.Println("并发累加结果：", num)
}

func Test_NotAllowSomeFail(t *testing.T) {
	// 不允许部分任务失败, 尽量阻止其他并发任务
	// AllowSomeFail 默认值就是false，可以不用设置
	g := TaskGroup{}

	var num int64
	g.Go(func() error {
		add(&num)
		return errors.New("error1")
	})

	g.Go(func() error {
		time.Sleep(time.Second * 1)
		add(&num)
		return nil
	})

	g.Go(func() error {
		time.Sleep(time.Second * 1)
		add(&num)
		return nil
	})

	g.Go(func() error {
		time.Sleep(time.Second * 1)
		add(&num)
		return nil
	})

	err := g.Wait()
	fmt.Println("首次捕捉到的err：", err)
	fmt.Println("并发累加结果：", num)
}

func Test_Panic(t *testing.T) {
	g := TaskGroup{}
	g.Go(func() error {
		panic("runtime panic")
	})
	err := g.Wait()
	fmt.Println("首次捕捉到的err：", err)
}
