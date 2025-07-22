package vector

type Vector[T any] interface {
	Get(i uint) T
	Set(i uint, val T) Vector[T]
	Ins(i uint, val T) Vector[T]
	Del(i uint) Vector[T]
	Iter(func(val T) bool)
	Weight() uint
	Height() uint
	Split(i uint) (Vector[T], Vector[T])
	Concat(other Vector[T]) Vector[T]
	Repr() []T
}

func NewVector[T any]() Vector[T] {
	return &vector[T]{node: nil}
}

type vector[T any] struct {
	node *node[T]
}

func (v *vector[T]) Get(i uint) T {
	return get(v.node, i)
}

func (v *vector[T]) Set(i uint, val T) Vector[T] {
	return &vector[T]{node: set(v.node, i, val)}
}

func (v *vector[T]) Ins(i uint, val T) Vector[T] {
	return &vector[T]{node: ins(v.node, i, val)}
}

func (v *vector[T]) Del(i uint) Vector[T] {
	return &vector[T]{node: del(v.node, i)}
}

func (v *vector[T]) Iter(f func(val T) bool) {
	iter(v.node, f)
}

func (v *vector[T]) Weight() uint {
	return weight(v.node)
}
func (v *vector[T]) Height() uint {
	return height(v.node)
}
func (v *vector[T]) Split(i uint) (Vector[T], Vector[T]) {
	n1, n2 := split(v.node, i)
	return &vector[T]{node: n1}, &vector[T]{node: n2}
}

func (v *vector[T]) Concat(other Vector[T]) Vector[T] {
	n1, n2 := v.node, other.(*vector[T]).node
	n3 := merge(n1, n2)
	return &vector[T]{node: n3}
}

func (v *vector[T]) Repr() []T {
	buffer := make([]T, 0, v.Weight())
	for val := range v.Iter {
		buffer = append(buffer, val)
	}
	return buffer
}
