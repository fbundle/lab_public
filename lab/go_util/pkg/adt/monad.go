package adt

type PureMonad[T any] interface {
	Monad() Monad[T]
}

type PureMonad2[T1 any, T2 any] interface {
	Monad2() Monad[Prod2[T1, T2]]
}

type Monad[T any] interface {
	Iter(func(T) bool)
}

type Iter[T any] func(func(T) bool)

func (i Iter[T]) Iter(f func(T) bool) {
	i(f)
}

func Pure[T any](v T) Monad[T] {
	return Iter[T](func(yield func(T) bool) {
		yield(v)
	})
}

func None[T any]() Monad[T] {
	return Iter[T](func(yield func(T) bool) {})
}

func Bind[T1 any, T2 any](m1 Monad[T1], f func(T1) Monad[T2]) Monad[T2] {
	return Iter[T2](func(yield func(T2) bool) {
		for v1 := range m1.Iter {
			m2 := f(v1)
			for v2 := range m2.Iter {
				if ok := yield(v2); !ok {
					return
				}
			}
		}
	})
}
