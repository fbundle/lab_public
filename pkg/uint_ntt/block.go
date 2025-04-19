package uint_ntt

// vec :
type vec[T any] struct {
	data []T
}

func makeVec[T any](n int) vec[T] {
	return vec[T]{make([]T, n)}
}

func (b vec[T]) clone() vec[T] {
	c := makeVec[T](b.len())
	copy(c.data, b.data)
	return c
}

func (b vec[T]) len() int {
	return len(b.data)
}

func (b vec[T]) get(i int) T {
	if i >= b.len() {
		var zero T
		return zero
	}
	return b.data[i]
}

func (b vec[T]) set(i int, v T) vec[T] {
	for i >= b.len() {
		var zero T
		b.data = append(b.data, zero)
	}
	b.data[i] = v
	return b
}

func (b vec[T]) slice(beg int, end int) vec[T] {
	for end > b.len()-1 {
		var zero T
		b.data = append(b.data, zero)
	}
	return vec[T]{b.data[beg:end]}
}
