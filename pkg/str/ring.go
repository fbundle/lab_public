package str

type Ring[T any] interface {
	Zero() T
	One() T
	Add(T) T
	Sub(T) T
	Neg() T
	Mul(T) T
}

type EuclideanDomain[T any] interface {
	Ring[T]
	DivMod(T) (T, T)
	Div(T) T
	Mod(T) T
}

type Field[T any] interface {
	Ring[T]
	DivField(T) T
}
