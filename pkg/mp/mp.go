package mp

const (
	BASE = 1 << 16
	// P
	// p = \phi_6(x) = x^2 - x + 1 where x = 2^32
	// then in F_p, 2^32 is the 6th root of unity
	// or 2 is the 192 th root of unity in F_p
	P = 18446744069414584321
	R = 2
	N = 192
)

type Uint3056 struct {
	time [N]uint16
	freq [N]uint64
}
