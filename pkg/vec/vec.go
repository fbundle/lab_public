package vec

// Zero : return Zero value of a type
func Zero[T any]() T {
	var zero T
	return zero
}

type Vec[T any] struct {
	Data []T
}

func MakeVec[T any](n int) Vec[T] {
	return Vec[T]{make([]T, n)}
}

func MakeVecFromIter[T any](iter Iter[T]) Vec[T] {
	v := MakeVec[T](0)
	for {
		value, remain := iter.Next()
		if !remain {
			break
		}
		v = v.Set(v.Len(), value)
	}
	return v
}

func (v Vec[T]) Clone() Vec[T] {
	w := MakeVec[T](v.Len())
	copy(w.Data, v.Data)
	return w
}

func (v Vec[T]) Len() int {
	return len(v.Data)
}

func (v Vec[T]) Get(i int) T {
	if i >= v.Len() {
		return Zero[T]()
	}
	return v.Data[i]
}

func (v Vec[T]) Set(i int, x T) Vec[T] {
	for i >= v.Len() {
		v.Data = append(v.Data, Zero[T]())
	}
	v.Data[i] = x
	return v
}

func (v Vec[T]) Slice(beg int, end int) Vec[T] {
	for end >= v.Len() {
		v.Data = append(v.Data, Zero[T]())
	}
	return Vec[T]{v.Data[beg:end]}
}

func (v Vec[T]) SliceRange(beg int, end int, step int) Vec[T] {
	s := Range{
		Beg:  beg,
		End:  end,
		Step: step,
	}
	ret := MakeVec[T](s.Len())
	for i := 0; i < s.Len(); i++ {
		ret.Data[i] = v.Get(s.Get(i))
	}
	return ret
}

func (v Vec[T]) Iterate() Iter[T] {
	return MakeIterFromFunc(func(i int) (value T, remain bool) {
		if i >= v.Len() {
			return Zero[T](), false
		}
		return v.Get(i), true
	})
}
