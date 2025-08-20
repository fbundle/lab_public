package named_tree

type Tree[T any] struct {
	Data     T
	children map[string]*Tree[T]
}

func (n *Tree[T]) Set(name string, child *Tree[T]) *Tree[T] {
	if n.children == nil {
		n.children = make(map[string]*Tree[T])
	}
	if _, ok := n.children[name]; ok {
		panic("child name already exists")
	}
	n.children[name] = child
	return child
}

func (n *Tree[T]) Get(name string) (*Tree[T], bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[name]
	return child, ok
}

func (n *Tree[T]) Del(name string) {
	if n.children == nil {
		panic("child name does not exist")
	}
	if _, ok := n.children[name]; !ok {
		panic("child name does not exist")
	}
	delete(n.children, name)
	if len(n.children) == 0 {
		n.children = nil
	}
}

func (n *Tree[T]) Iter(yield func(name string, child *Tree[T]) bool) {
	if n.children == nil {
		return
	}
	for name, child := range n.children {
		if ok := yield(name, child); !ok {
			return
		}
	}
}

func (n *Tree[T]) Resolve(path []string) *Tree[T] {
	if len(path) == 0 {
		return n
	}
	if n.children == nil {
		return nil
	}
	name, path := path[0], path[1:]
	child, ok := n.Get(name)
	if !ok {
		return nil
	}
	return child.Resolve(path)
}

func (n *Tree[T]) Walk(prefix []string) func(yield func(path []string, node *Tree[T]) bool) {
	return func(yield func(path []string, node *Tree[T]) bool) {
		if ok := yield(prefix, n); !ok {
			return
		}
		if n.children == nil {
			return
		}
		for name, child := range n.children {
			for path, node := range child.Walk(append(prefix, name)) {
				if ok := yield(path, node); !ok {
					return
				}
			}
		}
	}
}
