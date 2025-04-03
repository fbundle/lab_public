package ca

type PArdic interface {
	Add(PArdic) PArdic
	Neg() PArdic
	Sub(PArdic) PArdic
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

func NewPArdicFromList(prime int, list []int) PArdic {
	cache := make([]int, len(list))
	for i, x := range list {
		_, cache[i] = divmod(x, prime)
	}
	return &pArdic{
		prime: prime,
		iter:  zero,
		cache: cache,
	}
}

func (a *pArdic) get(i int) int {
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
		q, r := divmod(c+a.get(i)+b.get(i), a.prime)
		c = q
		i++
		return r
	})
}

func (a *pArdic) Neg() PArdic {
	i := 0
	return NewPArdic(a.prime, func() int {
		r := a.prime - a.get(i) - 1
		i++
		return r
	}).Add(NewPArdicFromList(a.prime, []int{1}))
}

func (a *pArdic) Sub(B PArdic) PArdic {
	return a.Add(B.Neg())
}

func (a *pArdic) Approx(n int) (int, []int) {
	approx := make([]int, n)
	s := 0
	x := 1
	for i := 0; i < n; i++ {
		s += a.get(i) * x
		x *= a.prime
		approx[i] = a.get(i)
	}
	return s, approx
}
