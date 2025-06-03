package monad

type Monad[T any] func() (value T, ok bool)

func Pure[T any](value T) Monad[T] {
	return func() (T, bool) {
		return value, true
	}
}

func None[T any]() Monad[T] {
	return func() (value T, ok bool) {
		return value, false
	}
}

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
