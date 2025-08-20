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
