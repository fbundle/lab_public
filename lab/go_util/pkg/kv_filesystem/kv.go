package kv_filesystem

type KVStore[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, val V) error
	Del(key K) error
}
