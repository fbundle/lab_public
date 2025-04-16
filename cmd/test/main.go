package main

import (
	"ca/pkg/pa"
	"ca/pkg/uint1792"
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

func testUint1792() {
	x := uint1792.FromString("0x318346193417412890342342")
	s := x.String()
	y := uint1792.FromString(s)
	fmt.Println(x.Add(uint1792.FromString("0x934174128903")).Mod(y))
}

func main() {
	testUint1792()
}
