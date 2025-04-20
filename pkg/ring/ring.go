package ring

import "ca/pkg/integer"

type Set[T any] interface {
	Equal(T) bool
}

type Order[T any] interface {
	Set[T]
	Cmp(T) int
}

type Ring[T any] interface {
	Set[T]
	Zero() T
	One() T
	Add(T) T
	Sub(T) T
	Neg() T
	Mul(T) T
}

type EuclideanDomain[T any] interface {
	Ring[T]
	Norm() integer.Int
	DivMod(T) (T, T)
}

type Field[T any] interface {
	Ring[T]
	DivField(T) T
}

// EuclideanAlgorithm : return a, b so that ax + by = 1
func EuclideanAlgorithm[T EuclideanDomain[T]](x T, y T) (T, T) {
	one := integer.One
	xNorm, yNorm := x.Norm(), y.Norm()
	if xNorm.Cmp(one) < 1 || yNorm.Cmp(one) < 1 {
		panic("Euclidean Algorithm only works for norm >= 2")
	}
	cmp := xNorm.Cmp(yNorm)
	if cmp < 0 {
		// x < y here
		b, a := EuclideanAlgorithm(y, x)
		return a, b
	}
	// x >= y here
	q, r := x.DivMod(y) // x = qb + r
	// ax + by = 1 <=> a(qb + r) + by = 1 <=> ar + (aq + b)y = 1
	a, b1 := EuclideanAlgorithm(r, y)
	b := b1.Sub(a.Mul(q))
	return a, b

}
