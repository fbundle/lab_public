package uint1792

import (
	"fmt"
	"strings"
)

const (
	R = 8 // 8 is 64-th primitive root of unity in mod P
	N = 64
	B = 1 << 28 // pick base B = 2^d, N * B * B < P so that multiplication won't overflow
)

// Block : polynomial of degree at most N-1 in F_p[X]
type Block = [N]uint64

var zeroBlock Block = [N]uint64{}

// Uint1792 : represents nonnegative integers by a_0 + a_1 B + a_2 B^2 + ... + a_{N-1} B^{N-1}
type Uint1792 struct {
	Time Block
}

var Zero Uint1792 = FromUint64(0)
var One Uint1792 = FromUint64(1)

func (a Uint1792) Zero() Uint1792 {
	return Uint1792{Time: zeroBlock}
}

func (a Uint1792) One() Uint1792 {
	return One
}

func FromUint64(x uint64) Uint1792 {
	return fromTime(Block{x})
}

func (a Uint1792) Uint64() uint64 {
	return a.Time[0] + a.Time[1]*B + a.Time[2]*B*B
}

func fromTime(time Block) Uint1792 {
	// trim time
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] = time[i+1] + q
		}
	}
	return Uint1792{
		Time: time,
	}
}

func FromString(s string) Uint1792 {
	if s[0:2] != "0x" {
		panic("string does not start with 0x")
	}
	s = strings.ToLower(s[2:])

	if len(s) > 448 {
		panic("string too long")
	}
	// convert string to base16
	var base16 []byte
	toBase16 := map[string]byte{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
		"a": 10,
		"b": 11,
		"c": 12,
		"d": 13,
		"e": 14,
		"f": 15,
	}
	for i := len(s) - 1; i >= 0; i-- {
		base16 = append(base16, toBase16[string(s[i])])
	}
	// convert base16 to base 2^28
	time := Block{}
	for len(base16)%7 != 0 {
		base16 = append(base16, 0)
	}
	for i := 0; i < len(base16)/7; i++ {
		var x uint64 = 0
		var b uint64 = 1
		for j := 0; j < 7; j++ {
			x += uint64(base16[7*i+j]) * b
			b *= 16
		}
		time[i] = x
	}

	return fromTime(time)
}

func (a Uint1792) String() string {
	// convert base 2^28 to base16
	var base16 []byte = nil
	for i := 0; i < N; i++ {
		x := a.Time[i]
		for j := 0; j < 7; j++ {
			q, r := x/16, x%16
			base16 = append(base16, byte(r))
			x = q
		}
		if x > 0 {
			panic("wrong")
		}
	}
	// convert base16 to string
	toChar := map[byte]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "a",
		11: "b",
		12: "c",
		13: "d",
		14: "e",
		15: "f",
	}
	out := ""
	if len(base16)%2 != 0 {
		base16 = append(base16, byte(0))
	}
	for i := len(base16) - 1; i >= 0; i-- {
		ch := toChar[base16[i]]
		out += ch
	}
	out = strings.TrimLeft(out, "0")
	if len(out) == 0 {
		out = "0"
	}
	return "0x" + out
}

func (a Uint1792) Add(b Uint1792) Uint1792 {
	time := Block{}
	for i := 0; i < N; i++ {
		time[i] = add(a.Time[i], b.Time[i])
	}
	return fromTime(time)
}

func (a Uint1792) Mul(b Uint1792) Uint1792 {
	aFreq, bFreq := time2freq(a.Time), time2freq(b.Time)

	freq := Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(aFreq[i], bFreq[i])
	}
	return fromTime(freq2time(freq))
}

func (a Uint1792) Sub(b Uint1792) Uint1792 {
	return a.Add(b.Neg())
}

func (a Uint1792) Neg() Uint1792 {
	// second complement
	time := Block{}
	for i := 0; i < N; i++ {
		time[i] = (^a.Time[i]) % B // flip bits and trim to B
	}
	return fromTime(time).Add(FromUint64(1))
}

// IsNeg : treat Uint1792 as Int1792
func (a Uint1792) IsNeg() bool {
	return a.Time[N-1]/(1<<27) > 0
}

func (a Uint1792) Abs() Uint1792 {
	if a.IsNeg() {
		return FromUint64(0).Sub(a)
	} else {
		return a
	}
}

func (a Uint1792) Sign() int {
	if a.IsNeg() {
		return -1
	}
	if a.Time == zeroBlock {
		return 0
	} else {
		return 1
	}

}

// shiftRight : a -> a / 2^{28 n}
func (a Uint1792) shiftRight(n int) Uint1792 {
	time := Block{}
	for i := 0; i < N; i++ {
		if 0 <= i+n && i+n < N {
			time[i] = a.Time[i+n]
		}
	}
	return fromTime(time)
}

// inv : approx the root of f(x) = m / x - a for m = 2^{28 n}
// using fixed point arithmetic
func (a Uint1792) inv(n int) Uint1792 {
	if a.shiftRight(n).Time != zeroBlock {
		panic(fmt.Sprintf("inv is supported only for 1 < a < 2^{28 %d}", n))
	}
	if a.Uint64() == 0 {
		panic("division by zero")
	}

	two := FromUint64(2)
	x := FromUint64(1)
	// Newton iteration
	// x_{n+1} = 2 x_n - (a x_n^2) / m
	for {
		// a / m = a.shiftRight(N/2)
		x1 := two.Mul(x).Sub(a.Mul(x).Mul(x).shiftRight(n))
		if x1 == x {
			break
		}
		x = x1
	}

	return x
}

func (a Uint1792) Div(b Uint1792) Uint1792 {
	x := b.inv(N / 2)                 // x ~ 2^896 / b
	return a.Mul(x).shiftRight(N / 2) // a/b = a (2^896 / b) / 2^896
}

func (a Uint1792) Mod(b Uint1792) Uint1792 {
	x := a.Div(b)
	return a.Sub(b.Mul(x))
}

func (a Uint1792) DivMod(b Uint1792) (Uint1792, Uint1792) {
	q := a.Div(b)
	r := a.Sub(b.Mul(q))
	return q, r
}

const (
	invR = 16140901060737761281 // precompute R^{-1}
	invN = 18158513693329981441 // precompute N^{-1}
	S    = 18410715272404008961 // precompute N^{-1} R^{-1}
)

func time2freq(time Block) Block {
	return dft(time, N, R)
}

func freq2time(freq Block) Block {
	time := Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = mul(f, invN)
	}
	// TODO - check why dft(time, N, mul(invR, invS)) seems to not work
	return time
}
