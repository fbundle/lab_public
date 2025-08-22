package vec

type Iter[T any] interface {
	Next() (value T, remain bool)
}

type iterFunc[T any] struct {
	i int
	f func(i int) (value T, remain bool)
}

func (i *iterFunc[T]) Next() (value T, remain bool) {
	i.i++
	return i.f(i.i - 1)
}

func MakeIterFromFunc[T any](f func(i int) (value T, remain bool)) Iter[T] {
	return &iterFunc[T]{
		i: 0,
		f: f,
	}
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

func ViewIter[T any](i Iter[T]) (j Iter[T], v Vec[T]) {
	v = MakeVecFromIter(i)
	return v.Iterate(), v
}
