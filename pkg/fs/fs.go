package fs

import "errors"

type Closer = func() error

type File interface {
	Size() uint64
	Read(offset uint64, length uint64, reader func([]byte)) error
	Write(offset uint64, length uint64, writer func([]byte)) error
	Truncate(length uint64) error
}

var ErrNotExist = errors.New("not_exist")
var ErrExist = errors.New("exist")

type Path = []Name
type Name = string

type FileSystem interface {
	Create(path Path) (File, error)
	Delete(path Path) error
	Load(path Path) (File, error)

	List(prefix Path) (func(yield func(name Name, file File) bool), error)
	Walk(prefix Path) (func(yield func(path Path, file File) bool), error)
}
