package uint_ntt

// block : polynomial in F_p[X]
type block[T any] struct {
	data []T
}

func makeBlock[T any](n int) block[T] {
	return block[T]{make([]T, n)}
}

func (b block[T]) clone() block[T] {
	c := makeBlock[T](b.len())
	copy(c.data, b.data)
	return c
}

func (b block[T]) len() int {
	return len(b.data)
}

func (b block[T]) get(i int) T {
	if i >= b.len() {
		var zero T
		return zero
	}
	return b.data[i]
}

func (b block[T]) set(i int, v T) block[T] {
	for i >= b.len() {
		var zero T
		b.data = append(b.data, zero)
	}
	b.data[i] = v
	return b
}

func (b block[T]) slice(beg int, end int) block[T] {
	for end > b.len()-1 {
		var zero T
		b.data = append(b.data, zero)
	}
	return block[T]{b.data[beg:end]}
}
