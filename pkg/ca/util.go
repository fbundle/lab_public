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
