package main

import (
	"ca/pkg/fib"
	"ca/pkg/integer"
	"ca/pkg/padic"
	"ca/pkg/ring"
	"ca/pkg/uint1792"
	"ca/pkg/uint_ntt"
	"fmt"
	"os"
	"strconv"
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

func testUint1792() {
	x := uint1792.FromString("0x318346193417412890342342")
	z := uint1792.FromString("0x484723895378245789")
	y := uint1792.FromString(x.String())
	fmt.Println(x.Add(z).Mod(y)) // (x + integer) % x
}

func testFib() {
	if len(os.Args) < 2 {
		panic("Usage: go run main.go <integer>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Error: argument must be an integer")
		return
	}

	fmt.Println(fib.Fib(uint_ntt.UintNTT{}, uint64(n)))
}

func testEA() {
	a, b := ring.EuclideanAlgorithm(integer.FromInt64(15), integer.FromInt64(46))
	fmt.Println(a, b)
}

func testUintNTT() {
	x := uint_ntt.FromString("0x318346193417412890342342")
	y := uint_ntt.FromString(x.String())
	fmt.Println(x, y)
	fmt.Println(x.Add(y))
	fmt.Println(x.Mul(y))
	fmt.Println(x.Sub(y))
	z := uint_ntt.FromString("0x539543980a084524")
	fmt.Println(z.Add(x).Mod(x))
}

func main() {
	testUintNTT()
}
