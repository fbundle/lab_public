package wbt

type Vector[T any] interface {
	Get(int) (T, bool)
	Set(int, T) Vector[T]
	Append(T) Vector[T]
	Len() int
	Slice(int, int) Vector[T]
}
