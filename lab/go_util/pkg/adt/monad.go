package adt

import "iter"

type Monad[T any] interface {
	Iter(func(T) bool)
}

type iterMonad[T any] func(func(T) bool)

func (im iterMonad[T]) Iter(f func(T) bool) {
	im(f)
}

func newIterMonad[T any](iter iter.Seq[T]) Monad[T] {
	return iterMonad[T](iter)
}

func Pure[T any](v T) Monad[T] {
	return newIterMonad[T](func(yield func(T) bool) {
		yield(v)
	})
}

func None[T any]() Monad[T] {
	return newIterMonad[T](func(yield func(T) bool) {})
}

func Bind[T1 any, T2 any](m1 Monad[T1], f func(T1) Monad[T2]) Monad[T2] {
	return newIterMonad[T2](func(yield func(T2) bool) {
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
