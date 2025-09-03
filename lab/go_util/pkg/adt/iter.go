package adt

import "iter"

func Pure[T any](v T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(v)
	}
}

func None[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

func Bind[T1 any, T2 any](i1 iter.Seq[T1], f func(T1) iter.Seq[T2]) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v1 := range i1 {
			i2 := f(v1)
			for v2 := range i2 {
				if ok := yield(v2); !ok {
					return
				}
			}
		}
	}
}
