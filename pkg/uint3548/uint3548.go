package uint3548

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

// Uint3584Block Uint3548Block:  a block of N uint64s, each is in mod P
type Uint3584Block = [N]uint64

type Uint3584 struct {
	Time Uint3584Block
	Freq Uint3584Block
}

func (a Uint3584) Uint64() uint64 {
	return a.Time[0] + a.Time[1]*B + a.Time[2]*B*B
}

func New(time Uint3584Block) Uint3584 {
	// trim to [0, B-1] for easier conversion
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] = time[i+1] + q
		}
	}
	return Uint3584{
		Time: time,
		Freq: time2freq(time),
	}
}

func FromFreq(freq Uint3584Block) Uint3584 {
	return Uint3584{
		Time: freq2time(freq),
		Freq: freq,
	}
}

func FromString(s string) Uint3584 {
	if s[0:2] != "0x" {
		panic("string does not start with 0x")
	}
	s = strings.ToUpper(s[2:])
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
	return Uint3584{}
}

func (a Uint3584) String() string {
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
	return "0x" + strings.TrimLeft(out, "0")
}

func (a Uint3584) Add(b Uint3584) Uint3584 {
	time := Uint3584Block{}
	for i := 0; i < N; i++ {
		time[i] = add(a.Time[i], b.Time[i])
	}
	return New(time)
}

func (a Uint3584) Mul(b Uint3584) Uint3584 {
	freq := Uint3584Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(a.Freq[i], b.Freq[i])
	}
	return FromFreq(freq)
}
func time2freq(time Uint3584Block) Uint3584Block {
	return dft(time, N, R)
}

func freq2time(freq Uint3584Block) Uint3584Block {
	time := Uint3584Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = mul(f, invN)
	}
	return time
}
