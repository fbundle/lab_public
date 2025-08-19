package vfs

import (
	"strings"
)

type PathFS interface {
	OpenOrCreate(path []string) (File, error)
	Unlink(path []string) error

	Walk(yield func(path []string, file File) bool)
}

func NewMemPathFS() PathFS {
	return &pathFS{
		fileMap: make(map[string]File),
	}
}

type pathFS struct {
	fileMap map[string]File
}

func (p *pathFS) OpenOrCreate(path []string) (File, error) {
	key := strings.Join(path, "/")
	file, ok := p.fileMap[key]
	if !ok {
		file = newMemFile()
		p.fileMap[key] = file
	}
	return file, nil
}

func (p *pathFS) Unlink(path []string) error {
	key := strings.Join(path, "/")
	delete(p.fileMap, key)
	return nil
}

func (p *pathFS) Walk(yield func(path []string, file File) bool) {
	for key, file := range p.fileMap {
		path := strings.Split(key, "/")
		if ok := yield(path, file); !ok {
			return
		}
	}
}
