package trie

import "slices"

type Trie[E comparable, T any] interface {
	Load(path []E) (value T, ok bool)
	Store(path []E, value T)
	Delete(path []E)

	IterBFS(prefix []E) func(yield func(path []E, value T) bool)
	IterDFS(prefix []E) func(yield func(path []E, value T) bool)
}

func NewSimpleTrie[E comparable, T any](value T) Trie[E, T] {
	return newSimpleTrie[E, T](value)
}
func newSimpleTrie[E comparable, T any](value T) *simpleTrie[E, T] {
	return &simpleTrie[E, T]{
		value:    value,
		children: make(map[E]*simpleTrie[E, T]),
	}
}
func zero[T any]() T {
	var z T
	return z
}

type simpleTrie[E comparable, T any] struct {
	value    T
	children map[E]*simpleTrie[E, T]
}

func (t *simpleTrie[E, T]) resolve(path []E) *simpleTrie[E, T] {
	if len(path) == 0 {
		return t
	}
	child, ok := t.children[path[0]]
	if !ok {
		return nil
	}
	return child.resolve(path[1:])
}
func (t *simpleTrie[E, T]) resolveWithCreate(path []E) *simpleTrie[E, T] {
	if len(path) == 0 {
		return t
	}
	child, ok := t.children[path[0]]
	if !ok {
		child = newSimpleTrie[E, T](zero[T]())
		t.children[path[0]] = child
	}
	return child.resolveWithCreate(path[1:])
}

func (t *simpleTrie[E, T]) Load(path []E) (value T, ok bool) {
	node := t.resolve(path)
	if node == nil {
		return zero[T](), false
	}
	return node.value, true
}

func (t *simpleTrie[E, T]) Store(path []E, value T) {
	node := t.resolveWithCreate(path)
	node.value = value
}

func (t *simpleTrie[E, T]) Delete(path []E) {
	if len(path) == 0 {
		panic("cannot delete root")
	}
	parent := t.resolve(path[:len(path)-1])
	if parent == nil {
		return
	}
	delete(parent.children, path[len(path)-1])
}

func (t *simpleTrie[E, T]) IterBFS(prefix []E) func(yield func(path []E, value T) bool) {
	root := t.resolve(prefix)
	if root == nil {
		return func(yield func(path []E, value T) bool) {
			return
		}
	}
	return func(yield func(path []E, value T) bool) {
		type trieWithPath struct {
			trie *simpleTrie[E, T]
			path []E
		}

		frontier := make([]trieWithPath, 0)
		frontier = append(frontier, trieWithPath{
			trie: root,
			path: prefix,
		})
		var node trieWithPath
		for len(frontier) > 0 {
			frontier, node = frontier[1:], frontier[0]
			if ok := yield(node.path, node.trie.value); !ok {
				return
			}
			for name, child := range node.trie.children {
				frontier = append(frontier, trieWithPath{
					trie: child,
					path: append(slices.Clone(node.path), name),
				})
			}
		}
	}
}

func (t *simpleTrie[E, T]) IterDFS(prefix []E) func(yield func(path []E, value T) bool) {
	root := t.resolve(prefix)
	if root == nil {
		return func(yield func(path []E, value T) bool) {
			return
		}
	}
	return func(yield func(path []E, value T) bool) {
		type trieWithPath struct {
			trie *simpleTrie[E, T]
			path []E
		}

		frontier := make([]trieWithPath, 0)
		frontier = append(frontier, trieWithPath{
			trie: root,
			path: prefix,
		})
		var node trieWithPath
		for len(frontier) > 0 {
			frontier, node = frontier[:len(frontier)-1], frontier[len(frontier)-1]
			if ok := yield(node.path, node.trie.value); !ok {
				return
			}
			for name, child := range node.trie.children {
				frontier = append(frontier, trieWithPath{
					trie: child,
					path: append(slices.Clone(node.path), name),
				})
			}
		}
	}
}
