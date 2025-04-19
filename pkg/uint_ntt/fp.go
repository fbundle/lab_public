package uint_ntt

import (
	"ca/pkg/vend/uint128"
	"sync"
)

// p : operations on finite field of order p
const (
	// p, g : p = 2^64 - 2^32 + 1, generator g of F_p^\times, g^{p-1} = 1 mod p
	p, g uint64 = 18446744069414584321, 7
)

var (
	primitiveRootCache     = map[uint64]uint64{}
	primitiveRootCacheLock = sync.Mutex{}
)

func getPrimitiveRoot(n uint64) uint64 {
	if (p-1)%n != 0 {
		panic("n must divide p-1")
	}
	primitiveRootCacheLock.Lock()
	defer primitiveRootCacheLock.Unlock()
	if _, ok := primitiveRootCache[n]; !ok {
		primitiveRootCache[n] = pow(g, (p-1)/n)
	}
	return primitiveRootCache[n]
}

// int : Fermat little theorem a^{p-1} = 1 mod p
func inv(a uint64) uint64 {
	return pow(a, p-2)
}

func add(a uint64, b uint64) uint64 {
	aLarge, bLarge := uint128.From64(a), uint128.From64(b)
	return aLarge.Add(bLarge).Mod64(p)
}

func sub(a uint64, b uint64) uint64 {
	return add(a, p-b%p)
}

func mul(a uint64, b uint64) uint64 {
	aLarge, bLarge := uint128.From64(a), uint128.From64(b)
	return aLarge.Mul(bLarge).Mod64(p)
}

func pow(a uint64, n uint64) uint64 {
	a = a % p

	var x uint64 = 1
	for n > 0 {
		if n%2 == 1 {
			x = mul(x, a)
		}
		a = mul(a, a)
		n /= 2
	}
	return x
}
