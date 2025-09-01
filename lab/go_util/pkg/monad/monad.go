package monad

type Iterator[T any] = func() (v T, ok bool) // TODO - make Iterator[T] = iter.Seq[T]
type Monad[T any] func() Iterator[T]

func (m Monad[T]) Iter(yield func(T) bool) {
	mi := m()
	for {
		v, ok := mi()
		if !ok {
			break
		}
		if ok := yield(v); !ok {
			break
		}
	}
}

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

// Insert is equivalent to a monad of length n
func (m Monad[T]) Insert(vs ...T) Monad[T] {
	return func() Iterator[T] {
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
	return func() Iterator[T] {
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
	return func() Iterator[T] {
		mi := m()
		for i := 0; i < n; i++ {
			mi()
		}
		return mi
	}
}

func (m Monad[T]) Head() (v T, ok bool) {
	mi := m()
	return mi()
}

func (m Monad[T]) Last() (v T, ok bool) {
	mi := m()
	ok = false
	for {
		v1, ok1 := mi()
		if !ok1 {
			break
		}
		v, ok = v1, true
	}
	return v, ok
}
