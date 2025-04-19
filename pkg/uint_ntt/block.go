package uint_ntt

// vec :
type vec[T any] struct {
	data []T
}

func makeVec[T any](n int) vec[T] {
	return vec[T]{make([]T, n)}
}

func mapVec[T1 any, T2 any](v vec[T1], f func(T1) T2) vec[T2] {
	w := makeVec[T2](v.len())
	for i := 0; i < v.len(); i++ {
		w = w.set(i, f(v.get(i)))
	}
	return w
}

func (v vec[T]) clone() vec[T] {
	w := makeVec[T](v.len())
	copy(w.data, v.data)
	return w
}

func (v vec[T]) len() int {
	return len(v.data)
}

func (v vec[T]) get(i int) T {
	if i >= v.len() {
		var zero T
		return zero
	}
	return v.data[i]
}

func (v vec[T]) set(i int, x T) vec[T] {
	for i >= v.len() {
		var zero T
		v.data = append(v.data, zero)
	}
	v.data[i] = x
	return v
}

func (v vec[T]) slice(beg int, end int) vec[T] {
	for end > v.len()-1 {
		var zero T
		v.data = append(v.data, zero)
	}
	return vec[T]{v.data[beg:end]}
}
