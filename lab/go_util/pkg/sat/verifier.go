package sat

func Verify(formula Formula, assignment Assignment) bool {
	s := &bcpState{
		Formula:    formula,
		Assignment: assignment,
	}
	return s.Verify()
}
