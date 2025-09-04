package adt

func NonEmpty[T any](s []T) Option[NonEmptySlice[T]] {
	if len(s) == 0 {
		return None[NonEmptySlice[T]]()
	} else {
		return Some[NonEmptySlice[T]](nonEmptySlice[T](s))
	}
}

type NonEmptySlice[T any] interface {
	Repr() []T
	Head() T
	Tail() []T
	Last() T
	Init() []T
}

type nonEmptySlice[T any] []T

func (s nonEmptySlice[T]) Repr() []T {
	return s
}

func (s nonEmptySlice[T]) Head() T {
	return s[0]
}

func (s nonEmptySlice[T]) Tail() []T {
	return s[1:]
}

func (s nonEmptySlice[T]) Last() T {
	return s[len(s)-1]
}

func (s nonEmptySlice[T]) Init() []T {
	return s[:len(s)-1]
}
