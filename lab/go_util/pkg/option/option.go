package option

type Option[T any] struct {
	Val T
	Err error
}

func (o Option[T]) Unwrap(err *error) T {
	*err = o.Err
	return o.Val
}

func Error[T any](err error) Option[T] {
	return Option[T]{
		Err: err,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Val: val,
		Err: nil,
	}
}

func Map[T1 any, T2 any](f func(T1) (T2, error)) func(Option[T1]) Option[T2] {
	return func(o Option[T1]) Option[T2] {
		if o.Err != nil {
			return Error[T2](o.Err)
		}
		val2, err := f(o.Val)
		if err != nil {
			return Error[T2](err)
		}
		return Some(val2)
	}
}
