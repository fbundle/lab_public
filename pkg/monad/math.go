package monad

// Fibonacci - Fibonacci sequence
var Fibonacci Monad[int] = func() Iterator[int] {
	a, b := 0, 1
	return func() (int, bool) {
		a, b = b, a+b
		return a, true
	}
}

// 3, 5, 7, 9, 11, ...
var oddNonUnit Monad[int] = Map(Natural, func(n int) int {
	return 2*n + 3
})

// Prime -
var Prime Monad[int] = Filter(oddNonUnit, func(n int) bool {
	if n < 3 {
		panic("n must be >= 3")
	}
	return Reduce(oddNonUnit, func(numFactors int, m int) (int, bool) {
		if numFactors > 0 {
			return numFactors, false // stop condition
		}
		if m*m > n {
			return numFactors, false // stop condition
		}
		if m%n == 0 {
			return numFactors + 1, true
		} else {
			return numFactors, true
		}
	}, 0) == 0
}).Insert(2)
