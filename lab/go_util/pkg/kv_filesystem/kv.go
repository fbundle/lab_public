package kv_filesystem

type KVStore interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, val []byte) error
	Del(key []byte) error
}
