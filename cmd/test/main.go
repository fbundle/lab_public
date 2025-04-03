package main

import (
	"ca/pkg/ca"
	"fmt"
)

const N = 10

func main() {
	x := ca.NewPArdicFromList(3, []int{2, 1}, nil) // 5
	y := ca.NewPArdicFromList(3, []int{1, 2}, nil) // 7
	fmt.Println(x.Approx(N))
	fmt.Println(y.Approx(N))
	fmt.Println(y.Mul(x).Approx(N)) // print 35
}
