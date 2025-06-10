package iterator

func FromChan[T any](c <-chan T) Iterator[T] {
	return func() (T, bool) {
		v, ok := <-c
		return v, ok
	}
}

// None is equivalent to an iterator of length 0
func None[T any]() Iterator[T] {
	return func() (v T, ok bool) {
		return zero[T](), false
	}
}

func Replicate[T any](v T) Iterator[T] {
	return func() (T, bool) {
		return v, true
	}
}

func Natural() Iterator[int] {
	n := 0
	return func() (int, bool) {
		n++
		return n - 1, true
	}
}

// Bind is equivalent to flatMap
func Bind[Ta any, Tb any](ma Iterator[Ta], f func(Ta) Iterator[Tb]) Iterator[Tb] {
	var mb Iterator[Tb] = nil
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

func Fold[T any, Tr any](m Iterator[T], f func(Tr, T) Tr, i Tr) Iterator[Tr] {
	return Bind(m, func(t T) Iterator[Tr] {
		i = f(i, t)
		return None[Tr]().Pure(i)
	})
}

func Map[Ta any, Tb any](ma Iterator[Ta], f func(Ta) Tb) Iterator[Tb] {
	return Bind(ma, func(ta Ta) Iterator[Tb] {
		return None[Tb]().Pure(f(ta))
	})
}

func Filter[T any](m Iterator[T], f func(T) bool) Iterator[T] {
	return Bind(m, func(t T) Iterator[T] {
		if f(t) {
			return None[T]().Pure(t)
		} else {
			return None[T]()
		}
	})
}

func Reduce[T any, Tr any](m Iterator[T], f func(Tr, T) Tr, i Tr) Tr {
	mr := Fold[T, Tr](m, f, i)
	v, ok := mr.Last()
	if !ok {
		return i
	}
	return v
}
