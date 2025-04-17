package str

type Ring[T any] interface {
	Add(T) T
	Mul(T) T
	Zero() T
	One() T
	Sub(T) T
	Neg() T
}
