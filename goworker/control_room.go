package goworker

import (
	"time"
)

type controlRoom struct {
	p *taskPool
}

func (c controlRoom) Run() {
	c.p.wg.Add(1)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				c.health()
			case <-c.p.closeChan:
				c.p.wg.Done()
				return
			}
		}
	}()
}

// health 健康检查
// 待处理任务chan容量满3/4后，新增一临时工作者来协助消费任务，临时工作者数量最多是core工作者数量的5倍
// 任务量不足任务chan容量1/4后，解雇一批临时工作者
func (c controlRoom) health() {
	// FIXME taskChan的任务量随时再变化，有可能导致计算不准，并发问题需要解决

	var isAdd bool

	// FIXME 需要调整优化
	if len(c.p.cancelFunc) < c.p.coreNum*5 && (c.p.taskSize*3/4) < len(c.p.taskChan) {
		c.p.addWorker()
		isAdd = true
	}

	if !isAdd && len(c.p.taskChan) > c.p.coreNum && len(c.p.taskChan) < (c.p.taskSize/4) {
		c.p.dismissWorker()
	}
}
