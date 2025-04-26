package sat

type bcpState struct {
	Formula    Formula
	Assignment Assignment
}

func (s *bcpState) Verify() bool {
	for _, clause := range s.Formula {
		if value, _, _ := s.clauseValue(clause); value != ValueTrue {
			return false
		}
	}
	return true
}

// Next : perform bcp Boolean Constraint propagation
func (s *bcpState) Next() (haveChance bool, activated []Literal) {
	for {
		unitProp := false
		for _, clause := range s.Formula {
			value, zeroCount, zeroFirstIdx := s.clauseValue(clause)
			if value == ValueFalse {
				return false, activated
			}
			if value == ValueTrue {
				continue
			}
			if value == ValueUnknown && zeroCount == 1 {
				activateLiteral := clause[zeroFirstIdx]
				s.Assignment[abs(activateLiteral)] = sign(activateLiteral)
				activated = append(activated, activateLiteral)
				unitProp = true
			}
		}
		if !unitProp {
			break
		}
	}
	return true, activated
}

func (s *bcpState) literalValue(literal Literal) Value {
	return s.Assignment[abs(literal)] * sign(literal)
}

// clauseValue:
// value : 1 if clause is sat, -1 if clause is unsat, 0 if unknown
// zeroCount : number of unknown literals (only when value=0)
// zeroFirstIdx : index of the first unknown literal (only when value=0)
func (s *bcpState) clauseValue(clause Clause) (value Value, zeroCount int, zeroFirstIdx int) {
	zeroFirstIdx = -1
	for idx, literal := range clause {
		v := s.literalValue(literal)
		if v == ValueTrue {
			return ValueTrue, 0, 0
		}
		if v == ValueUnknown {
			zeroCount++
			if zeroFirstIdx == -1 {
				zeroFirstIdx = idx
			}
		}
	}
	if zeroCount > 0 {
		return ValueUnknown, zeroCount, zeroFirstIdx
	}
	return ValueFalse, 0, 0
}
