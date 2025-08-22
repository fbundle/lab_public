package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func New() string {
	buffer := make([]byte, 32)
	n, err := io.ReadFull(rand.Reader, buffer)
	if n != len(buffer) {
		panic("read not full")
	}
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(buffer)
}
