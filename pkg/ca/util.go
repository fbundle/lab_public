package ca

func divmod(a int, n int) (int, int) {
	if n <= 0 || a < 0 {
		panic("n must be > 0, a must be >= 0")
	}
	return a / n, a % n
}

type Iter func() int

// shift : [1, 2, 3, 4, 5, ...] -> [v, v, v, 1, 2, 3, 4, 5, ...]
func shift(n int, v int, iter func() int) func() int {
	i := 0
	return func() int {
		if i < n {
			i++
			return v
		}
		return iter()
	}
}

func NewIterFromList(list []int, tail int) Iter {
	i := 0
	return func() int {
		if i < len(list) {
			v := list[i]
			i++
			return v
		}
		return tail
	}
}

// euclidean : find x, y so that ax + by = 1
func euclidean(a int, b int) (int, int) {
	if a < 0 || b < 0 {
		panic("a must be > 0, b must be > 0")
	}
	if a < b {
		y, x := euclidean(b, a)
		return x, y
	}
	// assume a >= b
	if b == 1 {
		return 0, 1
	}
	q, r := divmod(a, b)
	// a = qb + r
	// 1 = ax + by = (qb + r)x + by = rx + b(y + qx) = rx + b y_1
	x, y1 := euclidean(r, b)
	y := y1 - q*x
	return x, y
}

// mod : always return non-negative
func mod(a int, n int) int {
	if n <= 0 {
		panic("n must be > 0")
	}
	r := a % n
	if r < 0 {
		r = r + n
	}
	return r
}

// invmod : invert of a mod n
func invmod(a int, n int) int {
	x, _ := euclidean(a, n)
	// ax + ny = 1
	return mod(x, n)
}
