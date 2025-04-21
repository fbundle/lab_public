package main

import (
	"fmt"
)

// Iter is a generic iterator interface
type Iter[T any] interface {
	Next() (value T, ok bool)
	Map[S any](f func(T) S) Iter[S]
}

// sliceIter is a concrete implementation of Iter for slices
type sliceIter[T any] struct {
	data []T
	pos  int
}

func (it *sliceIter[T]) Next() (T, bool) {
	if it.pos >= len(it.data) {
		var zero T
		return zero, false
	}
	val := it.data[it.pos]
	it.pos++
	return val, true
}

func (it *sliceIter[T]) Map[S any](f func(T) S) Iter[S] {
	return &mapIter[T, S]{
		base: it,
		f:    f,
	}
}

// mapIter transforms an iterator of T into an iterator of S using function f
type mapIter[T any, S any] struct {
	base Iter[T]
	f    func(T) S
}

func (it *mapIter[T, S]) Next() (S, bool) {
	val, ok := it.base.Next()
	if !ok {
		var zero S
		return zero, false
	}
	return it.f(val), true
}

func (it *mapIter[T, S]) Map[U any](f func(S) U) Iter[U] {
	return &mapIter[S, U]{
		base: it,
		f:    f,
	}
}

// helper function to create an iterator from a slice
func FromSlice[T any](s []T) Iter[T] {
	return &sliceIter[T]{data: s}
}

func main() {
	iter := FromSlice([]int{1, 2, 3, 4, 5})

	mapped := iter.
		Map(func(x int) int { return x * 2 }).
		Map(func(x int) string { return fmt.Sprintf("val=%d", x) })

	for {
		v, ok := mapped.Next()
		if !ok {
			break
		}
		fmt.Println(v)
	}
}
