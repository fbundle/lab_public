package vec

type Iter[T any] interface {
	Next() (value T, remain bool)
}

type iterFunc[T any] struct {
	f func() (value T, remain bool)
}

func (i *iterFunc[T]) Next() (value T, remain bool) {
	return i.f()
}

func MakeIterFromFunc[T any](f func() (value T, remain bool)) Iter[T] {
	return &iterFunc[T]{f: f}
}

type chanIter[T any] struct {
	ch chan T
}

func (c *chanIter[T]) Next() (value T, remain bool) {
	value, remain = <-c.ch
	return value, remain
}

func MakeChanIter[T any](ch chan T) Iter[T] {
	return &chanIter[T]{ch: ch}
}
