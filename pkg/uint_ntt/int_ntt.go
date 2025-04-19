package uint_ntt

type IntNTT struct {
	Abs UintNTT
	Neg bool
}

func (a IntNTT) IsZero() bool {
	return a.Abs.IsZero()
}

func (a IntNTT) sameSign(b IntNTT) bool {
	return a.Neg == b.Neg
}

func (a IntNTT) Add(b IntNTT) IntNTT {
	if a.sameSign(b) {
		return IntNTT{
			Abs: a.Abs.Add(b.Abs),
			Neg: a.Neg,
		}
	} else {
		switch a.Abs.Cmp(b.Abs) {
		case +1:
			diff, _ := a.Abs.Sub(b.Abs)
			return IntNTT{
				Abs: diff,
				Neg: a.Neg,
			}
		case -1:
			diff, _ := b.Abs.Sub(a.Abs)
			return IntNTT{
				Abs: diff,
				Neg: b.Neg,
			}
		default: // cmp = 0
			return IntNTT{} // zero
		}
	}
}

func (a IntNTT) Sub(b IntNTT) IntNTT {
	return a.Add(IntNTT{
		Abs: b.Abs,
		Neg: !b.Neg, // flip sign of b
	})
}

func (a IntNTT) Mul(b IntNTT) IntNTT {
	return IntNTT{
		Abs: a.Abs.Mul(b.Abs),
		Neg: a.sameSign(b),
	}
}

func (a IntNTT) Equal(b IntNTT) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	return a.sameSign(b) && a.Abs.Cmp(b.Abs) == 0
}

func (a IntNTT) shiftRight(n int) IntNTT {
	return IntNTT{
		Abs: a.Abs.shiftRight(n),
		Neg: a.Neg,
	}
}

// inv : let m = 2^{16n}
// approx root of f(x) = m / x - a using Newton method
func (a IntNTT) pinv(n int) IntNTT {
	if a.IsZero() {
		panic("division by zero")
	}
	two := IntNTT{
		Abs: FromUint64(2),
		Neg: false,
	}
	x := IntNTT{
		Abs: FromUint64(1),
		Neg: false,
	}
	// Newton iteration
	// x_{n+1} = 2 x_n - (a x_n^2) / m
	for {
		// a / m = a.shiftRight(N/2)
		x1 := two.Mul(x).Sub(a.Mul(x).Mul(x).shiftRight(n))
		if x1.Equal(x) {
			break
		}
		x = x1
	}
	return a
}

func (a IntNTT) Div(b IntNTT) IntNTT {
	n := len(a.Abs.Time) + len(b.Abs.Time)
	x := b.pinv(n)
	return a.Mul(x).shiftRight(n)
}
