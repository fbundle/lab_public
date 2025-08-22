package fuse_util

import (
	"slices"
	"time"
)

type FileAttr struct {
	ID    uint64
	IsDir bool

	Path  []string
	Mtime time.Time
	Size  uint64
}

func (a FileAttr) Clone() FileAttr {
	newAttr := a
	newAttr.Path = slices.Clone(a.Path)
	return newAttr
}

type FileViewer interface {
	Attr() (FileAttr, error)
	Read(offset uint64, buffer []byte) (n int, err error)
}

type FileUpdater interface {
	Write(offset uint64, buffer []byte) (n int, err error)
	Trunc(size uint64) error
	UpdateAttr(func(FileAttr) FileAttr) error
}

type File interface {
	FileViewer
	FileUpdater
}

// FileStore - just a map[path]File - there is no hardlink or symlink
type FileStore interface {
	Create() (file File, err error)
	Delete(id uint64) (err error)
	Iterate(yield func(file File) bool) (err error)
}
