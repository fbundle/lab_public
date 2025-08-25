package fuse_util

import (
	"errors"
	"time"
)

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
	return 0, errors.New("no_permission")
}

func (d *dir) Write(offset uint64, buffer []byte) (n int, err error) {
	return 0, errors.New("no_permission")
}

func (d *dir) Trunc(size uint64) error {
	return errors.New("no_permission")
}

func (d *dir) UpdateAttr(f func(FileAttr) FileAttr) error {
	d.attr = f(d.attr)
	return nil
}
