package ca

type PArdic interface {
	Get(i int) int
	Add(PArdic) PArdic
	Neg() PArdic
	Sub(PArdic) PArdic
	MulDigit(b int) PArdic
	Mul(B PArdic) PArdic
	Iter() func() int
	Approx(n int) (int, []int)
}

var zero = func() int {
	return 0
}

type pArdic struct {
	prime int
	iter  func() int
	cache []int
}

func NewPArdic(prime int, iter func() int) PArdic {
	return &pArdic{
		prime: prime,
		iter:  iter,
		cache: nil,
	}
}

func divmod(a int, n int) (int, int) {
	if n <= 0 || a < 0 {
		panic("n must be > 0, a must be >= 0")
	}
	return a / n, a % n
}

func NewPArdicFromList(prime int, list []int, iter func() int) PArdic {
	cache := make([]int, len(list))
	for i, x := range list {
		_, cache[i] = divmod(x, prime)
	}
	if iter == nil {
		iter = zero
	}
	return &pArdic{
		prime: prime,
		iter:  iter,
		cache: cache,
	}
}

func (a *pArdic) Get(i int) int {
	for len(a.cache) <= i {
		a.cache = append(a.cache, a.iter())
	}
	return a.cache[i]
}

func (a *pArdic) Add(B PArdic) PArdic {
	b := B.(*pArdic)
	if a.prime != b.prime {
		panic("different bases")
	}
	i := 0
	c := 0
	return NewPArdic(a.prime, func() int {
		q, r := divmod(c+a.Get(i)+b.Get(i), a.prime)
		c = q
		i++
		return r
	})
}

func (a *pArdic) Neg() PArdic {
	i := 0
	return NewPArdic(a.prime, func() int {
		r := a.prime - a.Get(i) - 1
		i++
		return r
	}).Add(NewPArdicFromList(a.prime, []int{1}, nil))
}

func (a *pArdic) Sub(B PArdic) PArdic {
	return a.Add(B.Neg())
}

func (a *pArdic) MulDigit(b int) PArdic {
	i := 0
	c := 0
	return NewPArdic(a.prime, func() int {
		q, r := divmod(c+a.Get(i)*b, a.prime)
		c = q
		i++
		return r
	})
}

func (a *pArdic) Iter() func() int {
	i := 0
	return func() int {
		v := a.Get(i)
		i++
		return v
	}
}

// shiftLeft : [1, 2, 3, 4, 5, ...] -> [v, v, v, 1, 2, 3, 4, 5, ...]
func shiftLeft(n int, v int, iter func() int) func() int {
	i := 0
	return func() int {
		if i < n {
			i++
			return v
		}
		return iter()
	}
}

func (a *pArdic) Mul(B PArdic) PArdic {
	b := B.(*pArdic)
	if a.prime != b.prime {
		panic("different bases")
	}
	var partial []PArdic
	i := 0
	c := 0
	return NewPArdic(a.prime, func() int {
		partial = append(partial, NewPArdic(
			a.prime,
			shiftLeft(i, 0, a.MulDigit(b.Get(i)).Iter())),
		)
		s := 0
		for _, p := range partial {
			s += p.Get(i)
		}
		q, r := divmod(c+s, a.prime)
		c = q
		i++
		return r
	})
}

func (a *pArdic) Approx(n int) (int, []int) {
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
