package monad

type Iterator[T any] = func() (val T, ok bool) // TODO - make Iterator[T] = iter.Seq[T]
type Monad[T any] func() Iterator[T]

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

func (m Monad[T]) Head() (val T, ok bool) {
	mi := m()
	return mi()
}

func (m Monad[T]) Last() (val T, ok bool) {
	mi := m()
	ok = false
	for {
		v1, ok1 := mi()
		if !ok1 {
			break
		}
		val, ok = v1, true
	}
	return val, ok
}
