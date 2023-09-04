package limiter

import "context"

// ConcurrentLimiter 并发限速器
type ConcurrentLimiter interface {
	Get(ctx context.Context) (Releaser, error)
}

// LimitConcurrent 并发限速器实现
type LimitConcurrent struct {
	token chan struct{}
}

// NewConcurrentLimiter 创建一个并发控制器
func NewConcurrentLimiter(max int) ConcurrentLimiter {
	if max < 1 {
		max = 10
	}
	return &LimitConcurrent{
		token: make(chan struct{}, max),
	}
}

// Get 获取令牌
//
// 若传入的是context.Background()，会一直等待，直到成功获取
// 若传入一个超时的context,在指定时间没获取到将返回错误
func (l *LimitConcurrent) Get(ctx context.Context) (Releaser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case l.token <- struct{}{}:
		return newOnceReleaser(l.release), nil
	}
}

// release 释放令牌
func (l *LimitConcurrent) release() {
	<-l.token
}
