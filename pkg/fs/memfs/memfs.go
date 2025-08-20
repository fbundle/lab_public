package memfs

import (
	"errors"
	"strings"

	"github.com/fbundle/go_util/pkg/fs"
	"github.com/fbundle/go_util/pkg/sync_util"
)

var ErrPath = errors.New("path")

func NewFlatMemFS(makeFile func() fs.File) fs.FileSystem {
	return &flatMemFS{
		makeFile: makeFile,
		files:    sync_util.Map[string, fs.File]{},
	}
}

type flatMemFS struct {
	makeFile func() fs.File
	files    sync_util.Map[string, fs.File]
}

func (m *flatMemFS) Load(path []string) (fs.File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, ok := m.files.Load(key)
	if !ok {
		return nil, fs.ErrNotExist
	}
	return file, nil
}

func (m *flatMemFS) Create(path []string) (fs.File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, loaded := m.files.LoadOrStore(
		key, m.makeFile(),
	)
	if loaded {
		return nil, fs.ErrExist
	}
	return file, nil
}

func (m *flatMemFS) Delete(path []string) error {
	key, ok := pathToKey(path)
	if !ok {
		return ErrPath
	}
	m.files.Delete(key)
	return nil
}

func (m *flatMemFS) List(prefix []string) (func(yield func(name string, file fs.File) bool), error) {

	it, err := m.Walk(prefix)
	if err != nil {
		return nil, err
	}
	return func(yield func(name string, file fs.File) bool) {
		for path, file := range it {
			if len(path) == len(prefix)+1 { // 1 level below
				name := path[len(path)-1]
				if ok := yield(name, file); !ok {
					return
				}
			}
		}
	}, nil
}

func (m *flatMemFS) Walk(prefix []string) (func(yield func(path []string, file fs.File) bool), error) {
	prefixKey, ok := pathToKey(prefix)
	if !ok {
		return nil, ErrPath
	}
	prefixKey += "/"

	return func(yield func(path []string, file fs.File) bool) {
		for key, file := range m.files.Range {
			if strings.HasPrefix(key, prefixKey) {
				path := keyToPath(key)
				if ok := yield(path, file); !ok {
					return
				}
			}
		}
	}, nil
}
