package iter

type Morphism = func(x any) (any, bool)

type Iterator interface {
	Next() (x any, ok bool)
}

type sliceIterator struct {
	i int
	s []any
}

func (s *sliceIterator) Next() (x any, ok bool) {
	if s.i >= len(s.s) {
		return nil, false
	}
	s.i++
	return s.s[s.i-1], true
}

func MakeIteratorFromSlice(s []any) Iterator {
	return &sliceIterator{
		i: 0,
		s: s,
	}
}

func MakeSliceFromIterator(i Iterator) []any {
	var s []any
	for {
		x, ok := i.Next()
		if !ok {
			break
		}
		s = append(s, x)
	}
	return s
}

func MakeIteratorMore(i Iterator) IteratorMore {
	return &iteratorMore{
		i:  i,
		ms: nil,
	}
}

type IteratorMore interface {
	Iterator
	Apply(Morphism) IteratorMore
}

type iteratorMore struct {
	i  Iterator
	ms []Morphism
}

func (im *iteratorMore) Next() (any, bool) {
	for {
		x, ok := im.i.Next()
		if !ok {
			return nil, false
		}
		x, ok = func(y any) (any, bool) {
			var ok bool
			for _, m := range im.ms {
				x, ok = m(x)
				if !ok {
					return nil, false
				}
			}
			return x, true
		}(x)
		if !ok {
			continue
		}
		return x, true
	}

}

func (im *iteratorMore) Apply(f Morphism) IteratorMore {
	im.ms = append(im.ms, f)
	return im
}
