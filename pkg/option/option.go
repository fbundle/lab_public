package option

import "ca/pkg/monad"

type Option[T any] struct {
	value  T
	isFull bool
}

func None[T any]() Option[T] {
	return Option[T]{
		isFull: false,
	}
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value:  value,
		isFull: true,
	}
}

func Wrap[T1 any, T2 any](f func(T1) T2) func(Option[T1]) Option[T2] {
	return func(o Option[T1]) Option[T2] {
		if !o.isFull {
			return None[T2]()
		}
		return Some(f(o.value))
	}
}

func Match[T1 any, T2 any](o Option[T1], f func(T1) T2, g func() T2) T2 {
	if o.isFull {
		return f(o.value)
	} else {
		return g()
	}

}

func (o Option[T]) Monad() monad.Monad[T] {
	consume := false
	return func() (v T, ok bool) {
		if consume {
			return o.value, false
		}
		consume = true
		return o.value, true
	}
}
