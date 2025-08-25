package file

import (
	"io"
	"io/ioutil"
	"os"
)

type File interface {
	// Read :
	Read() ([]byte, error)
	// Write :
	Write([]byte) error
	// Sync :
	Sync() error
	// Close :
	Close() error
}

func New(path string) (File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	return &file{
		file: f,
	}, nil
}

type file struct {
	file *os.File
}

func (f *file) Close() error {
	return f.file.Close()
}

func (f *file) Sync() error {
	return f.file.Sync()
}

func (f *file) Read() ([]byte, error) {
	_, err := f.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f.file)
}

func (f *file) Write(b []byte) error {
	_, err := f.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	err = f.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.file.Write(b)
	return err
}
