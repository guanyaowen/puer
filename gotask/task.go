package gotask

import (
	"context"
	"sync"

	"github.com/guanyaowen/puer/limiter"
	"github.com/guanyaowen/puer/util"
)

// taskGroup是一个提供协程并发能力的任务组
//
// 提供了类似 err-group 的功能，doc：https://pkg.go.dev/golang.org/x/sync/errgroup
// 额外增强了限制并发数量、恐慌捕捉、任务出错打断、error捕捉返回等功能
//
// 基础使用样例看 example_test.go

type TaskGroup struct {
	ctx    context.Context
	cancel func()

	// AllowSomeFail 是否允许部分任务失败
	//
	// true: 有任务执行失败，但还是会继续执行剩余任务
	// false: 如果有任务失败，就会尽力阻止剩余任务的启动
	AllowSomeFail bool

	// 协程最大并发数, 默认是10
	Concurrent int
	wg         sync.WaitGroup

	// 并发限速器 令牌桶
	limiter limiter.ConcurrentLimiter

	errLock sync.Mutex
	err     error

	defaultOnce sync.Once
	errOnce     sync.Once
}

// WithContext 给任务组设置 context
func (g *TaskGroup) WithContext(ctx context.Context) context.Context {
	c, cancel := context.WithCancel(ctx)
	g.ctx = c
	g.cancel = cancel
	return c
}
func (g *TaskGroup) Err() error {
	g.errLock.Lock()
	defer g.errLock.Unlock()

	return g.err
}

func (g *TaskGroup) Go(fn func() error) {
	g.tgInit()

	if !g.AllowSomeFail && g.Err() != nil {
		return
	}

	// 限制并发
	releaser, err := g.limiter.Get(g.ctx)
	if err != nil {
		// ctx已关闭
		g.setError(err)
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer releaser.Release()
		if terr := util.NoPanic(fn)(); terr != nil {
			g.setError(terr)
		}
	}()
}

func (g *TaskGroup) tgInit() {
	g.defaultOnce.Do(func() {
		g.limiter = limiter.NewConcurrentLimiter(g.Concurrent)
		if g.ctx == nil {
			g.ctx, g.cancel = context.WithCancel(context.Background())
		}
	})
}

func (g *TaskGroup) setError(err error) {
	g.errOnce.Do(func() {
		g.errLock.Lock()
		defer g.errLock.Unlock()

		g.err = err

		// 尝试终止还未启动的协程
		if !g.AllowSomeFail && g.cancel != nil {
			g.cancel()
		}
	})
}

// Wait 等待执行结果
//
// err: 首次失败的error
func (g *TaskGroup) Wait() (err error) {
	g.tgInit()

	// 先等待所有任务处理完，再去关闭
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}
