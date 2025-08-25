package trie

import (
	"iter"
	"slices"
)

func New[E comparable, T any](value T) *Trie[E, T] {
	return &Trie[E, T]{
		value:    value,
		children: make(map[E]*Trie[E, T]),
	}
}

type Trie[E comparable, T any] struct {
	value    T
	children map[E]*Trie[E, T]
}

func (t *Trie[E, T]) ReduceAll(reducer func(parent T, child T) (newParent T)) T {
	for _, child := range t.children {
		t.value = reducer(t.value, child.ReduceAll(reducer))
	}
	return t.value
}

func (t *Trie[E, T]) ReducePartial(path []E, reducer func(parent T, child T) (newParent T)) {
	leaf := t.resolve(path)
	if leaf == nil {
		return
	}

	for i := len(path) - 1; i >= 0; i-- {
		node := t.resolve(path[:i])
		if node == nil {
			panic("unreachable")
		}
		for _, child := range node.children {
			node.value = reducer(node.value, child.value)
		}
	}
}

func (t *Trie[E, T]) resolve(path []E) *Trie[E, T] {
	if len(path) == 0 {
		return t
	}
	child, ok := t.children[path[0]]
	if !ok {
		return nil
	}
	return child.resolve(path[1:])
}

// Load - return the value attached to node at path
func (t *Trie[E, T]) Load(path []E) (value T, ok bool) {
	node := t.resolve(path)
	if node == nil {
		return value, false // NOTE - unused value
	}
	return node.value, true
}

// Store - store the value into the node at path
func (t *Trie[E, T]) Store(path []E, value T) (ok bool) {
	node := t.resolve(path)
	if node == nil {
		return false
	}
	node.value = value
	return true
}

// Insert - insert one child node - NOTE - Insert is not recursive
func (t *Trie[E, T]) Insert(path []E, value T) bool {
	parentPath, name := path[:len(path)-1], path[len(path)-1]

	parent := t.resolve(parentPath)
	if parent == nil {
		return false
	}
	if _, ok := parent.children[name]; ok {
		return false
	}
	parent.children[name] = New[E, T](value)
	return true
}

// Delete - delete one child node
func (t *Trie[E, T]) Delete(path []E) (ok bool) {
	parentPath, name := path[:len(path)-1], path[len(path)-1]

	parent := t.resolve(parentPath)
	if parent == nil {
		return false
	}
	if _, ok := parent.children[name]; !ok {
		return false
	}
	delete(parent.children, name)
	return true
}

// List - yield all childrens of a path
func (t *Trie[E, T]) List(prefix []E) func(yield func(name E, value T) bool) {
	root := t.resolve(prefix)
	if root == nil {
		return emptySeq2[E, T]()
	}
	return func(yield func(name E, value T) bool) {
		for name, child := range root.children {
			if ok := yield(name, child.value); !ok {
				return
			}
		}
	}
}

// Walk - yield all nodes in the subtree of a path
func (t *Trie[E, T]) Walk(prefix []E) func(yield func(path []E, value T) bool) {
	root := t.resolve(prefix)
	if root == nil {
		return emptySeq2[[]E, T]()
	}
	return func(yield func(path []E, value T) bool) {
		type trieWithPath struct {
			trie *Trie[E, T]
			path []E
		}

		frontier := make([]trieWithPath, 0)
		frontier = append(frontier, trieWithPath{
			trie: root,
			path: prefix,
		})
		var node trieWithPath
		for len(frontier) > 0 {
			frontier, node = frontier[:len(frontier)-1], frontier[len(frontier)-1] // DFS
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

func emptySeq2[K any, V any]() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {}
}
