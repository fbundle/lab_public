package line_slice

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Unmarshaler[T any] = func([]byte) (T, error)
type Marshaler[T any] = func(T) ([]byte, error)

// LineSlice - give jsonl file random access
type LineSlice[T any] interface {
	Close() error
	Get(i int) (T, error)
	Push(v T) error
}

type lineSlice[T any] struct {
	file        *os.File
	index       []int
	unmarshaler Unmarshaler[T]
	marshaler   Marshaler[T]
	delim       byte
}

func zero[T any]() T {
	var v T
	return v
}

func (l *lineSlice[T]) Get(i int) (T, error) {
	_, err := l.file.Seek(int64(l.index[i]), io.SeekStart)
	if err != nil {
		return zero[T](), err
	}
	reader := bufio.NewReader(l.file)
	line, err := reader.ReadBytes(l.delim)
	if err != nil && err != io.EOF {
		return zero[T](), err
	}
	if line[len(line)-1] == l.delim {
		line = line[:len(line)-1] // strip l.delim at the end
	}
	v, err := l.unmarshaler(line)
	if err != nil {
		return zero[T](), err
	}
	return v, nil
}

func (l *lineSlice[T]) Push(v T) error {
	line, err := l.marshaler(v)
	if err != nil {
		return err
	}
	offset, err := l.file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = l.file.Write(append(line, l.delim))
	if err != nil {
		return err
	}
	l.index = append(l.index, int(offset))
	return nil
}

func (l *lineSlice[T]) Close() error {
	return l.file.Close()
}

func NewLineSlice[T any](path string, unmarshaler Unmarshaler[T], marshaler Marshaler[T], delim byte) (LineSlice[T], error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	// build index
	index := make([]int, 0)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	i := 0
	for {
		line, err := reader.ReadBytes(delim)
		if err != nil {
			if len(line) > 0 {
				return nil, fmt.Errorf("PARTIAL_LINE_ERROR: %s %s %s", path, err.Error(), string(line))
			}
			if err != io.EOF {
				return nil, err
			}
			// len(line) == 0 and err == io.EOF
			break
		}
		index = append(index, i)
		i += len(line)
	}
	return &lineSlice[T]{
		file:        file,
		index:       index,
		unmarshaler: unmarshaler,
		marshaler:   marshaler,
		delim:       delim,
	}, nil
}
