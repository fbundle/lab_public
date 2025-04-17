package integer

import "math/big"

type Integer struct {
	bigint *big.Int
}

var Zero Integer = Integer{big.NewInt(0)}

var One Integer = Integer{big.NewInt(1)}

func (a Integer) Zero() Integer {
	return Zero
}

func (a Integer) One() Integer {
	return One
}

func (a Integer) String() string {
	return "0x" + a.bigint.Text(16)
}

func (a Integer) Add(b Integer) Integer {
	return Integer{(&big.Int{}).Add(a.bigint, b.bigint)}
}

func (a Integer) Sub(b Integer) Integer {
	return Integer{(&big.Int{}).Sub(a.bigint, b.bigint)}
}

func (a Integer) Mul(b Integer) Integer {
	return Integer{(&big.Int{}).Mul(a.bigint, b.bigint)}
}

func (a Integer) Neg() Integer {
	return Integer{(&big.Int{}).Neg(a.bigint)}
}

func (a Integer) Div(b Integer) Integer {
	return Integer{(&big.Int{}).Div(a.bigint, b.bigint)}
}

func (a Integer) Mod(b Integer) Integer {
	return Integer{(&big.Int{}).Mod(a.bigint, b.bigint)}
}

func (a Integer) DivMod(b Integer) (Integer, Integer) {
	q, r := (&big.Int{}).DivMod(a.bigint, b.bigint, &big.Int{})
	return Integer{q}, Integer{r}
}

func (a Integer) Cmp(b Integer) int {
	return a.bigint.Cmp(b.bigint)
}

func (a Integer) Equal(b Integer) bool {
	return a.bigint.Cmp(b.bigint) == 0
}
