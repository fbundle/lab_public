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
		p, ok := i.Head()
		if !ok {
			panic("not ok")
		}
		i = Filter(i.DropAtMost(1), func(n int) bool {
			return n%p != 0 // keep those n not divided by p
		})
		return p, true
	}
}
