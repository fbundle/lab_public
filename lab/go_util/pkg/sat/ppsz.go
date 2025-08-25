package sat

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
)

func SolvePPSZ(parentCtx context.Context, formula Formula, assumption Assignment) (context.Context, func()) {
	ctx, cancel := context.WithCancel(parentCtx)
	c := &solverCtx{
		ctx: ctx,
		r:   ValueUnknown,
		a:   nil,
	}
	r := rand.New(rand.NewSource(1234))
	mu := &sync.Mutex{}
	concurrent := runtime.NumCPU()
	for j := 0; j < concurrent; j++ {
		go func() {
			defer cancel()
			numVariable := formula.NumVariable()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				s := &bcpState{
					Formula:    formula,
					Assignment: NewAssignment(numVariable),
				}
				if assumption != nil {
					copy(s.Assignment, assumption)
				}
				for {
					// bcp
					haveChance, _ := s.Next()
					if !haveChance {
						break
					}
					// check success
					var zeroVariableList []int
					for i := 1; i < len(s.Assignment); i++ {
						if s.Assignment[i] == ValueUnknown {
							zeroVariableList = append(zeroVariableList, i)
						}
					}
					if len(zeroVariableList) == 0 {
						c.r = ValueTrue
						c.a = s.Assignment
						return
					}
					// guess
					mu.Lock()
					guess := zeroVariableList[r.Intn(len(zeroVariableList))]
					if r.Float32() < 0.5 {
						guess *= -1
					}
					mu.Unlock()
					s.Assignment[abs(guess)] = sign(guess)
				}
			}
		}()
	}
	return c, cancel
}
