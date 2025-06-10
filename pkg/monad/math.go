package monad

var Natural Monad[uint] = func() Iterator[uint] {
	var n uint = 0
	return func() (uint, bool) {
		n++
		return n - 1, true
	}
}

// Fibonacci - Fibonacci sequence
var Fibonacci Monad[uint] = func() Iterator[uint] {
	a, b := uint(0), uint(1)
	return func() (uint, bool) {
		a, b = b, a+b
		return a, true
	}
}

// 3, 5, 7, 9, 11, ...
var oddNonUnit Monad[uint] = Map(Natural, func(n uint) uint {
	return 2*n + 3
})

// Prime -
var Prime Monad[uint] = Filter(oddNonUnit, func(n uint) bool {
	if n < 3 {
		panic("n must be >= 3")
	}
	return Reduce(oddNonUnit, func(numFactors uint, m uint) (uint, bool) {
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
