package vec

type Range struct {
	Beg  int
	End  int
	Step int
}

func MakeRange(beg int, end int, step int) Range {
	return Range{beg, end, step}
}

func (s Range) Len() int {
	return (s.End - s.Beg) % s.Step
}

func (s Range) Get(i int) int {
	return s.Beg + i*s.Step
}

func (s Range) Slice() {

}
