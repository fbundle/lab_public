package fs

type Closer = func() error

type File interface {
	OpenRead() (FileReader, Closer, error)
	OpenWrite() (FileWriter, Closer, error)
}

type FileReader interface {
	Size() uint64
	Read(offset uint64, length uint64, reader func([]byte)) error
}

type FileWriter interface {
	FileReader
	Write(offset uint64, length uint64, writer func([]byte)) error
	Truncate(length uint64) error
}

type Name = string

type Path = []Name

type FileSystem interface {
	OpenOrCreate(path Path) (File, error)
	Delete(path Path) error

	List(prefix Path) func(yield func(name Name, file File) bool)
	Walk() func(yield func(path Path, file File) bool)
}
