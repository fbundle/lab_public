package dispatcher

func NewGoDispatcher() Dispatcher {
	return &goDispatcher{}
}

type goDispatcher struct{}

func (d *goDispatcher) Dispatch(task func()) bool {
	go task()
	return true
}
