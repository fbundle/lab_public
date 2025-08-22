package fuse_util_mem

import (
	"errors"
	"sync"
	"time"

	"github.com/fbundle/go_util/pkg/fuse_util"
)

func NewMemFileStore() fuse_util.FileStore {
	return &memFileStore{
		mu:     sync.RWMutex{},
		files:  make(map[uint64]*memFile),
		lastId: 0,
	}
}

type memFileStore struct {
	mu     sync.RWMutex
	files  map[uint64]*memFile
	lastId uint64
}

func (f *memFileStore) Create() (fuse_util.File, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.lastId++
	file := &memFile{
		mu:   sync.RWMutex{},
		data: nil,
		attr: fuse_util.FileAttr{
			ID:    f.lastId,
			IsDir: false,
			Path:  nil,
			Mtime: time.Now(),
			Size:  0,
		},
	}

	f.files[f.lastId] = file
	return file, nil
}

func (f *memFileStore) Delete(id uint64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.files[id]; !ok {
		return errors.New("file not found")
	}

	delete(f.files, id)
	return nil
}

func (f *memFileStore) Iterate(yield func(file fuse_util.File) bool) error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, file := range f.files {
		if ok := yield(file); !ok {
			return nil
		}
	}
	return nil
}
