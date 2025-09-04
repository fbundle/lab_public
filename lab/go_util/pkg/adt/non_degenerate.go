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

func NonNil[T any](pointer *T) Option[NonNilPointer[T]] {
	if pointer == nil {
		return None[NonNilPointer[T]]()
	} else {
		return Some[NonNilPointer[T]](nonNilPointer[T]{pointer: pointer})
	}
}

type NonNilPointer[T any] interface {
	Repr() *T
	Unwrap(*T)
}

type nonNilPointer[T any] struct {
	pointer *T
}

func (n nonNilPointer[T]) Repr() *T {
	return n.pointer
}

func (n nonNilPointer[T]) Unwrap(t *T) {
	if t != nil {
		*t = *n.pointer
	}
}
