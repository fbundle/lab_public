package fuse_util_mem

import (
	"errors"
	"sync"
	"time"

	"github.com/fbundle/go_util/pkg/fuse_util"
)

type memFile struct {
	mu   sync.RWMutex
	attr fuse_util.FileAttr
	data []byte
}

func (f *memFile) lockRead(m func()) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	m()
}

func (f *memFile) lockWrite(m func()) {
	f.mu.Lock()
	defer f.mu.Unlock()
	m()
	// update attr
	f.attr.Size = uint64(len(f.data))
	f.attr.Mtime = time.Now()
}

func (f *memFile) Attr() (attr fuse_util.FileAttr, err error) {
	f.lockRead(func() {
		attr = f.attr.Clone()
	})
	return attr, nil
}

func (f *memFile) UpdateAttr(updater func(fuse_util.FileAttr) fuse_util.FileAttr) (err error) {
	f.lockWrite(func() {
		newAttr := updater(f.attr)
		if f.attr.ID != newAttr.ID {
			err = errors.New("id cannot be changed")
			return
		}
		if newAttr.IsDir {
			err = errors.New("cannot change to directory")
			return
		}
		if f.attr.Size != newAttr.Size {
			err = errors.New("size cannot be changed from attr")
			return
		}

		f.attr = newAttr
	})
	return nil
}

func (f *memFile) Read(offset uint64, buffer []byte) (n int, err error) {
	f.lockRead(func() {
		if uint64(len(f.data)) <= offset {
			return
		}
		n = copy(buffer, f.data[offset:])
	})
	return n, err
}
func (f *memFile) Write(offset uint64, buffer []byte) (n int, err error) {
	f.lockWrite(func() {
		size := max(offset+uint64(len(buffer)), uint64(len(f.data)))
		if size > uint64(len(f.data)) {
			f.data = append(f.data, make([]byte, size-uint64(len(f.data)))...)
		}
		n = copy(f.data[offset:], buffer)
	})
	return n, err
}
func (f *memFile) Trunc(size uint64) error {
	f.lockWrite(func() {
		if size > uint64(len(f.data)) {
			f.data = append(f.data, make([]byte, size-uint64(len(f.data)))...)
		}

		f.data = f.data[0:size]
	})
	return nil
}
