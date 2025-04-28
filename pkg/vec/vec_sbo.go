package vec

const (
	MAX_BUFFER_LENGTH = 1024
)

// VecSBO : vector with small buffer optimzation
type VecSBO[T any] struct {
	Buffer       [MAX_BUFFER_LENGTH]T
	BufferLength int
	AddonData    []T
}

func MakeVecSBO[T any](length int) VecSBO[T] {
	addonLength := max(0, length-MAX_BUFFER_LENGTH)
	if addonLength == 0 {
		return VecSBO[T]{
			Buffer:       [1024]T{},
			BufferLength: length,
			AddonData:    nil,
		}
	}
	return VecSBO[T]{
		Buffer:       [MAX_BUFFER_LENGTH]T{},
		BufferLength: MAX_BUFFER_LENGTH,
		AddonData:    make([]T, addonLength),
	}
}

func (v VecSBO[T]) Len() int {
	return v.BufferLength + len(v.AddonData)
}

func (v VecSBO[T]) Clone() VecSBO[T] {
	w := MakeVecSBO[T](v.Len())
	copy(w.AddonData, v.AddonData)
	return w
}

func (v VecSBO[T]) Get(i int) T {
	if i >= v.Len() {
		return Zero[T]()
	}
	if i < MAX_BUFFER_LENGTH {
		return v.Buffer[i]
	}
	return v.AddonData[i-MAX_BUFFER_LENGTH]
}

func (v VecSBO[T]) Set(i int, x T) VecSBO[T] {
	for i >= v.Len() {
		if i < MAX_BUFFER_LENGTH {
			v.BufferLength = i + 1
		} else {
			v.AddonData = append(v.AddonData, Zero[T]())
		}
	}
	if i < MAX_BUFFER_LENGTH {
		v.Buffer[i] = x
	} else {
		v.AddonData[i-MAX_BUFFER_LENGTH] = x
	}
	return v
}

func (v VecSBO[T]) ToVec() Vec[T] {
	w := MakeVec[T](v.Len())
	for i := 0; i < v.Len(); i++ {
		w = w.Set(i, v.Get(i))
	}
	return w
}

func MakeVecSBOFromVec[T any](v Vec[T]) VecSBO[T] {
	w := MakeVecSBO[T](v.Len())
	for i := 0; i < v.Len(); i++ {
		w = w.Set(i, v.Get(i))
	}
	return w
}

func (v VecSBO[T]) Slice(beg int, end int) VecSBO[T] {
	w := v.ToVec()
	w.Slice(beg, end)
	z := MakeVecSBOFromVec(w)
	return z
}
