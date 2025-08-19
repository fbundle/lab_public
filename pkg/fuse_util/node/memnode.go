package node

import (
	"sync"

	"github.com/fbundle/go_util/pkg/sync_util"
)

type File interface {
	Size() int
	Read(offset int, buffer []byte) int
	Write(offset int, buffer []byte) int
	Truncate(size int)
}

type Node interface {
	File() File
	Iter(yield func(name string, child Node) bool)
	Get(name string) Node
	Set(name string, child Node)
	Del(name string)
}

func NewFile() File {
	return &file{
		mu:   sync.RWMutex{},
		data: nil,
	}
}

func NewFileNode(file File) Node {
	return &node{
		file: file,
	}
}
func NewDirNode() Node {
	return &node{
		children: &sync_util.Map[string, Node]{},
	}
}

// implementation

type file struct {
	mu   sync.RWMutex
	data []byte
}

func (f *file) Size() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.data)
}

func (f *file) Read(offset int, buffer []byte) int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return copy(buffer, f.data[offset:])
}

func (f *file) Write(offset int, buffer []byte) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.data) < offset+len(buffer) {
		f.data = append(f.data, make([]byte, offset+len(buffer)-len(f.data))...)
	}
	return copy(f.data[offset:], buffer)
}

func (f *file) Truncate(size int) {
	f.data = f.data[:size]
}

type node struct {
	file     File
	children *sync_util.Map[string, Node]
}

func (n *node) File() File {
	return n.file
}
func (n *node) Iter(yield func(name string, child Node) bool) {
	if n.children == nil {
		return
	}
	n.children.Range(func(key string, value Node) bool {
		return yield(key, value)
	})
}
func (n *node) Get(name string) Node {
	if n.children == nil {
		return nil
	}
	child, ok := n.children.Load(name)
	if !ok {
		return nil
	}
	return child
}

func (n *node) Set(name string, child Node) {
	if n.children == nil {
		return
	}
	n.children.Store(name, child)
}
func (n *node) Del(name string) {
	if n.children == nil {
		return
	}
	n.children.Delete(name)
}
