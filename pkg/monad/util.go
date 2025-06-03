package monad

func FromSlice[T any](s []T) Monad[T] {
	i := 0
	return func() (value T, ok bool) {
		if i >= len(s) {
			return value, false
		}
		value, ok = s[i], true
		i++
		return value, ok
	}
}

func ToSlice[T any](m Monad[T]) []T {
	var s []T
	for {
		v, ok := m()
		if !ok {
			break
		}
		s = append(s, v)
	}
	return s
}
