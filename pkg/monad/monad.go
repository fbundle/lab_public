package monad

type Monad[T any] func() (value T, ok bool)

// Pure is equivalent to an iterator of length n
func Pure[T any](values ...T) Monad[T] {
	i := 0
	return func() (value T, ok bool) {
		if i >= len(values) {
			return value, false
		}
		value, ok = values[i], true
		i++
		return value, ok
	}
}

// None is equivalent to an iterator of length 0
func None[T any]() Monad[T] {
	return Pure[T]()
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
				return b, false
			}
			mb = f(a)
		}
	}
}
