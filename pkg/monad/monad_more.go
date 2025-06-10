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

var Natural Monad[int] = func() Iterator[int] {
	n := 0
	return func() (int, bool) {
		n++
		return n - 1, true
	}
}

// Bind is equivalent to flatMap
func Bind[Tx any, Ty any](mx Monad[Tx], f func(Tx) Monad[Ty]) Monad[Ty] {
	return func() Iterator[Ty] {
		mxi := mx()
		var myi func() (Ty, bool) = nil
		return func() (y Ty, ok bool) {
			for {
				if myi != nil {
					y, ok = myi()
					if ok {
						return y, true
					}
					myi = nil
				}
				x, ok := mxi()
				if !ok {
					return zero[Ty](), false
				}
				myi = f(x)()
			}
		}
	}
}

func Fold[T any, Ta any](m Monad[T], f func(Ta, T) (Ta, bool), i Ta) Monad[Ta] {
	return func() Iterator[Ta] {
		mi := m()
		stopped := false
		return func() (ta Ta, ok bool) {
			if stopped {
				return zero[Ta](), false
			}
			v, ok := mi()
			if !ok {
				return zero[Ta](), false
			}
			i, ok = f(i, v)
			if ok {
				return i, true
			}
			stopped = true
			return zero[Ta](), false
		}
	}
}

func Map[Ta any, Tb any](ma Monad[Ta], f func(Ta) Tb) Monad[Tb] {
	return Bind(ma, func(ta Ta) Monad[Tb] {
		return None[Tb]().Insert(f(ta))
	})
}

func Filter[T any](m Monad[T], f func(T) bool) Monad[T] {
	return Bind(m, func(t T) Monad[T] {
		if f(t) {
			return None[T]().Insert(t)
		} else {
			return None[T]()
		}
	})
}

func Reduce[T any, Tr any](m Monad[T], f func(Tr, T) (Tr, bool), i Tr) Tr {
	mr := Fold[T, Tr](m, f, i)
	v, ok := mr.Last()
	if !ok {
		return i
	}
	return v
}
