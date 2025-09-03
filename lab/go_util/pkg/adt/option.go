package adt

type Option[T any] struct {
	Val T
	Ok  bool
}

func (o Option[T]) Unwrap(val *T) bool {
	if val != nil {
		*val = o.Val
	}
	return o.Ok
}
func (o Option[T]) Iter(yield func(T)) {
	if o.Ok {
		yield(o.Val)
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		Ok: false,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Val: val,
		Ok:  true,
	}
}
