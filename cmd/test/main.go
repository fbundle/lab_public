package main

import (
	"ca/pkg/mp"
	"ca/pkg/pa"
	"fmt"
)

func testPAdic() {
	const N = 10
	x := pa.NewPAdicFromInt(3, 23)
	y := pa.NewPAdicFromInt(3, 27)
	z := pa.NewPAdicFromInt(3, 92)
	fmt.Println(x.Approx(N))
	fmt.Println(y.Approx(N))
	fmt.Println(y.Mul(x).Approx(N)) // print 27 x 23 = 621
	fmt.Println(y.Sub(x).Approx(N)) // print 27 - 23 = 4
	fmt.Println(z.Div(x).Approx(N)) // print 92 / 23 = 4
}

func main() {
	a := mp.NewUint1024FromUint64(1231231123121)
	b := mp.NewUint1024FromUint64(6345353645645)
	_, _ = a, b
	fmt.Println(a)
	fmt.Println(mp.NewUint1024FromFreq(a.Freq))
	mp.TestDft()
}
