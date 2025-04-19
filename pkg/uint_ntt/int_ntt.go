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

func (a IntNTT) Div(b IntNTT) IntNTT {
	return IntNTT{
		Abs: a.Abs.Div(b.Abs),
		Neg: a.sameSign(b),
	}
}
