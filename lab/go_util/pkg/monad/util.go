package monad

func zero[T any]() T {
	var v T
	return v
}
