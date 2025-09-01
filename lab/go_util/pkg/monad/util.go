package monad

import "iter"

func zero[T any]() T {
	var v T
	return v
}

func (m Monad[T]) Iter(yield func(T) bool) {
	mi := m()
	for {
		v, ok := mi()
		if !ok {
			break
		}
		if ok := yield(v); !ok {
			break
		}
	}
}

func (m Monad[T]) Slice() []T {
	mi := m()
	var vs []T
	for {
		v, ok := mi()
		if !ok {
			break
		}
		vs = append(vs, v)
	}
	return vs
}

func (m Monad[T]) Chan() <-chan T {
	ch := make(chan T)
	go func() {
		mi := m()
		for {
			v, ok := mi()
			if !ok {
				break
			}
			ch <- v
		}
		close(ch)
	}()
	return ch
}

func FromChan[T any](chFunc func() <-chan T) Monad[T] {
	return func() Iterator[T] {
		ch := chFunc()
		return func() (val T, ok bool) {
			val, ok = <-ch
			return val, ok
		}
	}
}

func FromIter[T any](iFunc func() iter.Seq[T]) Monad[T] {
	return FromChan(func() <-chan T {
		i := iFunc()
		ch := make(chan T, 1)
		go func() {
			defer close(ch)
			for v := range i {
				ch <- v
			}
		}()
		return ch
	})
}
