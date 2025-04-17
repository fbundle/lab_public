package integer

import "math/big"

type Int struct {
	bigint *big.Int
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
	return "0x" + a.bigint.Text(16)
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
