package monad

// None is equivalent to a monad of length 0
func None[T any]() Monad[T] {
	return func() Iterator[T] {
		return func() (v T, ok bool) {
			return zero[T](), false
		}
	}
}

func Replicate[T any](v T) Monad[T] {
	return func() Iterator[T] {
		return func() (T, bool) {
			return v, true
		}
	}
}

func Natural() Monad[int] {
	return func() Iterator[int] {
		n := 0
		return func() (int, bool) {
			n++
			return n - 1, true
		}
	}
}

// Bind is equivalent to flatMap
func Bind[Ta any, Tb any](ma Monad[Ta], f func(Ta) Monad[Tb]) Monad[Tb] {
	return func() Iterator[Tb] {
		mai := ma()
		var mbi func() (Tb, bool) = nil
		return func() (b Tb, ok bool) {
			for {
				if mbi != nil {
					b, ok = mbi()
					if ok {
						return b, true
					}
					mbi = nil
				}
				a, ok := mai()
				if !ok {
					return zero[Tb](), false
				}
				mbi = f(a)()
			}
		}
	}
}

func Fold[T any, Tr any](m Monad[T], f func(Tr, T) Tr, i Tr) Monad[Tr] {
	return func() Iterator[Tr] {
		mi := m()
		return func() (tb Tr, ok bool) {
			v, ok := mi()
			if !ok {
				return zero[Tr](), false
			}
			i = f(i, v)
			return i, true
		}
	}
}

func Map[Ta any, Tb any](ma Monad[Ta], f func(Ta) Tb) Monad[Tb] {
	return Bind(ma, func(ta Ta) Monad[Tb] {
		return None[Tb]().Prepend(f(ta))
	})
}

func Filter[T any](m Monad[T], f func(T) bool) Monad[T] {
	return Bind(m, func(t T) Monad[T] {
		if f(t) {
			return None[T]().Prepend(t)
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
