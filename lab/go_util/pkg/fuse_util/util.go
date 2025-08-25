package fuse_util

import (
	"errors"
	"os"
	"reflect"
	"slices"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
)

const (
	defaultFileMode = 0o666
	defaultDirMode  = os.ModeDir | 0o777
)

func (m *memFS) updateAllMtimeWithoutLock() {
	m.inodePool.pathToNode.ReduceAll(mtimeReducer)
}

func (m *memFS) updateMtimeWithoutLock(path []string) {
	m.inodePool.pathToNode.ReducePartial(path, mtimeReducer)
}

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
