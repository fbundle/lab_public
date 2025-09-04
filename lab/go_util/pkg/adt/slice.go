package adt

func NonEmpty[T any](s []T) Option[NonEmptySlice[T]] {
	if len(s) == 0 {
		return None[NonEmptySlice[T]]()
	} else {
		return Some[NonEmptySlice[T]](s)
	}
}

type NonEmptySlice[T any] []T

func (s NonEmptySlice[T]) Head() T {
	return s[0]
}

func (s NonEmptySlice[T]) Tail() []T {
	return s[1:]
}

func (s NonEmptySlice[T]) Last() T {
	return s[len(s)-1]
}

func (s NonEmptySlice[T]) Init() []T {
	return s[:len(s)-1]
}
