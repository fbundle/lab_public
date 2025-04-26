package sat_test

import (
	"context"
	"fmt"
	"github.com/khanh-nguyen-code/go_util/pkg/sat"
	"math/rand"
	"testing"
	"time"
)

func TestCDCL(t *testing.T) {
	// x1 true, x2 false
	// (x1 or x2) and (not x2 or x3) and not x3
	formula := [][]int{
		{1, 2},
		{-2, 3},
		{-3},
	}
	ctx, cancel := sat.SolveCDCL(context.Background(), formula, nil)
	defer cancel()
	<-ctx.Done()
	fmt.Println(ctx.Value(sat.ContextKeySatisfiable))
	fmt.Println(ctx.Value(sat.ContextKeyAssignment))
}

func TestCDCLTimeout(t *testing.T) {
	numVariable := 1000
	numClause := 4 * numVariable
	var formula [][]int
	for i := 0; i < numClause; i++ {
		v1, v2, v3 := rand.Intn(numVariable)+1, rand.Intn(numVariable)+1, rand.Intn(numVariable)+1
		s1, s2, s3 := 1, 1, 1
		if rand.Intn(2) == 0 {
			s1 = -1
		}
		if rand.Intn(2) == 0 {
			s2 = -1
		}
		if rand.Intn(2) == 0 {
			s3 = -1
		}
		formula = append(formula, []int{v1 * s1, v2 * s2, v3 * s3})
	}
	timeout, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ctx, cancel := sat.SolveCDCL(timeout, formula, nil)
	defer cancel()
	<-ctx.Done()
	fmt.Println(ctx.Value(sat.ContextKeySatisfiable))
	fmt.Println(ctx.Value(sat.ContextKeyAssignment))
}
