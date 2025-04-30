package uint_ntt

func (a UintNTT) shiftRight(n int) UintNTT {
	if n > a.time.Len() {
		return UintNTT{}
	}
	cTime := a.time.Slice(n, a.time.Len()).Clone()
	return UintNTT{
		time: cTime,
	}
}

// inv : let m = 2^{16n}
// approx root of f(x) = m / x - a using Newton method
// error at most 1
func (a UintNTT) pinv(n int) UintNTT {
	if a.IsZero() {
		panic("division by zero")
	}
	x := FromUint64(1)
	// Newton iteration
	for {
		// x_{n+1} = x_n + x_n - (a x_n^2) / m
		left := x.Add(x)
		right := a.Mul(x).Mul(x).shiftRight(n)
		x1, ok := left.Sub(right)
		if !ok {
			// x is always on the left of the root - this will not happen
			panic("subtract overflow")
		}
		if x1.Cmp(x) == 0 {
			break
		}
		x = x1
	}
	return x
}
