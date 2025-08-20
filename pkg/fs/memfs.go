package fs

import (
	"errors"
	"strings"

	"github.com/fbundle/go_util/pkg/sync_util"
)

var ErrPath = errors.New("path")

func NewFlatMemFS(makeFile func() File) FileSystem {
	return &flatMemFS{
		makeFile: makeFile,
		files:    sync_util.Map[string, File]{},
	}
}

type flatMemFS struct {
	makeFile func() File
	files    sync_util.Map[string, File]
}

func (m *flatMemFS) Load(path []string) (File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, ok := m.files.Load(key)
	if !ok {
		return nil, ErrNotExist
	}
	return file, nil
}

func (m *flatMemFS) Create(path []string) (File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, loaded := m.files.LoadOrStore(
		key, m.makeFile(),
	)
	if loaded {
		return nil, ErrExist
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

func (m *flatMemFS) List(prefix []string) (func(yield func(name string, file File) bool), error) {

	it, err := m.Walk(prefix)
	if err != nil {
		return nil, err
	}
	return func(yield func(name string, file File) bool) {
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

func (m *flatMemFS) Walk(prefix []string) (func(yield func(path []string, file File) bool), error) {
	prefixKey, ok := pathToKey(prefix)
	if !ok {
		return nil, ErrPath
	}
	prefixKey += "/"

	return func(yield func(path []string, file File) bool) {
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
