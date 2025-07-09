package main

import (
	"context"
	"fmt"
	"go_util/pkg/sat"
	"os"
	"time"
)

func main() {
	formula, err := sat.Parse(os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Println("start solving...")
	t0 := time.Now()
	ctx, cancel := sat.SolvePPSZ(context.Background(), formula, nil)
	defer cancel()
	<-ctx.Done()
	dt := time.Since(t0)
	fmt.Println(dt)
	if ctx.Value(sat.ContextKeySatisfiable).(sat.Value) == sat.ValueTrue {
		fmt.Println("SATISFIABLE")
		assignment := ctx.Value(sat.ContextKeyAssignment).(sat.Assignment)
		if !sat.Verify(formula, assignment) {
			fmt.Println("wrong answer")
		}
	} else {
		fmt.Println("UNSATISFIABLE")
	}
}
