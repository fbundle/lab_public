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
	a := func() monad.Monad[int] {
		return monad.None[int]().Prepend(1, 2, 3, 4)
	}
	resultList := []interface{}{
		monad.Natural().TakeAtMost(10).Slice(),
		monad.Filter(monad.Natural(), func(n int) bool {
			return n%2 == 0
		}).TakeAtMost(10).Slice(),
		monad.Replicate(5).TakeAtMost(10).Slice(),
		a().TakeAtMost(0).Slice(),
		a().TakeAtMost(2).Slice(),
		a().TakeAtMost(5).Slice(),
		a().Prepend(9, 8, 7).Slice(),
		a().DropAtMost(0).Slice(),
		a().DropAtMost(2).Slice(),
		a().DropAtMost(5).Slice(),
		monad.Map(a(), func(x int) int {
			return x * 2
		}).Slice(),
		monad.Filter(a(), func(n int) bool {
			return n%2 == 0
		}).Slice(),
		monad.Reduce(a(), func(tr string, t int) string {
			return fmt.Sprintf("%s%d,", tr, t)
		}, ""), // TODO fix reduce
		monad.Fold(a(), func(tr string, t int) string {
			return fmt.Sprintf("%s%d,", tr, t)
		}, "").Slice(),
		monad.Bind(a(), func(ta int) monad.Monad[int] {
			return monad.Replicate(ta).TakeAtMost(ta)
		}).Slice(),
	}
	for _, result := range resultList {
		fmt.Println(result)
	}
}

func main() {
	testMonad()
}
