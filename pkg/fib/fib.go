package fib

type Ring[T any] interface {
	Zero() T
	One() T
	Add(T) T
	Mul(T) T
}

func boxProd[T Ring[T]](x [3]T, y [3]T) [3]T {
	a1, b1, c1 := x[0], x[1], x[2]
	a2, b2, c2 := y[0], y[1], y[2]
	return [3]T{
		a1.Mul(a2).Add(b1.Mul(b2)),
		a1.Mul(b2).Add(b1.Mul(c2)),
		b1.Mul(b2).Add(c1.Mul(c2)),
	}
}

func boxPow[T Ring[T]](x [3]T, n uint64) [3]T {
	dummy := x[0]
	if n == 0 {
		return [3]T{
			dummy.One(),
			dummy.Zero(),
			dummy.One(),
		}
	}
	if n == 1 {
		return x
	}
	if n%2 == 0 {
		half := boxPow(x, n/2)
		return boxProd(half, half)
	} else {
		half := boxPow(x, n/2)
		return boxProd(boxProd(half, half), x)
	}
}

// boxPowIter - somehow recursive version is faster
func boxPowIter[T Ring[T]](x [3]T, n uint64) [3]T {
	dummy := x[0]
	res := [3]T{
		dummy.One(),
		dummy.Zero(),
		dummy.One(),
	}

	for n > 0 {
		if n%2 == 1 {
			res = boxProd(res, x)
		}
		x = boxProd(x, x)
		n /= 2
	}

	return res
}

func Fib[T Ring[T]](dummy T, n uint64) T {
	x := [3]T{
		dummy.Zero(),
		dummy.One(),
		dummy.One(),
	}

	x = boxPow(x, n)
	return x[1]
}
