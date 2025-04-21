package main

import (
	"ca/pkg/fib"
	"ca/pkg/integer"
	"ca/pkg/padic"
	"ca/pkg/ring"
	"ca/pkg/uint_ntt"
	"ca/pkg/vec"
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

	_ = fib.Fib(uint_ntt.UintNTT{}, uint64(n))
}

func testEA() {
	a, b := ring.EuclideanAlgorithm(integer.FromInt64(15), integer.FromInt64(46))
	fmt.Println(a, b)
}

func testUintNTT() {
	x := uint_ntt.FromString("0x318346193417412890342342")
	z := uint_ntt.FromString("0x539543980a084524")
	fmt.Println(z)
	fmt.Println(z.Add(x).Mod(x))
}

func testVec() {
	i := vec.Map[int, int](vec.Range{Beg: 0, End: 64}.Iterate(), func(x int) (y int) {
		return 64 - x
	})
	i, v := vec.ViewIter(i)
	fmt.Println(v)
	i = vec.Filter[int](i, func(v int) bool {
		return v%2 == 0
	})
	i, v = vec.ViewIter(i)
	fmt.Println(v)
	z := vec.Reduce(i, func(i int, j int, x int, y int) int {
		return x + y
	})
	fmt.Println(z)
}

func main() {
	testVec()
}
