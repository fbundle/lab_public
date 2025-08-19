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
	Close() error
}

type Node interface {
	Iter(yield func(name string, child Node) bool)
	LookUp(name string) (Node, bool)

	Mkdir(name string) (Node, error)
	Rmdir(name string) error

	Create(name string) (Node, error)
	Unlink(name string) error

	Open(name string) (File, error)

	File() (File, bool)
}

func Resolve(path []string, node Node) (Node, bool) {
	if len(path) == 0 {
		return node, true
	}
	name := path[0]
	child, ok := node.LookUp(name)
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
		_, isFile := node.File()
		if isFile {
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
	return newNodeDir()
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
func (f *memFile) Close() error {
	return nil
}

func newNodeFile(file File) *node {
	return &node{
		file: file,
	}
}

func newNodeDir() *node {
	return &node{
		children: make(map[string]*node),
	}
}

type node struct {
	file     File
	children map[string]*node
}

func (n *node) LookUp(name string) (Node, bool) {
	child, ok := n.children[name]
	return child, ok
}

func (n *node) Iter(yield func(name string, child Node) bool) {
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

func (n *node) isFile() bool {
	return n.file != nil
}

func (n *node) mkChild(name string, child *node) (Node, error) {
	if n.isFile() {
		return nil, errors.New("node is a file")
	}
	if _, ok := n.children[name]; ok {
		return nil, errors.New("child exists")
	}
	n.children[name] = child
	return child, nil
}

func (n *node) rmChildIf(name string, cond func(child *node) error) error {
	if n.isFile() {
		return errors.New("node is a file")
	}
	if _, ok := n.children[name]; !ok {
		return errors.New("child not exists")
	}
	child := n.children[name]
	if err := cond(child); err != nil {
		return err
	}
	delete(n.children, name)
	return nil
}

func (n *node) Mkdir(name string) (Node, error) {
	return n.mkChild(name, newNodeDir())
}

func (n *node) Rmdir(name string) error {
	return n.rmChildIf(name, func(child *node) error {
		if child.isFile() {
			return errors.New("node is a file")
		}
		if len(child.children) > 0 {
			return errors.New("node not empty")
		}
		return nil
	})
}

func (n *node) Create(name string) (Node, error) {
	return n.mkChild(name, newNodeFile(newMemFile()))
}

func (n *node) Unlink(name string) error {
	return n.rmChildIf(name, func(child *node) error {
		if !child.isFile() {
			return errors.New("node is not a file")
		}
		return nil
	})
}

func (n *node) Open(name string) (File, error) {
	child, ok := n.children[name]
	if !ok {
		return nil, errors.New("child not exists")
	}
	if !child.isFile() {
		return nil, errors.New("child is not a file")
	}
	return child.file, nil
}

func (n *node) File() (File, bool) {
	return n.file, n.isFile()
}
