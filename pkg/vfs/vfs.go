package vfs

type File interface {
	Read(offset uint64, length uint64, reader func([]byte)) error
	Write(offset uint64, length uint64, writer func([]byte)) error
	Truncate(length uint64) error
	Close() error
}

type FileSystem interface {
	Mkdir(path string) error
	Rmdir(path string) error

	Iter(dir string) func(yield func(name string) bool)

	Create(path string) error
	Unlink(path string) error

	Open(path string) (File, error)
}

// in-memory file system implementation

func newMemFile() File {
	return &memFile{
		data: nil,
	}
}

type memFile struct {
	data []byte
}

func (f *memFile) Read(offset uint64, length uint64, reader func([]byte)) error {
	reader(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Write(offset uint64, length uint64, writer func([]byte)) error {
	writer(f.data[offset : offset+length])
	return nil
}
func (f *memFile) Truncate(length uint64) error {
	f.data = f.data[0:length]
	return nil
}
func (f *memFile) Close() error {
	return nil
}

func newNode(file File) FileSystem {
	var children map[string]FileSystem = nil
	if file != nil {
		children = make(map[string]FileSystem)
	}
	return &node{
		file:     file,
		children: children,
	}
}

type node struct {
	file     File
	children map[string]FileSystem
}

func (n *node) Mkdir(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Rmdir(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Iter(dir string) func(yield func(name string) bool) {
	//TODO implement me
	panic("implement me")
}

func (n *node) Create(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Unlink(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Open(path string) (File, error) {
	//TODO implement me
	panic("implement me")
}
