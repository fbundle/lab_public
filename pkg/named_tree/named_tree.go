package named_tree

type Tree[T any] struct {
	Data     T
	ChildMap map[string]*Tree[T]
}

func (n *Tree[T]) Set(name string, child *Tree[T]) *Tree[T] {
	if n.ChildMap == nil {
		n.ChildMap = make(map[string]*Tree[T])
	}
	if _, ok := n.ChildMap[name]; ok {
		panic("child name already exists")
	}
	n.ChildMap[name] = child
	return child
}

func (n *Tree[T]) Get(name string) (*Tree[T], bool) {
	if n.ChildMap == nil {
		return nil, false
	}
	child, ok := n.ChildMap[name]
	return child, ok
}

func (n *Tree[T]) Del(name string) {
	if n.ChildMap == nil {
		panic("child name does not exist")
	}
	if _, ok := n.ChildMap[name]; !ok {
		panic("child name does not exist")
	}
	delete(n.ChildMap, name)
	if len(n.ChildMap) == 0 {
		n.ChildMap = nil
	}
}

func (n *Tree[T]) Iter(yield func(name string, child *Tree[T]) bool) {
	if n.ChildMap == nil {
		return
	}
	for name, child := range n.ChildMap {
		if ok := yield(name, child); !ok {
			return
		}
	}
}

func (n *Tree[T]) Resolve(path []string) *Tree[T] {
	if len(path) == 0 {
		return n
	}
	if n.ChildMap == nil {
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
		if n.ChildMap == nil {
			return
		}
		for name, child := range n.ChildMap {
			for path, node := range child.Walk(append(prefix, name)) {
				if ok := yield(path, node); !ok {
					return
				}
			}
		}
	}
}
