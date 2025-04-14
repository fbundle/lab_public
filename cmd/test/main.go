package main

import (
	"ca/pkg/pa"
	"ca/pkg/uint3548"
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
	x := uint3548.TrimFromTime(uint3548.Uint3584Block{131414234243, 12314555, 123131, 5345777, 7646456})
	s := x.String()
	fmt.Println(s)
	y := uint3548.FromString(s)
	t := y.String()
	fmt.Println(t)
}
