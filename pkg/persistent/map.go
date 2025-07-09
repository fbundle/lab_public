package persistent

import (
	"go_util/pkg/persistent/internal/ordmap"
)

type Ordered[T any] interface {
	Less(T) bool
}

type Map[K Ordered[K], V any] interface {
	Get(k K) (V, bool)
	Set(k K, v V) Map[K, V]
	Len() uint
	Iter(func(K, V) bool)
}

func EmptyMap[K Ordered[K], V any]() Map[K, V] {
	return &pmap[K, V]{
		impl: ordmap.New[K, V](),
	}
}

type pmap[K Ordered[K], V any] struct {
	impl *ordmap.Node[K, V]
}

func (m *pmap[K, V]) Get(k K) (V, bool) {
	return m.impl.Get(k)
}

func (m *pmap[K, V]) Set(k K, v V) Map[K, V] {
	impl := m.impl.Insert(k, v)
	return &pmap[K, V]{
		impl: impl,
	}
}
func (m *pmap[K, V]) Len() uint {
	return uint(m.impl.Len())
}
func (m *pmap[K, V]) Iter(f func(K, V) bool) {
	for i := m.impl.Iterate(); !i.Done(); i.Next() {
		k, v := i.GetKey(), i.GetValue()
		if ok := f(k, v); !ok {
			break
		}
	}
}
