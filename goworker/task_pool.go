package goworker

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/guanyaowen/puer/util"
	"github.com/guanyaowen/puer/util/maths"
)

var _ TaskPool = (*taskPool)(nil)

type TaskPool interface {
	Run()

	Stop()

	// PushTask 推送一个任务到任务池中
	PushTask(task Task)
}

type taskPool struct {
	// coreNum 核心工作协程数，只有协程池退出时核心工作协程才会退出
	coreNum int

	// taskSize 任务队列容量最大
	taskSize int
	// taskChan 任务队列
	taskChan chan Task

	// 临时工作协程退出控制
	cancelMutex sync.Mutex
	cancelFunc  []context.CancelFunc

	// 协程池退出通知
	closeChan chan struct{}

	// 工作协程控制
	wg sync.WaitGroup

	runOnce   sync.Once
	closeOnce sync.Once
}

// NewTaskPool 常驻工作池
func NewTaskPool(coreNum, taskChanSize int) TaskPool {
	return newTaskPool(coreNum, taskChanSize)
}

func newTaskPool(coreNum, taskChanSize int) (pool *taskPool) {
	taskChanSize = maths.Max(taskChanSize, 10)
	coreNum = maths.Max(coreNum, 1)
	pool = &taskPool{
		taskSize:  taskChanSize,
		coreNum:   coreNum,
		taskChan:  make(chan Task, taskChanSize),
		closeChan: make(chan struct{}),
	}
	return
}

func (p *taskPool) PushTask(task Task) {
	p.taskChan <- task
}

func (p *taskPool) Run() {
	p.runOnce.Do(func() {
		for i := 0; i < p.coreNum; i++ {
			p.coreWorker()
		}
		p.controlRoom()
	})
}

func (p *taskPool) Stop() {
	p.closeOnce.Do(func() {
		close(p.closeChan)
		close(p.taskChan)
		p.cancelMutex.Lock()
		defer p.cancelMutex.Unlock()
		for _, cancelFunc := range p.cancelFunc {
			cancelFunc()
		}
	})
	p.wg.Wait()
}

// controlRoom 资源监控
func (p *taskPool) controlRoom() {
	(&controlRoom{p: p}).Run()
}

// coreWorker 核心工作者
func (p *taskPool) coreWorker() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for t := range p.taskChan {
			p.doTask(t)
		}
	}()
}

// addWorker 新增工作者
func (p *taskPool) addWorker() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	func() {
		p.cancelMutex.Lock()
		defer p.cancelMutex.Unlock()
		p.cancelFunc = append(p.cancelFunc, cancelFunc)
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case t, ok := <-p.taskChan:
				if ok {
					p.doTask(t)
				}
			case <-ctx.Done():
				return
			case <-p.closeChan:
				return
			}
		}
	}()
}

// dismissWorker 解雇部分临时工作者
func (p *taskPool) dismissWorker() {
	p.cancelMutex.Lock()
	defer p.cancelMutex.Unlock()

	dismissNum := int(math.Ceil(float64(len(p.cancelFunc) / 4)))

	for i := 0; i < len(p.cancelFunc) && i < dismissNum; i++ {
		p.cancelFunc[i]()
	}

	p.cancelFunc = p.cancelFunc[dismissNum+1:]
}

// doTask 处理单个任务
func (p *taskPool) doTask(t Task) {
	ctx := t.Ctx()
	if _, ok := ctx.Deadline(); !ok {
		// ctx里没有过期时间，需要设置一个默认过期时间，防止工作者协程被永久阻塞死
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	ch := make(chan error, 1)
	go func() {
		ch <- util.NoPanic(t.Do)()
	}()

	var err error
	select {
	case err = <-ch:
	case <-ctx.Done():
		err = fmt.Errorf("task execute timeout, %v", t.Ctx().Err())
	}

	if err != nil {
		t.SetErr(err)
	}

	t.Done()
}
