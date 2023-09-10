package observer

import (
	"fmt"
	"testing"
	"time"
)

type TestObserver struct {
	id int
}

func (t TestObserver) OnNotify(event interface{}) {
	fmt.Printf("TestObserver %d receive：%v \n", t.id, event)
}

func TestNewNotifier(t *testing.T) {
	notifier := NewNotifier(5)

	o1 := &TestObserver{1}
	o2 := &TestObserver{2}
	o3 := &TestObserver{3}

	notifier.Register(o1)
	notifier.Register(o2)
	notifier.Register(o3)

	ticker := time.NewTicker(1 * time.Second).C
	stop := time.NewTimer(10 * time.Second).C

	for {
		select {
		case ti := <-ticker:
			notifier.Notify(ti.Unix())
		case <-stop:
			return
		}
	}
}
