package wbt

type Vector[T any] interface {
	Get(int) T
	Set(int, T) Vector[T]
	Append(T) Vector[T]
	Len() int
	Slice(int, int) Vector[T]
	Repr() []T
}

type vector[T any] struct {
	offset int
	length int
	omap   OrderedMap[int, T]
}

func (v *vector[T]) Get(i int) T {
	if i < 0 || i >= v.length {
		panic("index out of range")
	}
	val, _ := v.omap.Get(i + v.offset)
	return val
}

func (v *vector[T]) Set(i int, x T) Vector[T] {
	if i < 0 || i >= v.length {
		panic("index out of range")
	}
	omap := v.omap.Set(i+v.offset, x)
	return &vector[T]{
		length: v.length,
		omap:   omap,
	}
}

func (v *vector[T]) Append(x T) Vector[T] {
	omap := v.omap.Set(v.length+v.offset, x)
	return &vector[T]{
		length: v.length + 1,
		omap:   omap,
	}
}

func (v *vector[T]) Len() int {
	return v.length
}

func (v *vector[T]) Slice(i int, j int) Vector[T] {
	omap1, _ := v.omap.Split(j + v.offset)
	_, omap2 := omap1.Split(i + v.offset)
	return &vector[T]{
		offset: v.offset + i,
		length: j - i,
		omap:   omap2,
	}
}

func (v *vector[T]) Repr() []T {
	m := make([]T, v.length)
	for key, val := range v.omap.Iter {
		m[key-v.offset] = val
	}
	return m
}

func NewVector[T any](length int) Vector[T] {
	return &vector[T]{
		offset: 0,
		length: length,
		omap:   NewOrderedMap[int, T](),
	}
}
