package monad

func FromChan[T any](c <-chan T) Monad[T] {
	return func() (T, bool) {
		v, ok := <-c
		return v, ok
	}
}

// None is equivalent to an monad of length 0
func None[T any]() Monad[T] {
	return func() (v T, ok bool) {
		return zero[T](), false
	}
}

func Replicate[T any](v T) Monad[T] {
	return func() (T, bool) {
		return v, true
	}
}

func Natural() Monad[int] {
	n := 0
	return func() (int, bool) {
		n++
		return n - 1, true
	}
}

// Bind is equivalent to flatMap
func Bind[Ta any, Tb any](ma Monad[Ta], f func(Ta) Monad[Tb]) Monad[Tb] {
	var mb Monad[Tb] = nil
	return func() (b Tb, ok bool) {
		for {
			if mb != nil {
				b, ok = mb()
				if ok {
					return b, true
				}
				mb = nil
			}
			a, ok := ma()
			if !ok {
				return zero[Tb](), false
			}
			mb = f(a)
		}
	}
}

func Fold[T any, Tr any](m Monad[T], f func(Tr, T) Tr, i Tr) Monad[Tr] {
	return func() (tb Tr, ok bool) {
		v, ok := m.Next()
		if !ok {
			return zero[Tr](), false
		}
		i = f(i, v)
		return i, true
	}
}

func Map[Ta any, Tb any](ma Monad[Ta], f func(Ta) Tb) Monad[Tb] {
	return Bind(ma, func(ta Ta) Monad[Tb] {
		return None[Tb]().Pure(f(ta))
	})
}

func Filter[T any](m Monad[T], f func(T) bool) Monad[T] {
	return Bind(m, func(t T) Monad[T] {
		if f(t) {
			return None[T]().Pure(t)
		} else {
			return None[T]()
		}
	})
}

func Reduce[T any, Tr any](m Monad[T], f func(Tr, T) Tr, i Tr) Tr {
	mr := Fold[T, Tr](m, f, i)
	v, ok := mr.Last()
	if !ok {
		return i
	}
	return v
}
