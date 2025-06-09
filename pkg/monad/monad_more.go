package monad

func Prepend[T any](m Monad[T], s []T) Monad[T] {
	var v T
	return func() (value T, ok bool) {
		if len(s) == 0 {
			return m()
		}
		v, s = s[0], s[1:]
		return v, true
	}
}
