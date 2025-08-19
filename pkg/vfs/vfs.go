package vfs

import "errors"

type File interface {
	Read(offset uint64, length uint64, reader func([]byte)) error
	Write(offset uint64, length uint64, writer func([]byte)) error
	Truncate(length uint64) error
	Close() error
}

type Node interface {
	Iter(yield func(name string, child Node) bool)
	LookUp(name string) (Node, bool)

	Mkdir(name string) error
	Rmdir(name string) error

	Create(name string) error
	Unlink(name string) error

	Open(name string) (File, error)
}

// in-memory file system implementation

func newMemFile() File {
	return &memFile{
		data: nil,
	}
}

type memFile struct {
	data []byte
}

func (f *memFile) Read(offset uint64, length uint64, reader func([]byte)) error {
	reader(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Write(offset uint64, length uint64, writer func([]byte)) error {
	writer(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Truncate(length uint64) error {
	f.data = f.data[0:length]
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
	for name, child := range n.children {
		if ok := yield(name, child); !ok {
			return
		}
	}
}

func (n *node) isFile() bool {
	return n.file != nil
}

func (n *node) mkChild(name string, child *node) error {
	if n.isFile() {
		return errors.New("node is a file")
	}
	if _, ok := n.children[name]; ok {
		return errors.New("child exists")
	}
	n.children[name] = child
	return nil
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

func (n *node) Mkdir(name string) error {
	return n.mkChild(name, newNodeDir())
}

func (n *node) Rmdir(name string) error {
	return n.rmChildIf(name, func(child *node) error {
		if child.isFile() {
			return errors.New("node is a file")
		}
		return nil
	})
}

func (n *node) Create(name string) error {
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
