package wbt

type Node[T Comparable[T]] interface {
	Get(T) (T, bool)
	Set(T) Node[T]
	Del(T) Node[T]
	Split(T) (Node[T], Node[T])
	Iter(func(T) bool)
	Weight() int
	Height() int // for debug only
}

func New[T Comparable[T]]() Node[T] {
	return &wbt[T]{node: nil}
}

type wbt[T Comparable[T]] struct {
	node *node[T]
}

func (w *wbt[T]) Weight() int {
	return int(weight(w.node))
}

func (w *wbt[T]) Height() int {
	return int(height(w.node))
}

func (w *wbt[T]) Get(keyIn T) (T, bool) {
	return get(w.node, keyIn)
}

func (w *wbt[T]) Set(keyIn T) Node[T] {
	return &wbt[T]{node: set(w.node, keyIn)}
}
func (w *wbt[T]) Del(keyIn T) Node[T] {
	return &wbt[T]{node: del(w.node, keyIn)}
}

func (w *wbt[T]) Split(keyIn T) (Node[T], Node[T]) {
	l, r := split(w.node, keyIn)
	return &wbt[T]{node: l}, &wbt[T]{node: r}
}
func (w *wbt[T]) Iter(f func(T) bool) {
	iter(w.node, f)
}
