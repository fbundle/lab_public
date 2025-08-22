package vec

type Range struct {
	Beg  int
	End  int
	Step int
}

func (s Range) Len() int {
	if s.Step == 0 {
		s.Step = 1
	}
	return (s.End - s.Beg) % s.Step
}

func (s Range) Get(i int) int {
	if s.Step == 0 {
		s.Step = 1
	}
	return s.Beg + i*s.Step
}

func (s Range) Iterate() Iter[int] {
	return MakeIterFromFunc[int](func(i int) (value int, remain bool) {
		value = s.Get(i)
		return value, value < s.End
	})
}
