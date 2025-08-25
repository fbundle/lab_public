package proto

import (
	"errors"
	"io"
	"reflect"
)

func mustBePtrOfStruct(i interface{}) error {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return errors.New("must be pointer to struct")
	}
	if reflect.TypeOf(i).Elem().Kind() != reflect.Struct {
		return errors.New("must be pointer to struct")
	}
	return nil
}

// readUntil : read bytes from reader
// messages are separated by separator or io.EOF
// if err is io.EOF, b is empty
func readUntil(reader io.Reader, separator byte) (b []byte, err error) {
	buf := make([]byte, 1)
	n := 0
	for {
		n, err = reader.Read(buf)
		if err == nil && buf[0] == separator {
			break
		}
		b = append(b, buf[:n]...)
		if err != nil {
			break
		}
	}
	if len(b) > 0 && err == io.EOF {
		err = nil
	}
	return b, err
}
