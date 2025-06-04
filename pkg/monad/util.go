package monad

func ToSlice[T any](m Monad[T]) []T {
	var s []T
	for {
		v, ok := m()
		if !ok {
			break
		}
		s = append(s, v)
	}
	return s
}

func FromChan[T any](c <-chan T) Monad[T] {
	return func() (T, bool) {
		v, ok := <-c
		return v, ok
	}
}

func ToChan[T any](m Monad[T]) <-chan T {
	ch := make(chan T)
	go func() {
		for {
			v, ok := m()
			if !ok {
				break
			}
			ch <- v
		}
		close(ch)
	}()
	return ch
}
