package kv_filesystem

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
