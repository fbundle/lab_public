package monad

// Fibonacci - Fibonacci sequence
var Fibonacci Monad[int] = func() Iterator[int] {
	a, b := 0, 1
	return func() (int, bool) {
		a, b = b, a+b
		return a, true
	}
}

// Prime - prime sieve
var Prime Monad[int] = Filter(Map(Natural, func(n int) int {
	return 2*n + 3
}), func(n int) bool {
	if n < 3 {
		panic("n must be >= 3")
	}
	i := 3
	for i*i <= n {
		if n%i == 0 {
			return false
		}
		i += 2
	}
	return true
}).Insert(2)
