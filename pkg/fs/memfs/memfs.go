package memfs

import (
	"errors"
	"strings"

	"github.com/fbundle/go_util/pkg/fs"
	"github.com/fbundle/go_util/pkg/sync_util"
)

var ErrPath = errors.New("path")
var ErrNotExist = errors.New("not_exist")
var ErrExist = errors.New("exist")

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

func (fs *flatMemFS) Load(path []string) (fs.File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, ok := fs.files.Load(key)
	if !ok {
		return nil, ErrNotExist
	}
	return file, nil
}

func (fs *flatMemFS) Create(path []string) (fs.File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}
	file, loaded := fs.files.LoadOrStore(
		key, fs.makeFile(),
	)
	if loaded {
		return nil, ErrExist
	}
	return file, nil
}

func (fs *flatMemFS) Delete(path []string) error {
	key, ok := pathToKey(path)
	if !ok {
		return ErrPath
	}
	fs.files.Delete(key)
	return nil
}

func (fs *flatMemFS) List(prefix []string) (func(yield func(name string, file fs.File) bool), error) {

	it, err := fs.Walk(prefix)
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

func (fs *flatMemFS) Walk(prefix []string) (func(yield func(path []string, file fs.File) bool), error) {
	prefixKey, ok := pathToKey(prefix)
	if !ok {
		return nil, ErrPath
	}
	prefixKey += "/"

	return func(yield func(path []string, file fs.File) bool) {
		for key, file := range fs.files.Range {
			if strings.HasPrefix(key, prefixKey) {
				path := keyToPath(key)
				if ok := yield(path, file); !ok {
					return
				}
			}
		}
	}, nil
}
