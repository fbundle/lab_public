package ordered_map

import "github.com/fbundle/lab_public/lab/go_util/pkg/adt"

func (m OrderedMap[K, V]) Monad2() adt.Monad[adt.Prod2[K, V]] {
	return adt.Iter[adt.Prod2[K, V]](func(yield func(adt.Prod2[K, V]) bool) {
		for k, v := range m.Iter {
			if ok := yield(adt.NewProd2(k, v)); !ok {
				return
			}
		}
	})
}
