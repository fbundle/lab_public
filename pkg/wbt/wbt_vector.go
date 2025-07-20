package wbt

type Vector[T any] interface {
	Get(int) T
	Set(int, T) Vector[T]
	Iter(func(int, T) bool)
	Len() int
	Slice(int, int) Vector[T]
	Repr() []T
	PushBack(T) Vector[T]
	PushFront(T) Vector[T]
}

type vector[T any] struct {
	offset int
	length int
	data   OrderedMap[int, T]
}

func NewVector[T any](length int) Vector[T] {
	return &vector[T]{
		offset: 0,
		length: length,
		data:   NewOrderedMap[int, T](),
	}
}

func (v *vector[T]) Len() int {
	return v.length
}

func (v *vector[T]) Get(i int) T {
	if i < 0 || i >= v.length {
		panic("index out of range")
	}
	val, _ := v.data.Get(v.offset + i)
	return val
}

func (v *vector[T]) Iter(f func(int, T) bool) {
	for i := 0; i < v.length; i++ {
		if !f(i, v.Get(i)) {
			break
		}
	}
}

func (v *vector[T]) Set(i int, val T) Vector[T] {
	if i < 0 || i >= v.length {
		panic("index out of range")
	}
	return &vector[T]{
		offset: v.offset,
		length: v.length,
		data:   v.data.Set(v.offset+i, val),
	}
}

func (v *vector[T]) Slice(beg int, end int) Vector[T] {
	if beg < 0 || end > v.length || beg > end {
		panic("index out of range")
	}
	data := v.data
	data, _ = data.Split(v.offset + end)
	_, data = data.Split(v.offset + beg)
	return &vector[T]{
		offset: v.offset + beg,
		length: end - beg,
		data:   data,
	}
}

func (v *vector[T]) Repr() []T {
	vs := make([]T, v.length)
	v.Iter(func(i int, val T) bool {
		vs[i] = val
		return true
	})
	return vs
}

func (v *vector[T]) PushBack(val T) Vector[T] {
	return &vector[T]{
		offset: v.offset,
		length: v.length + 1,
		data:   v.data.Set(v.offset+v.length, val),
	}
}

func (v *vector[T]) PushFront(val T) Vector[T] {
	return &vector[T]{
		offset: v.offset - 1,
		length: v.length + 1,
		data:   v.data.Set(v.offset-1, val),
	}
}
