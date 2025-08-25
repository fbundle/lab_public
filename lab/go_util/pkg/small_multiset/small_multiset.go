package small_multiset

type Element[T any] struct {
	Value T
	Index int
}
type Multiset[T any] struct {
	Data []*Element[T]
}

func New[T any]() *Multiset[T] {
	return &Multiset[T]{}
}

func (s *Multiset[T]) Add(v T) *Element[T] {
	e := &Element[T]{
		Value: v,
		Index: len(s.Data),
	}
	s.Data = append(s.Data, e)
	return e
}

func (s *Multiset[T]) Del(e *Element[T]) *Element[T] {
	if e.Index < 0 || e.Index >= len(s.Data) {
		return nil
	}
	tail := s.Data[len(s.Data)-1]
	tail.Index = e.Index
	s.Data[e.Index] = tail
	s.Data = s.Data[:len(s.Data)-1]
	e.Index = -1
	return e
}

func (s *Multiset[T]) Len() int {
	return len(s.Data)
}
