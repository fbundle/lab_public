package monad

// Fibonacci - Fibonacci sequence
var Fibonacci Monad[int] = func() Iterator[int] {
	a, b := 0, 1
	return func() (int, bool) {
		a, b = b, a+b
		return a, true
	}
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	i := 3
	for i*i <= n {
		if n%i == 0 {
			return false
		}
		i += 2
	}
	return true
}

// Prime - prime sieve
var Prime Monad[int] = Filter(Map(Natural, func(n int) int {
	return 2*n + 3
}), isPrime).Insert(2)
