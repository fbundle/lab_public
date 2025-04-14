package uint1792

import (
	"strings"
)

const (
	R, N = 8, 64   // 8 is 64-th primitive root of unity in mod P
	B    = 1 << 28 // choose base 2^d so that N * B * B < P - this guarantees that multiplication won't overflow
)
const (
	invR = 16140901060737761281 // precompute R^{-1}
	invN = 18158513693329981441 // precompute N^{-1}
)

// Uint1792Block Uint3548Block:  a block of N uint64s, each is in mod P
type Uint1792Block = [N]uint64

// Uint1792 : represents nonnegative integers by a_0 + a_1 B + a_2 B^2 + ...
type Uint1792 struct {
	Time Uint1792Block
}

func FromUint64(x uint64) Uint1792 {
	return fromTime(Uint1792Block{x})
}

func (a Uint1792) Uint64() uint64 {
	return a.Time[0] + a.Time[1]*B + a.Time[2]*B*B
}

func fromTime(time Uint1792Block) Uint1792 {
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
	time := Uint1792Block{}
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
	time := Uint1792Block{}
	for i := 0; i < N; i++ {
		time[i] = add(a.Time[i], b.Time[i])
	}
	return fromTime(time)
}

func (a Uint1792) Mul(b Uint1792) Uint1792 {
	aFreq, bFreq := time2freq(a.Time), time2freq(b.Time)

	freq := Uint1792Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(aFreq[i], bFreq[i])
	}
	return fromTime(freq2time(freq))
}

func (a Uint1792) Sub(b Uint1792) Uint1792 {
	// second complement for b
	bTime := b.Time
	for i := 0; i < N; i++ {
		bTime[i] = (^bTime[i]) % B // flip bits and trim to B
	}
	bNeg := fromTime(bTime).Add(FromUint64(1))
	return a.Add(bNeg)
}

// Abs : treat Uint1792 as a signed integer
func (a Uint1792) Abs() Uint1792 {
	if a.Time[N-1]/(1<<27) > 0 {
		return FromUint64(0).Sub(a)
	} else {
		return a
	}
}

func (a Uint1792) ShiftLeft(n int) Uint1792 {
	time := Uint1792Block{}
	for i := 0; i < N; i++ {
		if 0 <= i-n && i-n < N {
			time[i] = a.Time[i-n]
		}
	}
	return fromTime(time)
}

func (a Uint1792) ShiftRight(n int) Uint1792 {
	return a.ShiftLeft(-n)
}

// Inv : Newton-Raphson Division for 1 < a < 2^896
// return x so that ax = 2^896 = m using fixed point arithmetic
func (a Uint1792) Inv() Uint1792 {
	// a / m = a.ShiftRight(N/2)
	zero := Uint1792Block{}
	if a.ShiftRight(N/2).Time != zero {
		panic("inv if only a < 2^896")
	}
	if a.Uint64() == 0 {
		panic("division by zero")
	}

	two := FromUint64(2)
	x := FromUint64(1)
	//Newton-Raphson iterations
	// x_{n+1} = 2 x_n - (a x_n^2) / m
	for {
		x1 := two.Mul(x).Sub(a.Mul(x).Mul(x).ShiftRight(N / 2))
		if x1 == x {
			break
		}
		x = x1
	}

	return x
}

func (a Uint1792) Div(b Uint1792) Uint1792 {
	x := b.Inv() // x = 2^896 / b
	return a.Mul(x).ShiftRight(896 / 28)
}

func time2freq(time Uint1792Block) Uint1792Block {
	return dft(time, N, R)
}

func freq2time(freq Uint1792Block) Uint1792Block {
	time := Uint1792Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = mul(f, invN)
	}
	return time
}
