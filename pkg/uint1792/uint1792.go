package uint1792

import (
	"fmt"
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
	Time    Uint1792Block
	Freq    Uint1792Block
	hasFreq bool
}

func FromUint64(x uint64) Uint1792 {
	return fromTime(trimTime(Uint1792Block{x}))
}

func (a Uint1792) Uint64() uint64 {
	return a.Time[0] + a.Time[1]*B + a.Time[2]*B*B
}

func trimTime(time Uint1792Block) Uint1792Block {
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] = time[i+1] + q
		}
	}
	return time
}

func fromTime(time Uint1792Block) Uint1792 {
	return Uint1792{
		Time:    time,
		Freq:    Uint1792Block{},
		hasFreq: false,
	}
}

func FromString(s string) Uint1792 {
	if s[0:2] != "0x" {
		panic("string does not start with 0x")
	}
	s = strings.ToUpper(s[2:])

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
		"A": 10,
		"B": 11,
		"C": 12,
		"D": 13,
		"E": 14,
		"F": 15,
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
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
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
	return fromTime(trimTime(time))
}

func (a Uint1792) Mul(b Uint1792) Uint1792 {
	aFreq := a.Freq
	if !a.hasFreq {
		aFreq = time2freq(a.Time)
	}
	bFreq := b.Freq
	if !b.hasFreq {
		bFreq = time2freq(b.Time)
	}
	freq := Uint1792Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(aFreq[i], bFreq[i])
	}
	time := trimTime(freq2time(freq))
	return fromTime(time)
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

// Inv : TODO - Newton-Raphson Division
// - return x so that ax = 2^1792 = m using fixed point arithmetic
// implement using Uint3584 add sub mul
// after that, implement Div and Mod
func (a Uint1792) Inv() Uint1792 {
	x := a
	for {
		// x = x * (2 - a * x)
		x1 := x.Mul(FromUint64(2).Sub(a.Mul(x)))
		if x == x1 {
			break
		}
		fmt.Println(x.Sub(x1).Abs())
		x = x1
	}

	return x
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
