package persistent

import "slices"

type List[T any] interface {
	Push(v T) List[T]
	Pop() (List[T], T)
	Depth() uint
	Iter(func(T) bool)
	Repr() []T
}

func EmptyList[T any]() List[T] {
	return &list[T]{
		depth: 0,
	}
}

type list[T any] struct {
	depth uint
	value T
	next  *list[T]
}

func (l *list[T]) Depth() uint {
	return l.depth
}

func (l *list[T]) Push(v T) List[T] {
	return &list[T]{
		depth: l.depth + 1,
		value: v,
		next:  l,
	}
}

func (l *list[T]) Pop() (List[T], T) {
	return l.next, l.value
}

func (l *list[T]) Iter(f func(T) bool) {
	var li List[T] = l
	var v T
	for li.Depth() > 0 {
		li, v = li.Pop()
		if ok := f(v); !ok {
			break
		}
	}
}

func (l *list[T]) Repr() []T {
	var res []T
	l.Iter(func(v T) bool {
		res = append(res, v)
		return true
	})
	slices.Reverse(res)
	return res
}
