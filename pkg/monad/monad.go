package monad

type Monad[T any] func() (v T, ok bool)

func (m Monad[T]) Next() (v T, ok bool) {
	return m()
}

func (m Monad[T]) Slice() []T {
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

func (m Monad[T]) Chan() <-chan T {
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

// Pure is equivalent to a monad of length n
func (m Monad[T]) Pure(vs ...T) Monad[T] {
	i := 0
	return func() (v T, ok bool) {
		if i >= len(vs) {
			return m()
		}
		i++
		return vs[i-1], true
	}
}

func (m Monad[T]) TakeAtMost(n int) Monad[T] {
	return func() (v T, ok bool) {
		if n <= 0 {
			return zero[T](), false
		}
		n--
		return m()
	}
}

func (m Monad[T]) DropAtMost(n int) Monad[T] {
	for i := 0; i < n; i++ {
		m()
	}
	return m
}

func (m Monad[T]) Last() (v T, ok bool) {
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
