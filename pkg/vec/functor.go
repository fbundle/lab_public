package vec

// Wrap makes Vec a functor
func Wrap[T any](f func(T) T) func(Vec[T]) Vec[T] {
	return func(v1 Vec[T]) Vec[T] {
		v2 := MakeVec[T](v1.Len())
		for i := 0; i < v1.Len(); i++ {
			v2.Set(i, f(v1.Get(i)))
		}
		return v2
	}
}
