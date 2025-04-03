package main

import (
	"ca/pkg/ca"
	"fmt"
)

const N = 10

func main() {
	one := ca.NewPArdicFromList(3, []int{1}, nil)
	two := ca.NewPArdicFromList(3, []int{2}, nil)
	x := ca.NewPArdicFromList(3, []int{2, 1}, nil) // 5
	y := ca.NewPArdicFromList(3, []int{1, 2}, nil) // 7
	_, _, _, _ = x, y, one, two
	fmt.Println(y.Mul(x).Approx(N))
}
