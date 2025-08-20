package memfs

import (
	"errors"
	"strings"

	"github.com/fbundle/go_util/pkg/fs"
	"github.com/fbundle/go_util/pkg/sync_util"
)

var ErrPath = errors.New("path")

var _ fs.FileSystem = (*flatMemFS)(nil)

type flatMemFS struct {
	makeFile func() fs.File
	files    sync_util.Map[string, fs.File]
}

func (n *flatMemFS) OpenOrCreate(path []string) (fs.File, error) {
	key, ok := pathToKey(path)
	if !ok {
		return nil, ErrPath
	}

	file, loaded := n.files.LoadOrStore(
		key, n.makeFile(),
	)
	_ = loaded
	return file, nil
}

func (n *flatMemFS) Delete(path []string) error {
	key, ok := pathToKey(path)
	if !ok {
		return ErrPath
	}
	n.files.Delete(key)
	return nil
}

func (n *flatMemFS) List(prefix []string) (func(yield func(name string, file fs.File) bool), error) {
	//TODO implement me
	panic("implement me")
}

func (n *flatMemFS) Walk(prefix []string) (func(yield func(path []string, file fs.File) bool), error) {
	prefixKey, ok := pathToKey(prefix)
	if !ok {
		return nil, ErrPath
	}
	prefixKey += "/"

	return func(yield func(path []string, file fs.File) bool) {
		for key, file := range n.files.Range {
			if strings.HasPrefix(key, prefixKey) {
				path := keyToPath(key)
				if ok := yield(path, file); !ok {
					return
				}
			}
		}
	}, nil
}
