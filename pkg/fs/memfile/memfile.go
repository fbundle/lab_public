package memfile

import (
	"errors"
	"sync"

	"github.com/fbundle/go_util/pkg/fs"
)

func NewMemFile() fs.File {
	return &memFile{
		mu: sync.RWMutex{},
		file: memOpenedFile{
			data: nil,
		},
	}
}

type memFile struct {
	mu   sync.RWMutex
	file memOpenedFile
}

func (f *memFile) OpenRead() (fs.FileReader, fs.Closer, error) {
	f.mu.RLock()
	return &f.file, func() error {
		f.mu.RUnlock()
		return nil
	}, nil
}

func (f *memFile) OpenWrite() (fs.FileWriter, fs.Closer, error) {
	f.mu.Lock()
	return &f.file, func() error {
		f.mu.Unlock()
		return nil
	}, nil
}

type memOpenedFile struct {
	data []byte
}

func (f *memOpenedFile) Size() uint64 {
	return uint64(len(f.data))
}

func (f *memOpenedFile) Read(offset uint64, length uint64, reader func([]byte)) error {
	if offset > uint64(len(f.data)) {
		return errors.New("index out of range")
	}
	remaining := uint64(len(f.data)) - offset
	if length > remaining {
		length = remaining
	}
	reader(f.data[offset : offset+length])
	return nil
}
func (f *memOpenedFile) Write(offset uint64, length uint64, writer func([]byte)) error {
	if offset+length > uint64(len(f.data)) {
		f.data = append(f.data, make([]byte, offset+length-uint64(len(f.data)))...)
	}
	writer(f.data[offset : offset+length])
	return nil
}
func (f *memOpenedFile) Truncate(length uint64) error {
	if length > uint64(len(f.data)) {
		f.data = append(f.data, make([]byte, length-uint64(len(f.data)))...)
	}
	if length < uint64(len(f.data)) {
		f.data = f.data[0:length]
	}
	return nil
}
