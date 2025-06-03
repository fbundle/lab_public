package option

type Option[T any] struct {
	value  T
	isFull bool
}

func None[T any]() Option[T] {
	return Option[T]{
		isFull: false,
	}
}

func Some[T any](value T) Option[T] {
	return Option[T]{
		value:  value,
		isFull: true,
	}
}

func (o Option[T]) IsFull() bool {
	return o.isFull
}

func (o Option[T]) Value() T {
	return o.value
}

func (o Option[T]) Unwrap() T {
	return o.value
}

func (o Option[T]) UnwrapOr(def T) T {
	if o.isFull {
		return o.value
	}
	return def
}

func (o Option[T]) UnwrapOrElse(def func() T) T {
	if o.isFull {
		return o.value
	}
	return def()
}

func (o Option[T]) Expect(msg string) T {
	if o.isFull {
		return o.value
	}
	panic(msg)
}

func Wrap[T any](f func(T) T) func(Option[T]) Option[T] {
	return func(o Option[T]) Option[T] {
		if o.isFull {
			return Some(f(o.value))
		}
		return o
	}
}
