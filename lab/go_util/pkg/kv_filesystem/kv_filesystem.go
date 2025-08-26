package kv_filesystem

import (
	"encoding/binary"
)

type keyType byte

const (
	keyTypeFile  keyType = 0
	keyTypeBlock keyType = 1
)

func wrapKey(t keyType, i uint64) []byte {
	b := make([]byte, 9)
	b[0] = byte(t)
	binary.LittleEndian.PutUint64(b[1:9], i)
	return b
}

func unwrapKey(b []byte) (t keyType, i uint64) {
	t = keyType(b[0])
	i = binary.LittleEndian.Uint64(b[1:9])
	return t, i
}

func size[T any](s []T) uint64 {
	return uint64(len(s))
}

type file struct {
	ID        uint64   `json:"id"`
	BlockSize uint64   `json:"block_size"`
	BlockList []uint64 `json:"block_list"`
	Other     any      `json:"other"`
}

type kvFilesystem struct {
	kvstore KVStore
	files   map[uint64]*file
}

func (fs *kvFilesystem) Read(id uint64, offset uint64, buffer []byte) (n int, err error) {
	f := fs.files[id]
	begOffset, endOffset := offset, offset+uint64(len(buffer))

	begBlockOffset := begOffset / f.BlockSize
	endBlockOffset := endOffset / f.BlockSize
	extendedBuffer := make([]byte, f.BlockSize*(endBlockOffset-begBlockOffset))

	m, err := fs.readExactBlocks(id, begBlockOffset, endBlockOffset, extendedBuffer)
	if m > 0 {
		// shift
		extendedBuffer = extendedBuffer[begOffset-begBlockOffset*f.BlockSize:]
		n = copy(buffer, extendedBuffer)
	}
	return n, err
}

func (fs *kvFilesystem) readExactBlocks(id uint64, begBlockOffset uint64, endBlockOffset uint64, buffer []byte) (uint64, error) {
	f := fs.files[id]

	i := uint64(0)

	for idx := begBlockOffset; idx < endBlockOffset; idx++ {
		if idx >= size(f.BlockList) {
			break
		}

		bufferOffset := i * f.BlockSize
		block, err := fs.kvstore.Get(wrapKey(keyTypeBlock, f.BlockList[idx]))
		if err != nil {
			return i, err
		}
		copy(buffer[bufferOffset:], block)
		i++
	}
	return i, nil
}
