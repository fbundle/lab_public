package kv_filesystem

import (
	"encoding/json"
	"errors"
)

var ErrFileNotFound = errors.New("file_not_found")
var ErrBlockNotFound = errors.New("block_not_found")

type FileID uint64
type Block uint64

type file struct {
	ID        FileID  `json:"id"`
	Size      uint64  `json:"size"`
	BlockSize uint64  `json:"block_size"`
	BlockList []Block `json:"block_list"`
	Other     any     `json:"other"`
}

func NewBlockFS(kv KVStore[Block, []byte]) *blockFS {
	fs := &blockFS{
		kv:    kv,
		files: make(map[FileID]*file),
	}
	// load files into cache
	defer fs.writeFileMeta()

	key := Block(wrapKey(keyTypeFile, 0))
	b, ok, err := kv.Get(key)
	if err != nil {
		panic(err)
	}
	if !ok {
		return fs
	}
	err = json.Unmarshal(b, &fs.files)
	if err != nil {
		panic(err)
	}
	return fs
}

type blockFS struct {
	kv    KVStore[Block, []byte]
	files map[FileID]*file
}

func (fs *blockFS) writeFileMeta() {
	b, err := json.Marshal(fs.files)
	if err != nil {
		panic(err)
	}
	key := Block(wrapKey(keyTypeFile, 0))
	err = fs.kv.Set(key, b)
	if err != nil {
		panic(err)
	}
}

func (fs *blockFS) Size(id FileID) (uint64, error) {
	f, ok := fs.files[id]
	if !ok {
		return 0, ErrFileNotFound
	}
	return f.Size, nil
}

func (fs *blockFS) Read(id FileID, offset uint64, buffer []byte) error {
	f, ok := fs.files[id]
	if !ok {
		return ErrFileNotFound
	}

	begOffset := offset
	endOffset := offset + size(buffer) - 1

	begBlockIdx := begOffset / f.BlockSize
	endBlockIdx := endOffset / f.BlockSize

	blockList := f.BlockList[begBlockIdx : endBlockIdx+1]
	readBuffer := make([]byte, size(blockList)*f.BlockSize)
	err := readBlocks(fs.kv, readBuffer, f.BlockSize, blockList...)
	if err != nil {
		return err
	}

	copy(buffer, readBuffer[begOffset%f.BlockSize:])
	return nil
}

func (fs *blockFS) Write(id FileID, offset uint64, buffer []byte) error {
	f, ok := fs.files[id]
	if !ok {
		return ErrFileNotFound
	}

	defer fs.writeFileMeta()

	begOffset := offset
	endOffset := offset + size(buffer) - 1

	begBlockIdx := begOffset / f.BlockSize
	endBlockIdx := endOffset / f.BlockSize

	blockList := f.BlockList[begBlockIdx : endBlockIdx+1]
	writeBuffer := make([]byte, size(blockList)*f.BlockSize)

	err := readBlocks(fs.kv, writeBuffer[:f.BlockSize], f.BlockSize, blockList[0])
	if err != nil {
		return err
	}
	err = readBlocks(fs.kv, writeBuffer[size(writeBuffer)-f.BlockSize:], f.BlockSize, blockList[size(blockList)-1])
	if err != nil {
		return err
	}

	copy(writeBuffer[begOffset%f.BlockSize:], buffer)

	return writeBlocks(fs.kv, writeBuffer, f.BlockSize, blockList...)
}

func (fs *blockFS) Trunc(id FileID, length uint64) error {
	f, ok := fs.files[id]
	if !ok {
		return ErrFileNotFound
	}
	defer fs.writeFileMeta()

	endOffset := length - 1
	endBlockIdx := endOffset / f.BlockSize

	for idx := endBlockIdx + 1; idx < size(f.BlockList); idx++ {
		block := f.BlockList[idx]
		err := fs.kv.Del(block)
		if err != nil {
			return err
		}
	}

	f.BlockList = f.BlockList[:endBlockIdx+1]
	f.Size = length
	return nil
}

func readBlocks(kv KVStore[Block, []byte], buffer []byte, blockSize uint64, blockList ...Block) error {
	for i, block := range blockList {
		blockData, ok, err := kv.Get(block)
		if err != nil {
			return err
		}
		if !ok {
			return ErrBlockNotFound
		}

		copy(buffer[uint64(i)*blockSize:], blockData)
	}
	return nil
}
func writeBlocks(kv KVStore[Block, []byte], buffer []byte, blockSize uint64, blockList ...Block) error {
	for i, block := range blockList {
		if err := kv.Set(block, buffer[uint64(i)*blockSize:]); err != nil {
			return err
		}
	}
	return nil
}

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
