package node

import (
	"sync"

	"github.com/fbundle/go_util/pkg/sync_util"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type File interface {
	Read(offset int, buffer []byte) int
	Write(offset int, buffer []byte) int
	Truncate(size int)
	Attr() *fuse.Attr
}

type Node interface {
	File() File
	Iter(yield func(name string, child Node) bool)
	Get(name string) Node
	Set(name string, child Node)
	Del(name string)
	Attr() *fuse.Attr
}

func NewFile() File {
	return &file{
		mu:   sync.RWMutex{},
		data: nil,
		attr: &fuse.Attr{
			Mode: fuse.S_IFREG | 0777,
		},
	}
}

func NewFileNode(file File) Node {
	return &node{
		file: file,
		attr: file.Attr(),
	}
}
func NewDirNode() Node {
	return &node{
		children: &sync_util.Map[string, Node]{},
		attr: &fuse.Attr{
			Mode: fuse.S_IFDIR | 0777,
		},
	}
}

// implementation

type file struct {
	mu   sync.RWMutex
	data []byte
	attr *fuse.Attr
}

func (f *file) Attr() *fuse.Attr {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.attr
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
		f.attr.Size = uint64(len(f.data))
	}
	return copy(f.data[offset:], buffer)
}

func (f *file) Truncate(size int) {
	f.data = f.data[:size]
	f.attr.Size = uint64(len(f.data))
}

type node struct {
	file     File
	children *sync_util.Map[string, Node]
	attr     *fuse.Attr
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
func (n *node) Attr() *fuse.Attr {
	return n.attr
}
