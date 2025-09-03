package iter_util

import "iter"

func FlatMap[T1 any, T2 any](i1 iter.Seq[T1], f func(T1) iter.Seq[T2]) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v1 := range i1 {
			v2s := f(v1)
			for v2 := range v2s {
				if ok := yield(v2); !ok {
					return
				}
			}
		}
	}
}

func Map[T1 any, T2 any](i1 iter.Seq[T1], f func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v1 := range i1 {
			if ok := yield(f(v1)); !ok {
				return
			}
		}
	}
}

func Filter[T any](i1 iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v1 := range i1 {
			if f(v1) {
				if ok := yield(v1); !ok {
					return
				}
			}
		}
	}
}

func Fold[T any](i1 iter.Seq[T], init T, f func(T, T) T) T {
	for v1 := range i1 {
		init = f(init, v1)
	}
	return init
}

func ToSlice[T any](i1 iter.Seq[T]) []T {
	var s []T
	for v1 := range i1 {
		s = append(s, v1)
	}
	return s
}

func FromSlice[T any](s []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v1 := range s {
			if ok := yield(v1); !ok {
				return
			}
		}
	}
}
