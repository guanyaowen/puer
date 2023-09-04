package worker

import (
	"context"

	"github.com/guanyaowen/puer/util"
	"github.com/guanyaowen/puer/util/jsons"
)

var _ Task = (*task)(nil)

type Task interface {
	// TaskId 获取任务ID
	TaskId() string

	// Ctx 任务退出上下文，控制超时退出、手动退出等
	Ctx() context.Context

	// Do 执行任务
	Do() error

	// Wait 等待任务执行结束
	Wait() error

	// Done 结束任务
	Done()

	// Err 获取error
	Err() error

	// SetErr 设置Err
	SetErr(err error)
}

// task 任务
type task struct {
	ctx context.Context

	taskId string

	task func() error

	err *TaskErr

	close chan struct{}
}

func newTaskId() string {
	return util.RandomNumber()
}

func NewTask(ctx context.Context, f func() error) Task {
	return &task{
		ctx:    ctx,
		taskId: newTaskId(),
		task:   f,
		close:  make(chan struct{}),
	}
}

func (t *task) Wait() error {
	<-t.close
	return t.Err()
}

func (t *task) Err() error {
	return t.err
}

func (t *task) Ctx() context.Context {
	return t.ctx
}

func (t *task) TaskId() string {
	return t.taskId
}

func (t *task) Do() (err error) {
	return t.task()
}

func (t *task) Done() {
	close(t.close)
}

func (t *task) SetErr(err error) {
	t.err = &TaskErr{
		Err:    err.Error(),
		TaskId: t.taskId,
	}
}

type TaskErr struct {
	TaskId string `json:"task_id"`
	Err    string `json:"err"`
}

func (t TaskErr) Error() string {
	return jsons.ToJson(t)
}
