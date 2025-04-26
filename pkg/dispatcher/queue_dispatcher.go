package dispatcher

func NewQueueDispatcher(waitCount int, runCount int) Dispatcher {
	d := &queueDispatcher{
		waitQueue: make(chan func(), waitCount),
		runSlot:   make(chan struct{}, runCount),
	}
	go func() {
		for task := range d.waitQueue {
			d.runSlot <- struct{}{}
			go func(task func()) {
				defer func() {
					<-d.runSlot
				}()
				task()
			}(task)
		}
	}()
	return d
}

type queueDispatcher struct {
	waitQueue chan func()
	runSlot   chan struct{}
}

func (d *queueDispatcher) Dispatch(task func()) (ok bool) {
	select {
	case d.waitQueue <- task:
		return true
	default:
		return false
	}
}
