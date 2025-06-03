package main

import (
	"ca/pkg/fib"
	"ca/pkg/int_ntt"
	"ca/pkg/integer"
	"ca/pkg/monad"
	"ca/pkg/padic"
	"ca/pkg/ring"
	"ca/pkg/tup"
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

	fmt.Println(fib.Fib(int_ntt.Nat{}, uint64(n)))
}

func testEA() {
	a, b := ring.EuclideanAlgorithm(integer.FromInt64(15), integer.FromInt64(46))
	fmt.Println(a, b)
}

func testIntNTT() {
	x := int_ntt.FromString("0x318346193417412890342342")
	z := int_ntt.FromString("0x539543980a084524")
	fmt.Println(z)
	fmt.Println(z.Add(x).Mod(x))
}

func testTup() {
	t := tup.MakeTup(2, 2.5)
	fmt.Println(t)
}

func testVecFunctor() {
	add1 := func(x int) int {
		return x + 1
	}

	v1 := vec.MakeVecFromSlice([]int{1, 2, 3, 4, 5})

	add1vec := vec.Wrap(add1)
	v2 := add1vec(v1)
	fmt.Println(v2)

}

func testMonad() {
	bind := func(x int) monad.Monad[int] {
		if x%2 == 0 {
			return monad.None[int]()
		} else {
			return monad.FromSlice([]int{x, x + 1, x + 2})
		}
	}
	m := monad.FromSlice([]int{1, 2, 3})
	m = monad.Bind(m, bind)
	s := monad.ToSlice(m)
	fmt.Println(s)
	m = monad.None[int]()
	m = monad.Bind(m, bind)
	s = monad.ToSlice(m)
	fmt.Println(s)
}

func main() {
	testMonad()
}
