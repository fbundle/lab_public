package vec

func Map[T1 any, T2 any](v Vec[T1], f func(T1) T2) {
	w := Make[T2](v.Len())
	for i := 0; i < v.Len(); i++ {
		w = w.Set(i, f(v.Get(i)))
	}
}

func Filter[T any](v Vec[T], f func(T) bool) Vec[T] {
	w := Make[T](0)
	for i := 0; i < v.Len(); i++ {
		if f(v.Get(i)) {
			w = w.Set(w.Len(), v.Get(i))
		}
	}
	return w
}

func Reduce[T any](v Vec[T], f func(T, T) T, acc T) T {
	for i := 0; i < v.Len(); i++ {
		acc = f(v.Get(i), acc)
	}
	return acc
}
