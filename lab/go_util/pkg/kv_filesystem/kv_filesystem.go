package kv_filesystem

import "errors"

type keyType byte

const (
	keyTypeFile  keyType = 0
	keyTypeBlock keyType = 1
)

func wrapKey(t keyType, i uint64) uint64 {
	return (uint64(t) << 56) + i
}

func unwrapKey(key uint64) (t keyType, i uint64) {
	return keyType(key >> 56), key & 0x00FFFFFFFFFFFFFF
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
	kvstore KVStore[uint64, []byte]
	files   map[uint64]*file
}

func (fs *kvFilesystem) readBlocks(blockSize uint64, blockList []uint64, buffer []byte) error {
	if size(blockList)*blockSize != size(buffer) {
		return errors.New("invalid buffer size")
	}
	for i, block := range blockList {
		blockData, err := fs.kvstore.Get(block)
		if err != nil {
			return err
		}
		copy(buffer[uint64(i)*blockSize:], blockData)
	}
	return nil
}
func (fs *kvFilesystem) writeBlocks(blockSize uint64, blockList []uint64, buffer []byte) error {
	if size(blockList)*blockSize != size(buffer) {
		return errors.New("invalid buffer size")
	}
	for i, block := range blockList {
		if err := fs.kvstore.Set(block, buffer[uint64(i)*blockSize:]); err != nil {
			return err
		}
	}
	return nil
}
