package uint_ntt

import (
	"math/bits"
	"strings"
)

const (
	base = 1 << 16 // pick base = 2^d, N * base * base < p so that multiplication won't overflow
)

// Block : polynomial in F_p[X]
type Block []uint64

func (b Block) get(i int) uint64 {
	if i >= len(b) {
		return 0
	}
	return b[i]
}

// UintNTT : represents nonnegative integers by a_0 + a_1 base + a_2 base^2 + ... + a_{N-1} base^{N-1}
type UintNTT struct {
	Time Block
}

func (a UintNTT) Zero() UintNTT {
	return UintNTT{Time: Block{}}
}

func (a UintNTT) One() UintNTT {
	return FromUint64(1)
}

func FromUint64(x uint64) UintNTT {
	return fromTime(Block{x})
}

func (a UintNTT) Uint64() uint64 {
	return a.Time[0] + a.Time[1]*base + a.Time[2]*base*base
}

func fromTime(time Block) UintNTT {
	// trim time
	var q, r uint64 = 0, 0
	for i := 0; i < len(time); i++ {
		q, r = time[i]/base, time[i]%base
		time[i] = r
		if i+1 < len(time) {
			time[i+1] = add(time[i+1], q)
		}
	}
	time = append(time, q)
	for len(time) > 0 && time[len(time)-1] == 0 {
		time = time[:len(time)-1]
	}
	return UintNTT{
		Time: time,
	}
}

func FromString(s string) UintNTT {
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
	if base != 1<<16 {
		panic("not implemented")
	}
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

func (a UintNTT) String() string {
	if base != 1<<16 {
		panic("not implemented")
	}
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

func (a UintNTT) Add(b UintNTT) UintNTT {
	l := max(len(a.Time), len(b.Time))
	cTime := make(Block, l)
	for i := 0; i < l; i++ {
		cTime[i] = add(a.Time.get(i), b.Time.get(i))
	}
	return fromTime(cTime)
}

func (a UintNTT) Mul(b UintNTT) UintNTT {
	l := nextPowerOfTwo(uint64(len(a.Time) + len(b.Time)))

	aFreq, bFreq := time2freq(a.Time, l), time2freq(b.Time, l)
	freq := Block{}
	for i := 0; i < int(l); i++ {
		freq = append(freq, mul(aFreq.get(i), bFreq.get(i)))
	}
	time := freq2time(freq, l)
	return fromTime(time)
}

func (a UintNTT) Sub(b UintNTT) (UintNTT, bool) {
	l := max(len(a.Time), len(b.Time))
	cTime := make(Block, l)
	copy(cTime, a.Time)
	cTime = append(cTime, 1) // borrow 1
	for i := 0; i < l; i++ {
		// x in [0, 2^{32}-1]
		x := sub(cTime.get(i)+cTime.get(i+1)*base, b.Time.get(i))
		cTime[i] = x % base
		cTime[i+1] = x / base
	}
	overflow := cTime[len(cTime)-1] < 1
	cTime = cTime[:len(cTime)-1]
	return fromTime(cTime), overflow
}

func (a UintNTT) IsZero() bool {
	for i := 0; i < len(a.Time); i++ {
		if a.Time[i] != 0 {
			return false
		}
	}
	return true
}

func (a UintNTT) Cmp(b UintNTT) int {
	l := max(len(a.Time), len(b.Time))
	for i := l - 1; i >= 0; i-- {
		if a.Time.get(i) > b.Time.get(i) {
			return +1
		}
		if a.Time.get(i) < b.Time.get(i) {
			return -1
		}
	}
	return 0
}

func (a UintNTT) shiftRight(n int) UintNTT {
	if n > len(a.Time) {
		return UintNTT{}
	}
	l := len(a.Time) - n
	cTime := make(Block, l)
	copy(cTime, a.Time[n:])
	return fromTime(cTime)
}

// inv : let m = 2^n, approx root of f(x) = - a + m / x
func (a UintNTT) inv(n int) UintNTT {
	if a.IsZero() {
		panic("division by zero")
	}
	two := FromUint64(2)
	x := FromUint64(1)
	// Newton iteration
	// x_{n+1} = 2 x_n - (a x_n^2) / m
	for {
		// a / m = a.shiftRight(N/2)
		x1, overflow := two.Mul(x).Sub(a.Mul(x).Mul(x).shiftRight(n))
		if overflow {
			panic("subtraction overflow")
		}
		if x1.Cmp(x) == 0 {
			break
		}
		x = x1
	}
	return x
}

func (a UintNTT) Div(b UintNTT) UintNTT {
	n := 2 * len(b.Time)
	x := b.inv(n)
	return a.Mul(x).shiftRight(n)
}

func (a UintNTT) Mod(b UintNTT) UintNTT {
	x := a.Div(b)
	m, overflow := a.Sub(b.Mul(x))
	if overflow {
		panic("subtraction overflow")
	}
	return m
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
func time2freq(time Block, length uint64) Block {
	// extend  into powers of 2
	l := nextPowerOfTwo(length)
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

func freq2time(freq Block, length uint64) Block {
	// extend  into powers of 2
	l := nextPowerOfTwo(length)
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
