package uint3548

import "ca/pkg/uint128"

const (
	P = 18446744069414584321 // p = 2^64 - 2^32 + 1 and 2 is the 192-th primitive root of unity
)

func add(a uint64, b uint64) uint64 {
	aLarge, bLarge := uint128.From64(a), uint128.From64(b)
	return aLarge.Add(bLarge).Mod64(P)
}

func sub(a uint64, b uint64) uint64 {
	return add(a, P-b%P)
}

func mul(a uint64, b uint64) uint64 {
	aLarge, bLarge := uint128.From64(a), uint128.From64(b)
	return aLarge.Mul(bLarge).Mod64(P)
}

func pow(a uint64, n uint64) uint64 {
	var powLarge func(a uint128.Uint128, n uint64) uint64
	powLarge = func(a uint128.Uint128, n uint64) uint64 {
		if n >= P {
			return powLarge(a, n%P)
		}
		if n == 0 {
			return 1
		}
		if n == 1 {
			return a.Mod64(P)
		}
		if n%2 == 0 {
			h := uint128.From64(powLarge(a, n/2))
			return h.Mul(h).Mod64(P)
		} else {
			h := uint128.From64(powLarge(a, n/2))
			return uint128.From64(h.Mul(h).Mod64(P)).Mul(a).Mod64(P)
		}
	}
	return powLarge(uint128.From64(a), n)
}
