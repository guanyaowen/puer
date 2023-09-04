package goworker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/guanyaowen/puer/util"
)

func TestNewTaskPool(t *testing.T) {
	tp := NewTaskPool(5, 100)
	tp.Run()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	for i := 0; i < 40; i++ {
		i := i
		ta := NewTask(ctx, func() error {
			fmt.Println(i)
			return nil
		})
		tp.PushTask(ta)
	}
	tp.Stop()
}

func TestNewTaskPoolErr(t *testing.T) {
	tp := NewTaskPool(5, 100)
	tp.Run()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	tasks := make([]Task, 0, 10)
	for i := 0; i < 5; i++ {
		i := i
		ta := NewTask(ctx, func() error {
			fmt.Println(i)
			return errors.New("error")
		})
		tp.PushTask(ta)
		tasks = append(tasks, ta)
	}

	for _, t := range tasks {
		if err := t.Wait(); err != nil {
			fmt.Printf("err:%v \n", err)
		}
	}

	tp.Stop()
}

func TestNewTaskPoolPanic(t *testing.T) {
	tp := NewTaskPool(5, 100)
	tp.Run()

	func() {
		defer util.MeasureFuncExecTime()()
		tasks := make([]Task, 0, 10)
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		for i := 0; i < 1; i++ {
			ta := NewTask(ctx, func() error {
				panic("panic xxxx")
				return nil
			})
			tp.PushTask(ta)
			tasks = append(tasks, ta)
		}

		for _, t := range tasks {
			if err := t.Wait(); err != nil {
				fmt.Printf("err:%v \n", err)
			}
		}
	}()

	tp.Stop()
}

func TestNewTaskPool_addWorker(t *testing.T) {
	tp := NewTaskPool(10, 100)
	tp.Run()

	func() {
		defer util.MeasureFuncExecTime()()
		tasks := make([]Task, 0, 100)
		for i := 0; i < 100; i++ {
			ta := NewTask(context.Background(), func() error {
				time.Sleep(1 * time.Second)
				return nil
			})
			tp.PushTask(ta)
			tasks = append(tasks, ta)
		}

		for _, t := range tasks {
			if err := t.Wait(); err != nil {
				fmt.Printf("err:%v \n", err)
			}
		}
	}()

	tp.Stop()
}
