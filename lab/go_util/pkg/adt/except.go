package adt

type Except[T any] struct {
	val T
	err error
}

func (e Except[T]) Unwrap(val *T) error {
	if val != nil {
		*val = e.val
	}
	return e.err
}

func Err[T any](err error) Except[T] {
	return Except[T]{
		err: err,
	}
}

func Ok[T any](val T) Except[T] {
	return Except[T]{
		val: val,
		err: nil,
	}
}

func Bind[T1 any, T2 any](f func(T1) Except[T2]) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return f(e1.val)
	}
}

func Map[T1 any, T2 any](f func(T1) T2) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return Ok[T2](f(e1.val))
	}
}

func Seq[T1 any, T2 any](f Except[func(T1) T2]) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if f.err != nil {
			return Err[T2](f.err)
		}
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return Ok[T2](f.val(e1.val))
	}
}
