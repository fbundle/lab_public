package dispatcher

type Dispatcher interface {
	Dispatch(func()) bool
}
