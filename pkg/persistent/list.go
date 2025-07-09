package persistent

type List[T any] interface {
	Push(v T) List[T]
	Pop() (List[T], T)
	Depth() uint
	Iter(func(int, T) bool)
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

func (l *list[T]) Iter(f func(int, T) bool) {
	var ll List[T] = l
	var v T
	for ll.Depth() > 0 {
		i := int(ll.Depth() - 1)
		ll, v = ll.Pop()
		if ok := f(i, v); !ok {
			break
		}
	}
}
