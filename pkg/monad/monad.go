package monad

type Monad[T any] func() func() (v T, ok bool)

func (m Monad[T]) Slice() []T {
	mi := m()
	var vs []T
	for {
		v, ok := mi()
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
		mi := m()
		for {
			v, ok := mi()
			if !ok {
				break
			}
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// Prepend is equivalent to a monad of length n
func (m Monad[T]) Prepend(vs ...T) Monad[T] {
	return func() func() (v T, ok bool) {
		mi := m()
		i := 0
		return func() (v T, ok bool) {
			if i >= len(vs) {
				return mi()
			}
			i++
			return vs[i-1], true
		}
	}
}

func (m Monad[T]) TakeAtMost(n int) Monad[T] {
	return func() func() (v T, ok bool) {
		mi := m()
		return func() (v T, ok bool) {
			if n <= 0 {
				return zero[T](), false
			}
			n--
			return mi()
		}
	}
}

func (m Monad[T]) DropAtMost(n int) Monad[T] {
	return func() func() (v T, ok bool) {
		mi := m()
		for i := 0; i < n; i++ {
			mi()
		}
		return mi
	}
}

func (m Monad[T]) Last() (v T, ok bool) {
	mi := m()
	empty := true
	for {
		v, ok = mi()
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
