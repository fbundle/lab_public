package uint_ntt

// vector :
type vector[T any] struct {
	data []T
}

func makeVector[T any](n int) vector[T] {
	return vector[T]{make([]T, n)}
}

func (b vector[T]) clone() vector[T] {
	c := makeVector[T](b.len())
	copy(c.data, b.data)
	return c
}

func (b vector[T]) len() int {
	return len(b.data)
}

func (b vector[T]) get(i int) T {
	if i >= b.len() {
		var zero T
		return zero
	}
	return b.data[i]
}

func (b vector[T]) set(i int, v T) vector[T] {
	for i >= b.len() {
		var zero T
		b.data = append(b.data, zero)
	}
	b.data[i] = v
	return b
}

func (b vector[T]) slice(beg int, end int) vector[T] {
	for end > b.len()-1 {
		var zero T
		b.data = append(b.data, zero)
	}
	return vector[T]{b.data[beg:end]}
}
