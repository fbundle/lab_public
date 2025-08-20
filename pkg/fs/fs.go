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

type FileSystem interface {
	Create(path []string) (File, error)
	Delete(path []string) error
	Load(path []string) (File, error)

	List(prefix []string) (func(yield func(name string, file File) bool), error)
	Walk(prefix []string) (func(yield func(path []string, file File) bool), error)
}
