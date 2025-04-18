package uint_ntt

import (
	"math/bits"
	"strings"
)

const (
	B = 1 << 16 // pick base B = 2^d, N * B * B < p so that multiplication won't overflow
)

// TODO - set N variable and find R on the fly

// Block : polynomial in F_p[X]
type Block []uint64

func (b Block) get(i int) uint64 {
	if i >= len(b) {
		return 0
	}
	return b[i]
}

var zeroBlock []uint64

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
	var q, r uint64 = 0, 0
	for i := 0; i < len(time); i++ {
		q, r = time[i]/B, time[i]%B
		time[i] = r
		if i+1 < len(time) {
			time[i+1] = time[i+1] + q
		}
	}
	time = append(time, q)
	for len(time) > 0 && time[len(time)-1] == 0 {
		time = time[:len(time)-1]
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
	// convert base16 (2^4) to base 2^16 then trim
	for len(base16)%4 != 0 {
		base16 = append(base16, byte(0))
	}

	time := Block{}
	for i := 0; i < len(base16); i += 4 {
		var x uint64 = 0
		x += uint64(base16[i])
		x += uint64(base16[i+1]) * 16
		x += uint64(base16[i+2]) * 16 * 16
		x += uint64(base16[i+3]) * 16 * 16 * 16
		time = append(time, x)
	}

	return fromTime(time)
}

func (a Uint1792) String() string {
	// convert base 2^16 to base16 (2^4)
	var base16 []byte = nil
	for i := 0; i < len(a.Time); i++ {
		x := a.Time[i]
		base16 = append(base16, byte(x%16))
		x /= 16
		base16 = append(base16, byte(x%16))
		x /= 16
		base16 = append(base16, byte(x%16))
		x /= 16
		base16 = append(base16, byte(x%16))
		x /= 16
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
	l := max(len(a.Time), len(b.Time))
	time := make(Block, l)
	for i := 0; i < l; i++ {
		time[i] = add(a.Time.get(i), b.Time.get(i))
	}
	return fromTime(time)
}

func (a Uint1792) Mul(b Uint1792) Uint1792 {
	aFreq, bFreq := time2freq(a.Time), time2freq(b.Time)
	l := min(len(aFreq), len(bFreq))
	freq := Block{}
	for i := 0; i < l; i++ {
		freq = append(freq, mul(aFreq.get(i), bFreq.get(i)))
	}
	return fromTime(freq2time(freq))
}

func nextPowerOfTwo(x uint64) uint64 {
	if x == 0 {
		return 1
	}
	if x > 1<<63 {
		panic("next power of 2 overflows uint64")
	}
	return 1 << (64 - bits.LeadingZeros64(x-1))
}
func time2freq(time Block) Block {
	// extend  into powers of 2
	l := nextPowerOfTwo(uint64(len(time)))
	for len(time) < int(l) {
		time = append(time, 0)
	}

	ω := getPrimitiveRoot(l)
	freq := CooleyTukeyFFT(time, ω)

	for len(freq) > 0 && freq[len(freq)-1] == 0 {
		freq = freq[:len(freq)-1]
	}
	return freq
}

func freq2time(freq Block) Block {
	// extend  into powers of 2
	l := nextPowerOfTwo(uint64(len(freq)))
	for len(freq) < int(l) {
		freq = append(freq, 0)
	}

	time := Block{}
	ω := getPrimitiveRoot(l)
	for _, f := range CooleyTukeyFFT(freq, inv(ω)) {
		time = append(time, mul(f, inv(l)))
	}

	for len(time) > 0 && time[len(time)-1] == 0 {
		time = time[:len(time)-1]
	}
	return time
}
