package vec

type Vec[T any] struct {
	data []T
}

func Make[T any](n int) Vec[T] {
	return Vec[T]{make([]T, n)}
}

func (v Vec[T]) Iter(f func(int, T) bool) {
	for i := 0; i < v.Len(); i++ {
		if !f(i, v.Get(i)) {
			break
		}
	}
}

func (v Vec[T]) Clone() Vec[T] {
	w := Make[T](v.Len())
	copy(w.data, v.data)
	return w
}

func (v Vec[T]) Len() int {
	return len(v.data)
}

func (v Vec[T]) Get(i int) T {
	if i >= v.Len() {
		var zero T
		return zero
	}
	return v.data[i]
}

func (v Vec[T]) Set(i int, x T) Vec[T] {
	for i >= v.Len() {
		var zero T
		v.data = append(v.data, zero)
	}
	v.data[i] = x
	return v
}

func (v Vec[T]) Slice(beg int, end int) Vec[T] {
	for end > v.Len()-1 {
		var zero T
		v.data = append(v.data, zero)
	}
	return Vec[T]{v.data[beg:end]}
}
