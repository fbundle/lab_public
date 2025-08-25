package option

type Option[T any] struct {
	val T
	err error
}

func (o Option[T]) Unwrap(err *error) T {
	*err = o.err
	return o.val
}

func Error[T any](err error) Option[T] {
	return Option[T]{
		err: err,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		val: val,
		err: nil,
	}
}

func Map[T1 any, T2 any](f func(T1) (T2, error)) func(Option[T1]) Option[T2] {
	return func(o Option[T1]) Option[T2] {
		if o.err != nil {
			return Error[T2](o.err)
		}
		val2, err := f(o.val)
		if err != nil {
			return Error[T2](err)
		}
		return Some(val2)
	}
}
