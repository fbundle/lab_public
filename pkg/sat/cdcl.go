package sat

import (
	"context"
	"github.com/irifrance/gini"
	"github.com/irifrance/gini/z"
)

func lit2zLit(l Literal) z.Lit {
	switch {
	case l > 0:
		return z.Var(abs(l)).Pos()
	case l < 0:
		return z.Var(abs(l)).Neg()
	default:
		panic("l must not be 0")
	}
}
func SolveCDCL(parentCtx context.Context, formula Formula, assumption Assignment) (context.Context, func()) {
	ctx, cancel := context.WithCancel(parentCtx)
	c := &solverCtx{
		ctx: ctx,
		r:   ValueUnknown,
		a:   nil,
	}
	go func() {
		defer cancel()
		g := gini.NewVc(formula.NumVariable(), formula.NumClause())
		for _, clause := range formula {
			if len(clause) == 0 {
				continue
			}
			for _, lit := range clause {
				g.Add(lit2zLit(lit))
			}
			g.Add(z.LitNull)
		}
		assumptionLitList := make([]z.Lit, 0)
		for v := 1; v < len(assumption); v++ {
			if assumption[v] == ValueUnknown {
				continue
			}
			assumptionLitList = append(assumptionLitList, lit2zLit(v*assumption[v]))
		}
		g.Assume(assumptionLitList...)
		s := g.GoSolve()
		go func() {
			<-ctx.Done()
			s.Stop()
		}()
		switch s.Wait() {
		case 1:
			c.r = ValueTrue
		case -1:
			c.r = ValueFalse
		case 0:
			c.r = ValueUnknown
		default:
			panic("unexpected output from gini")
		}
		if c.r == ValueTrue {
			c.a = NewAssignment(formula.NumVariable())
			for v := 1; v < len(c.a); v++ {
				if g.Value(lit2zLit(v)) {
					c.a[v] = ValueTrue
				} else {
					c.a[v] = ValueFalse
				}
			}
		}
	}()
	return c, cancel
}
