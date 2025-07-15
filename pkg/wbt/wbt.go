package wbt

type WBT[T Key[T]] interface {
	Get(T) (T, bool)
	Set(T) WBT[T]
	Iter(func(T) bool)
}

func New[T Key[T]]() WBT[T] {
	return &wbt[T]{node: nil}
}

type wbt[T Key[T]] struct {
	node *node[T]
}

func (w *wbt[T]) Get(keyIn T) (T, bool) {
	return get(w.node, keyIn)
}

func (w *wbt[T]) Set(keyIn T) WBT[T] {
	return &wbt[T]{node: set(w.node, keyIn)}
}

func (w *wbt[T]) Iter(f func(T) bool) {
	iter(w.node, f)
}
