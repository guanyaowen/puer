package limiter

import (
	"sync"
)

// Releaser 释放资源/令牌
// 如使用 ConcurrentLimiter 获取的令牌就实现了该接口
type Releaser interface {
	// Release 释放
	Release()
}

// newOnceReleaser 创建一个只允许执行一次的 Releaser
func newOnceReleaser(fn func()) Releaser {
	return &onceReleaser{
		releaseFunc: fn,
	}
}

type onceReleaser struct {
	once        sync.Once
	releaseFunc func()
}

func (c *onceReleaser) Release() {
	c.once.Do(c.releaseFunc)
}
