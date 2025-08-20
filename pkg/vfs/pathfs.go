package vfs

import "strings"

type PathFS interface {
	OpenOrCreate(path []string) (File, error)
	Delete(path []string) error

	Walk(yield func(path []string, file File) bool)
}

type memPathFS struct {
	files map[string]File
}

func NewMemPathFS() PathFS {
	return &memPathFS{
		files: make(map[string]File),
	}
}

func (p *memPathFS) OpenOrCreate(path []string) (File, error) {
	key := strings.Join(path, "/")
	file, ok := p.files[key]
	if !ok {
		file = newMemFile()
		p.files[key] = file
	}
	return file, nil
}

func (p *memPathFS) Delete(path []string) error {
	key := strings.Join(path, "/")
	delete(p.files, key)
	return nil
}

func (p *memPathFS) Walk(yield func(path []string, file File) bool) {
	for key, file := range p.files {
		path := strings.Split(key, "/")
		if ok := yield(path, file); !ok {
			return
		}
	}
}
