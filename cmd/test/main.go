package main

import (
	"ca/pkg/fib"
	"ca/pkg/integer"
	"ca/pkg/padic"
	"ca/pkg/ring"
	"ca/pkg/uint_ntt"
	"fmt"
)

func testPAdic() {
	const N = 10
	x := padic.NewPAdicFromInt(3, 23)
	y := padic.NewPAdicFromInt(3, 27)
	z := padic.NewPAdicFromInt(3, 92)
	fmt.Println(x.Approx(N))
	fmt.Println(y.Approx(N))
	fmt.Println(y.Mul(x).Approx(N)) // print 27 x 23 = 621
	fmt.Println(y.Sub(x).Approx(N)) // print 27 - 23 = 4
	fmt.Println(z.Div(x).Approx(N)) // print 92 / 23 = 4
}

func testUint() {
	x := uint_ntt.FromString("0x318346193417412890342342")
	z := uint_ntt.FromString("0x484723895378245789")
	y := uint_ntt.FromString(x.String())
	fmt.Println(x.Add(z).Mod(y)) // (x + integer) % x
}

func testFib() {
	fmt.Println(fib.Fib(integer.Zero, uint64(20)))
}

func testEA() {
	a, b := ring.EuclideanAlgorithm(integer.FromInt64(15), integer.FromInt64(46))
	fmt.Println(a, b)
}

func main() {
	testUint()
}
