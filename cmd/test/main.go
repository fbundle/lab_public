package main

import (
	"ca/pkg/ca"
	"fmt"
)

func main() {
	x := ca.NewPArdicFromList(3, []int{2})
	y := ca.NewPArdicFromList(3, []int{1})
	_ = x
	fmt.Println(y.Neg().Add(x).Approx(10))
}
