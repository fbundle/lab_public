package monad

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
