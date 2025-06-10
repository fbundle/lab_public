package iterator

type Iterator[T any] func() (v T, ok bool)

func (m Iterator[T]) Next() (v T, ok bool) {
	return m()
}

func (m Iterator[T]) Slice() []T {
	var vs []T
	for {
		v, ok := m()
		if !ok {
			break
		}
		vs = append(vs, v)
	}
	return vs
}

func (m Iterator[T]) Chan() <-chan T {
	ch := make(chan T)
	go func() {
		for {
			v, ok := m()
			if !ok {
				break
			}
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// Pure is equivalent to an iterator of length n
func (m Iterator[T]) Pure(vs ...T) Iterator[T] {
	i := 0
	return func() (v T, ok bool) {
		if i >= len(vs) {
			return m()
		}
		i++
		return vs[i-1], true
	}
}

func (m Iterator[T]) TakeAtMost(n int) Iterator[T] {
	return func() (v T, ok bool) {
		if n <= 0 {
			return zero[T](), false
		}
		n--
		return m()
	}
}

func (m Iterator[T]) DropAtMost(n int) Iterator[T] {
	for i := 0; i < n; i++ {
		m()
	}
	return m
}

func (m Iterator[T]) Last() (v T, ok bool) {
	empty := true
	for {
		v, ok = m()
		if !ok {
			break
		}
		empty = false
	}
	if empty {
		return zero[T](), false
	}
	return v, true
}
