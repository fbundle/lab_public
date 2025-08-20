package vfs

import (
	"errors"

	"github.com/fbundle/go_util/pkg/named_tree"
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

	Resolve(path []string) Node
	Walk(prefix []string) func(yield func(path []string, node Node) bool)
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

func newNode() Node {
	return newNodeWithFile(nil)
}

func newNodeWithFile(file File) *node {
	tree := &named_tree.Tree[File]{
		Data: file,
	}
	return (*node)(tree)
}

type node named_tree.Tree[File]

func toTree(n *node) *named_tree.Tree[File] {
	return (*named_tree.Tree[File])(n)
}
func toNode(n *named_tree.Tree[File]) *node {
	return (*node)(n)
}

func (n *node) Iter(yield func(name string, child Node) bool) {
	toTree(n).Iter(func(name string, child *named_tree.Tree[File]) bool {
		return yield(name, toNode(child))
	})
}

func (n *node) FnChild(name string) (Node, bool) {
	child, ok := toTree(n).Get(name)
	return toNode(child), ok
}

func (n *node) MkChild(name string) (Node, error) {
	child := newNodeWithFile(nil)
	_, err := toTree(n).Set(name, toTree(child))
	return child, err
}

func (n *node) RmChild(name string) error {
	return toTree(n).Del(name)
}

func (n *node) OpenFile() (File, error) {
	file := toTree(n).Data
	if file == nil {
		return nil, errors.New("file not exist")
	}
	// TODO - open file
	return file, nil
}

func (n *node) CloseFile() error {
	file := toTree(n).Data
	if file == nil {
		return errors.New("file not exist")
	}
	// TODO - close file
	return nil
}

func (n *node) File() File {
	file := toTree(n).Data
	return file
}

func (n *node) Resolve(path []string) Node {
	return toNode(toTree(n).Resolve(path))
}

func (n *node) Walk(prefix []string) func(yield func(path []string, node Node) bool) {
	return func(yield func(path []string, node Node) bool) {
		toTree(n).Walk(prefix)(func(path []string, node *named_tree.Tree[File]) bool {
			return yield(path, toNode(node))
		})
	}
}
