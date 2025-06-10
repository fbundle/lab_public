package monad

var Fibonacci Monad[int] = func() Iterator[int] {
	a, b := 0, 1
	return func() (int, bool) {
		a, b = b, a+b
		return a, true
	}
}

var Prime Monad[int] = func() Iterator[int] {
	i := Natural.DropAtMost(2)
	return func() (int, bool) {

	}
}
