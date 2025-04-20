package vec

// Zero : return Zero value of a type
func Zero[T any]() T {
	var zero T
	return zero
}

type Vec[T any] struct {
	Data []T
}

func Make[T any](n int) Vec[T] {
	return Vec[T]{make([]T, n)}
}

func (v Vec[T]) Clone() Vec[T] {
	w := Make[T](v.Len())
	copy(w.Data, v.Data)
	return w
}

func (v Vec[T]) Len() int {
	return len(v.Data)
}

func (v Vec[T]) Get(i int) T {
	if i >= v.Len() {
		return Zero[T]()
	}
	return v.Data[i]
}

func (v Vec[T]) Set(i int, x T) Vec[T] {
	for i >= v.Len() {
		v.Data = append(v.Data, Zero[T]())
	}
	v.Data[i] = x
	return v
}

func (v Vec[T]) Slice(beg int, end int) Vec[T] {
	for end >= v.Len() {
		v.Data = append(v.Data, Zero[T]())
	}
	return Vec[T]{v.Data[beg:end]}
}
