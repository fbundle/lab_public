package adt

type Slice[T any] []T

func (s Slice[T]) NonEmpty() Option[NonEmptySlice[T]] {
	if len(s) == 0 {
		return None[NonEmptySlice[T]]()
	} else {
		return Some(NonEmptySlice[T](s))
	}

}

type NonEmptySlice[T any] []T

func (s NonEmptySlice[T]) Head() T {
	return s[0]
}

func (s NonEmptySlice[T]) Tail() Slice[T] {
	return Slice[T](s[1:])
}

func (s NonEmptySlice[T]) Last() T {
	return s[len(s)-1]
}

func (s NonEmptySlice[T]) Init() Slice[T] {
	return Slice[T](s[:len(s)-1])
}
