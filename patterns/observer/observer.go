package observer

import "github.com/guanyaowen/puer/util/maths"

// 观察者模式

// Observer 观察者
type Observer interface {
	// OnNotify 收到通知时
	OnNotify(event interface{})
}

// Notifier 通知者
type Notifier interface {
	// Register 注册观察者
	Register(Observer)

	// Deregister 注销观察者
	Deregister(Observer)

	// Notify 事件通知
	Notify(interface{})
}

var _ Notifier = (*baseNotifier)(nil)

type baseNotifier struct {
	observers map[Observer]struct{}
}

// NewNotifier 创建一个通知者
func NewNotifier(observersSize int) Notifier {
	return newBaseNotifier(observersSize)
}

func newBaseNotifier(size int) *baseNotifier {
	return &baseNotifier{
		make(map[Observer]struct{}, maths.Max(size, 0)),
	}
}

func (b *baseNotifier) Register(observer Observer) {
	b.observers[observer] = struct{}{}
}

func (b *baseNotifier) Deregister(observer Observer) {
	delete(b.observers, observer)
}

func (b *baseNotifier) Notify(event interface{}) {
	for observer, _ := range b.observers {
		observer.OnNotify(event)
	}
}
