package mp

func powmod(a, n, m uint64) uint64 {
	if a >= m || n >= m {
		return powmod(a%m, n%m, m)
	}
	if n == 0 {
		return 1
	}
	if n == 1 {
		return a % m
	}
	if n%2 == 0 {
		half := powmod(a, n/2, m)
		return (half * half) % m
	} else {
		half := powmod(a, n/2, m)
		return (half * half * a) % m

	}
}
