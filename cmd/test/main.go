package main

import (
	"ca/pkg/ca"
	"fmt"
)

const N = 10

func main() {
	x := ca.NewPArdicFromInt(3, 23)
	y := ca.NewPArdicFromInt(3, 27)
	fmt.Println(x.Approx(N))
	fmt.Println(y.Approx(N))
	fmt.Println(y.Mul(x).Approx(N)) // print 23 x 27
}
