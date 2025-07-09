package persistent

import (
	"go_util/pkg/persistent/internal/ordmap"
	"golang.org/x/exp/constraints"
)

type Map[K constraints.Ordered, V any] interface {
	Get(k K) (V, bool)
	Set(k K, v V) Map[K, V]
	Len() uint
	Iter(func(K, V) bool)
	Repr() map[K]V
}

func EmptyMap[K constraints.Ordered, V any]() Map[K, V] {
	return &pmap[K, V]{
		impl: ordmap.New[comp[K], V](),
	}
}

type comp[K constraints.Ordered] struct {
	k K
}

func (k comp[K]) Less(k2 comp[K]) bool {
	return k.k < k2.k
}

type pmap[K constraints.Ordered, V any] struct {
	impl *ordmap.Node[comp[K], V]
}

func (m *pmap[K, V]) Get(k K) (V, bool) {
	return m.impl.Get(comp[K]{k: k})
}

func (m *pmap[K, V]) Set(k K, v V) Map[K, V] {
	impl := m.impl.Insert(comp[K]{k: k}, v)
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
		if ok := f(k.k, v); !ok {
			break
		}
	}
}

func (m *pmap[K, V]) Repr() map[K]V {
	res := make(map[K]V, m.Len())
	m.Iter(func(k K, v V) bool {
		res[k] = v
		return true
	})
	return res
}
