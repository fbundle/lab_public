package integer

import "math/big"

type Int struct {
	bigint *big.Int
}

func FromInt64(x int64) Int {
	return Int{bigint: big.NewInt(x)}
}

var Zero Int = Int{big.NewInt(0)}

var One Int = Int{big.NewInt(1)}

func (a Int) Zero() Int {
	return Zero
}

func (a Int) One() Int {
	return One
}

func (a Int) String() string {
	text := a.bigint.Text(16)
	if text[0] == '-' {
		return "-0x" + text[1:]
	} else {
		return "0x" + text
	}
}

func (a Int) Add(b Int) Int {
	return Int{(&big.Int{}).Add(a.bigint, b.bigint)}
}

func (a Int) Sub(b Int) Int {
	return Int{(&big.Int{}).Sub(a.bigint, b.bigint)}
}

func (a Int) Mul(b Int) Int {
	return Int{(&big.Int{}).Mul(a.bigint, b.bigint)}
}

func (a Int) Neg() Int {
	return Int{(&big.Int{}).Neg(a.bigint)}
}

func (a Int) Div(b Int) Int {
	return Int{(&big.Int{}).Div(a.bigint, b.bigint)}
}

func (a Int) Mod(b Int) Int {
	return Int{(&big.Int{}).Mod(a.bigint, b.bigint)}
}

func (a Int) DivMod(b Int) (Int, Int) {
	q, r := (&big.Int{}).DivMod(a.bigint, b.bigint, &big.Int{})
	return Int{q}, Int{r}
}

func (a Int) Cmp(b Int) int {
	return a.bigint.Cmp(b.bigint)
}

func (a Int) Equal(b Int) bool {
	return a.bigint.Cmp(b.bigint) == 0
}

func (a Int) Norm() Int {
	return Int{(&big.Int{}).Abs(a.bigint)}
}
