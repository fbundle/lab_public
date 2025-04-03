package ca

// PAdic : p-adic integers
type PAdic interface {
	Get(int) int
	Add(PAdic) PAdic
	Neg() PAdic
	Sub(PAdic) PAdic
	Iter() Iter
	Mul(PAdic) PAdic
	Norm() int
	Div(PAdic) PAdic
	Inv() PAdic
	Approx(n int) (int, []int)
}

type padic struct {
	prime int
	iter  Iter
	cache []int
}

func NewPAdic(prime int, iter Iter) PAdic {
	carry := 0
	return &padic{
		prime: prime,
		iter: func() int {
			val := iter()
			q, r := divmod(carry+val, prime)
			carry = q
			return r
		},
		cache: nil,
	}
}

func NewPAdicFromInt(prime int, v int) PAdic {
	return NewPAdic(prime, NewIterFromList([]int{v}, 0))
}

func (a *padic) Get(i int) int {
	for len(a.cache) <= i {
		a.cache = append(a.cache, a.iter())
	}
	return a.cache[i]
}

func (a *padic) Add(B PAdic) PAdic {
	b := B.(*padic)
	if a.prime != b.prime {
		panic("different bases")
	}
	i := 0
	return NewPAdic(a.prime, func() int {
		val := a.Get(i) + b.Get(i)
		i++
		return val
	})
}

func (a *padic) Neg() PAdic {
	i := 0
	return (&padic{
		prime: a.prime,
		iter: func() int {
			r := a.prime - a.Get(i) - 1
			i++
			return r
		},
	}).Add(&padic{
		prime: a.prime,
		iter:  NewIterFromList([]int{1}, 0),
	})
}

func (a *padic) Sub(B PAdic) PAdic {
	return a.Add(B.Neg())
}

func (a *padic) Iter() Iter {
	i := 0
	return func() int {
		v := a.Get(i)
		i++
		return v
	}
}

func (a *padic) Mul(B PAdic) PAdic {
	b := B.(*padic)
	if a.prime != b.prime {
		panic("different bases")
	}
	i := 0
	return NewPAdic(a.prime, func() int {
		val := 0
		for j1 := 0; j1 <= i; j1++ {
			j2 := i - j1
			val += a.Get(j1) * b.Get(j2)
		}
		i++
		return val
	})
}

func (a *padic) Div(B PAdic) PAdic {
	b := B.(*padic)
	if a.prime != b.prime {
		panic("different bases")
	}
	return a.Mul(b.Inv())
}

// Inv : [1, 1, 1, 1, ...] = 1 / (1 - p)
// ab = 1 / (1 - p) => 1 = a b(1-p)
func (a *padic) Inv() PAdic {
	return a.inv1().Mul(
		NewPAdicFromInt(a.prime, 1).Sub(NewPAdicFromInt(a.prime, a.prime)),
	)
}

// inv1 : find b so that ab = [1, 1, 1, 1, ...]
// b exists if and only if a_0 != 0
// if we add p^{-1}, then p-adic integers become p-adic number
// ... + a_{-1} p^{-1} + a_0 + a_1 p + ...
func (a *padic) inv1() PAdic {
	// a_0 b_0 + carry = 1						=> b_0 = inv_a_0 n_0
	// a_0 b_1 + a_1 b_0 + carry = 1			=> b_1 = inv_a_0 (1 - a_1 b_0 - carry) = inv_a_0 n_1
	// a_0 b_2 + a_1 b_1 + a_2 b_0 + carry = 1	=> b_2 = int_a_0 (1 - a_1 b_1 - a_2 b_0 - carry) = inv_a_0 n_2
	// ...
	if a.Get(0) == 0 {
		panic("division by zero")
	}
	inv_a_0 := invmod(a.Get(0), a.prime)
	var bList []int
	i := 0
	carry := 0
	return NewPAdic(a.prime, func() int {
		n_i := 1 - carry
		for j := 0; j < i; j++ {
			n_i -= a.Get(i-j) * bList[j]
		}
		b_i := mod(inv_a_0*n_i, a.prime)
		bList = append(bList, b_i)
		total := carry
		for j := 0; j <= i; j++ {
			total += a.Get(i-j) * bList[j]
		}
		// total mod p^i is must be
		q, r := divmod(total, a.prime)
		if r != 1 {
			panic("runtime error")
		}
		carry = q
		i++
		return b_i
	})
}

func (a *padic) Approx(n int) (int, []int) {
	approx := make([]int, n)
	s := 0
	x := 1
	for i := 0; i < n; i++ {
		s += a.Get(i) * x
		x *= a.prime
		approx[i] = a.Get(i)
	}
	return s, approx
}

func (a *padic) Norm() int {
	i := 0
	for {
		v := a.Get(i)
		if v != 0 {
			return i
		}
		i++
	}
}
