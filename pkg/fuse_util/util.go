package fuse_util

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"slices"
	"time"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
)

const (
	defaultFileMode = 0o666
	defaultDirMode  = os.ModeDir | 0o777
)

func getField[T any](o any, name string) (t T, err error) {
	v := reflect.ValueOf(o)
	// If it's a pointer, dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t, ok := v.FieldByName(name).Interface().(T)
	if !ok {
		return t, errors.New("internal error")
	}
	return t, nil
}

func getPathWithParent(m *memFS, op any) ([]string, error) {
	parent, err := getField[fuseops.InodeID](op, "Parent")
	if err != nil {
		return nil, err
	}
	name, err := getField[string](op, "Name")
	if err != nil {
		return nil, err
	}

	node, ok := m.inodePool.getNodeFromInode(parent)
	if !ok {
		return nil, fuse.ENOENT
	}
	return append(slices.Clone(node.path), name), nil
}

func getInodeAttributes(a FileAttr) fuseops.InodeAttributes {
	var mode os.FileMode
	if a.IsDir {
		mode = defaultDirMode
	} else {
		mode = defaultFileMode
	}

	return fuseops.InodeAttributes{
		Size:  a.Size,
		Nlink: 1,
		Mode:  mode,
		Mtime: a.Mtime,
	}
}
func deepCopy[T any](src T) T {
	b, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	var dst T
	err = json.Unmarshal(b, &dst)
	if err != nil {
		panic(err)
	}
	return dst
}

func mtimeReducer(parent node, child node) (newParent node) {
	parentMtime := mustAttr(parent.file).Mtime
	childMtime := mustAttr(child.file).Mtime
	if parentMtime.After(childMtime) {
		return parent
	}
	err := parent.file.UpdateAttr(func(attr FileAttr) FileAttr {
		attr.Mtime = childMtime
		return attr
	})
	if err != nil {
		panic(err)
	}
	return parent
}

func mustAttr(file File) FileAttr {
	meta, err := file.Attr()
	if err != nil {
		panic(err)
	}
	return meta
}

func newDir() File {
	return &dir{
		attr: FileAttr{
			ID:    0,
			IsDir: true,

			Path:  nil,
			Mtime: time.Now(),
			Size:  0,
		},
	}
}

type dir struct {
	attr FileAttr
}

func (d *dir) Attr() (FileAttr, error) {
	return d.attr, nil
}

func (d *dir) Read(offset uint64, buffer []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *dir) Write(offset uint64, buffer []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *dir) Trunc(size uint64) error {
	//TODO implement me
	panic("implement me")
}

func (d *dir) UpdateAttr(f func(FileAttr) FileAttr) error {
	d.attr = f(d.attr)
	return nil
}
