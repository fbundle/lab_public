package ring

import "ca/pkg/integer"

type Computable[T any] interface {
	Equal(T) bool
}

type Order[T any] interface {
	Computable[T]
	Cmp(T) int
}

type Ring[T any] interface {
	Computable[T]
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
	Div(T) T
	Mod(T) T
}

type Field[T any] interface {
	Ring[T]
	DivField(T) T
}

// EuclideanAlgorithm : return a, b so that ax + by = 1
func EuclideanAlgorithm[T EuclideanDomain[T]](x T, y T) (T, T) {
	zero := x.Zero()
	one := x.One()
	if x.Equal(zero) || y.Equal(zero) {
		panic("euclidean algorithm: zero")
	}
	if x.Equal(one) {
		return one, zero
	}
	if y.Equal(one) {
		return zero, one
	}
	cmp := x.Norm().Cmp(y.Norm())
	if cmp == 0 {
		panic("euclidean algorithm: x=y => just find the inverse of x, y")
	}
	if cmp < 0 {
		// x < y here
		b, a := EuclideanAlgorithm(y, x)
		return a, b
	}
	// x > y here
	q, r := x.DivMod(y) // x = qb + r
	// ax + by = 1 <=> a(qb + r) + by = 1 <=> ar + (aq + b)y = 1
	a, b1 := EuclideanAlgorithm(r, y)
	b := b1.Sub(a.Mul(q))
	return a, b

}
