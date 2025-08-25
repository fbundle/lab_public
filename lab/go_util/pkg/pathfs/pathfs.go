package pathfs

import (
	"errors"
	"strings"
)

var ErrPath = errors.New("path")

const (
	PathSeparator = "/"
)

type Path = []string

func ensurePath(path []string) bool {
	for _, name := range path {
		if strings.Contains(name, PathSeparator) {
			return false
		}
	}
	return true
}

type PathFS interface {
	OpenOrCreate(path Path) (File, error)
	Delete(path Path) error

	Walk(yield func(path Path, file File) bool)
}

type memPathFS struct {
	files map[string]File
}

func NewMemPathFS() PathFS {
	return &memPathFS{
		files: make(map[string]File),
	}
}

func (p *memPathFS) OpenOrCreate(path Path) (File, error) {
	if !ensurePath(path) {
		return nil, ErrPath
	}

	key := strings.Join(path, PathSeparator)
	file, ok := p.files[key]
	if !ok {
		file = newMemFile()
		p.files[key] = file
	}
	return file, nil
}

func (p *memPathFS) Delete(path Path) error {
	if !ensurePath(path) {
		return ErrPath
	}
	key := strings.Join(path, PathSeparator)
	delete(p.files, key)
	return nil
}

func (p *memPathFS) Walk(yield func(path Path, file File) bool) {
	for key, file := range p.files {
		path := strings.Split(key, PathSeparator)
		if ok := yield(path, file); !ok {
			return
		}
	}
}
