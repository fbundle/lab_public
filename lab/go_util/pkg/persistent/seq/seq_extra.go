package seq

func PushFront[T any](s Seq[T], vals ...T) Seq[T] {
	for i := len(vals) - 1; i >= 0; i-- {
		s = s.Ins(0, vals[i])
	}
	return s
}
func PushBack[T any](s Seq[T], vals ...T) Seq[T] {
	for i := 0; i < len(vals); i++ {
		s = s.Ins(s.Len(), vals[i])
	}
	return s
}

func PopFront[T any](s Seq[T]) Seq[T] {
	return s.Del(0)
}

func PopBack[T any](s Seq[T]) Seq[T] {
	return s.Del(s.Len() - 1)
}

func IndexOf[T any](s Seq[T], pred func(T) bool) int {
	index := -1
	for i, val := range s.Iter {
		if pred(val) {
			index = i
			break
		}
	}
	return index
}
func Contains[T any](s Seq[T], pred func(T) bool) bool {
	return IndexOf(s, pred) >= 0
}

func FlatMap[T1 any, T2 any](s Seq[T1], f func(T1))

func Slice[T any](s Seq[T], beg int, end int) Seq[T] {
	if beg > end {
		panic("slice out of range")
	}
	s, _ = s.Split(end)
	_, s = s.Split(beg)
	return s
}

func Merge[T any](ss ...Seq[T]) Seq[T] {
	if len(ss) == 0 {
		return Empty[T]()
	}
	s := ss[0]
	for i := 1; i < len(ss); i++ {
		s = s.Merge(ss[i])
	}
	return s
}
