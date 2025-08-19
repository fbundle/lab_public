package vfs

import (
	"errors"
	"sort"
)

type File interface {
	Length() uint64
	Read(offset uint64, length uint64, reader func([]byte)) error
	Write(offset uint64, length uint64, writer func([]byte)) error
	Truncate(length uint64) error
}

type Node interface {
	Iter(yield func(name string, child Node) bool)

	FnChild(name string) (Node, bool)
	MkChild(name string) (Node, error)
	RmChild(name string) error

	OpenFile() (File, error)
	CloseFile() error

	File() File
}

func Resolve(path []string, node Node) (Node, bool) {
	if len(path) == 0 {
		return node, true
	}
	name := path[0]
	child, ok := node.FnChild(name)
	if !ok {
		return child, false
	}
	return Resolve(path[1:], child)
}

func Walk(node Node) func(yield func(path []string, node Node) bool) {
	return walk(node, nil)
}

func walk(node Node, prefix []string) func(yield func(path []string, node Node) bool) {
	return func(yield func(path []string, node Node) bool) {
		file := node.File()
		if file != nil {
			if ok := yield(prefix, node); !ok {
				return
			}
		}
		for name, child := range node.Iter {
			for path, node := range walk(child, append(prefix, name)) {
				if ok := yield(path, node); !ok {
					return
				}
			}
		}
	}
}

// in-memory file system implementation

func NewMemFS() Node {
	return newNode()
}

func newMemFile() File {
	return &memFile{
		data: nil,
	}
}

type memFile struct {
	data []byte
}

func (f *memFile) Length() uint64 {
	return uint64(len(f.data))
}

func (f *memFile) Read(offset uint64, length uint64, reader func([]byte)) error {
	if offset > uint64(len(f.data)) {
		return errors.New("index out of range")
	}
	length = min(length, uint64(len(f.data))-offset)
	reader(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Write(offset uint64, length uint64, writer func([]byte)) error {
	if offset+length > uint64(len(f.data)) {
		f.data = append(f.data, make([]byte, offset+length-uint64(len(f.data)))...)
	}
	writer(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Truncate(length uint64) error {
	if length > uint64(len(f.data)) {
		f.data = append(f.data, make([]byte, length-uint64(len(f.data)))...)
	}
	if length < uint64(len(f.data)) {
		f.data = f.data[0:length]
	}
	return nil
}

func newNode() *node {
	return &node{
		children: nil,
	}
}

type node struct {
	file     File
	children map[string]*node
}

func (n *node) FnChild(name string) (Node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[name]
	return child, ok
}

func (n *node) Iter(yield func(name string, child Node) bool) {
	if n.children == nil {
		return
	}
	names := make([]string, 0, len(n.children))
	for k := range n.children {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, name := range names {
		if ok := yield(name, n.children[name]); !ok {
			return
		}
	}
}

func (n *node) MkChild(name string) (Node, error) {
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	if _, ok := n.children[name]; ok {
		return nil, errors.New("child exists")
	}
	child := newNode()
	n.children[name] = child
	return child, nil
}

func (n *node) RmChild(name string) error {
	if n.children == nil {
		return errors.New("child not exists")
	}
	if _, ok := n.children[name]; !ok {
		return errors.New("child not exists")
	}
	delete(n.children, name)

	if len(n.children) == 0 {
		n.children = nil
	}
	return nil
}

func (n *node) OpenFile() (File, error) {
	if n.file == nil {
		n.file = newMemFile()
	}
	return n.file, nil
}

func (n *node) CloseFile() error {
	if n.file == nil {
		return errors.New("file not exists")
	}
	n.file = nil
	return nil
}

func (n *node) File() File {
	return n.file
}
