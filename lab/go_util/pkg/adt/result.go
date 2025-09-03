package adt

type Result[T any] struct {
	Val T
	Err error
}

func (r Result[T]) Unwrap(val *T) error {
	if val != nil {
		*val = r.Val
	}
	return r.Err
}

func (r Result[T]) Monad() Monad[T] {
	return Iter[T](func(yield func(T) bool) {
		if r.Err != nil {
			return
		}
		yield(r.Val)
	})
}

func Error[T any](err error) Result[T] {
	return Result[T]{
		Err: err,
	}
}

func Some[T any](val T) Result[T] {
	return Result[T]{
		Val: val,
		Err: nil,
	}
}
