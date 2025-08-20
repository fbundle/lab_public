package memfs

import (
	"errors"

	"github.com/fbundle/go_util/pkg/fs"
	"github.com/fbundle/go_util/pkg/sync_util"
)

var ErrPath = errors.New("path")

var _ fs.FileSystem = (*flatFS)(nil)

type flatFS struct {
	makeFile func() fs.File
	files    sync_util.Map[string, fs.File]
}

func (n *flatFS) OpenOrCreate(path []string) (fs.File, error) {
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

func (n *flatFS) Delete(path []string) error {
	key, ok := pathToKey(path)
	if !ok {
		return ErrPath
	}
	n.files.Delete(key)
	return nil
}

func (n *flatFS) List(prefix []string) (func(yield func(name string, file fs.File) bool), error) {
	//TODO implement me
	panic("implement me")
}

func (n *flatFS) Walk(prefix []string) (func(yield func(path []string, file fs.File) bool), error) {
	prefixKey, ok := pathToKey(prefix)
	if !ok {
		return nil, ErrPath
	}
	prefixKey += "/"

}
