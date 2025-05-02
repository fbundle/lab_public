package int_ntt

type Int struct {
	Abs Nat
	Neg bool
}

func (a Int) IsZero() bool {
	return a.Abs.IsZero()
}

func (a Int) sameSign(b Int) bool {
	return a.Neg == b.Neg
}

func (a Int) Add(b Int) Int {
	if a.sameSign(b) {
		return Int{
			Abs: a.Abs.Add(b.Abs),
			Neg: a.Neg,
		}
	} else {
		switch a.Abs.Cmp(b.Abs) {
		case +1:
			diff, _ := a.Abs.Sub(b.Abs)
			return Int{
				Abs: diff,
				Neg: a.Neg,
			}
		case -1:
			diff, _ := b.Abs.Sub(a.Abs)
			return Int{
				Abs: diff,
				Neg: b.Neg,
			}
		default: // cmp = 0
			return Int{} // zero
		}
	}
}

func (a Int) Sub(b Int) Int {
	return a.Add(Int{
		Abs: b.Abs,
		Neg: !b.Neg, // flip sign of b
	})
}

func (a Int) Mul(b Int) Int {
	return Int{
		Abs: a.Abs.Mul(b.Abs),
		Neg: a.sameSign(b),
	}
}
func (a Int) Div(b Int) Int {
	return Int{
		Abs: a.Abs.Div(b.Abs),
		Neg: a.sameSign(b),
	}
}

func (a Int) Mod(b Int) Int {
	if b.Neg || b.IsZero() {
		panic("only mod positive number")
	}

	mod := Int{
		Abs: a.Abs.Mod(b.Abs),
		Neg: a.Neg,
	}
	if mod.Neg {
		mod = mod.Add(b)
	}
	return mod
}

func (a Int) Equal(b Int) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	return a.sameSign(b) && a.Abs.Cmp(b.Abs) == 0
}
