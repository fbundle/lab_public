package sat

type Variable = int

type Literal = int

type Clause = []Literal

type Formula []Clause

func (formula Formula) NumClause() int {
	return len(formula)
}

func (formula Formula) NumVariable() int {
	numVar := 0
	for _, clause := range formula {
		for _, literal := range clause {
			variable := abs(literal)
			if numVar < variable {
				numVar = variable
			}
		}
	}
	return numVar
}
