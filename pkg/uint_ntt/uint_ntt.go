package uint_ntt

import (
	"strings"
)

const (
	base = 1 << 16 // pick base = 2^d, max_n * base * base < p so that multiplication won't overflow
)

// UintNTT : represents nonnegative integers by a_0 + a_1 base + a_2 base^2 + ... + a_{N-1} base^{N-1}
type UintNTT struct {
	time block[uint64]
}

func (a UintNTT) Zero() UintNTT {
	return UintNTT{time: block[uint64]{}}
}

func (a UintNTT) One() UintNTT {
	return FromUint64(1)
}

func FromUint64(x uint64) UintNTT {
	return fromTime(makeBlock[uint64](1).set(0, x))
}

func (a UintNTT) Uint64() uint64 {
	return a.time.get(0) + a.time.get(1)*base + a.time.get(2)*base*base + a.time.get(3)*base*base*base
}

func fromTime(time block[uint64]) UintNTT {
	// reduce to base
	var q, r uint64 = 0, 0
	for i := 0; i < time.len(); i++ {
		q, r = time.get(i)/base, time.get(i)%base
		time = time.set(i, r)
		if i+1 < time.len() {
			time = time.set(i+1, add(time.get(i+1), q))
		}
	}
	time = time.set(time.len(), q)
	time = trimZeros(time)
	return UintNTT{
		time: time,
	}
}

// trimZeros : trim zeros at high degree
func trimZeros(block block[uint64]) block[uint64] {
	for block.len() > 0 && block.get(block.len()-1) == 0 {
		block = block.slice(0, block.len()-1)
	}
	return block
}

func FromString(s string) UintNTT {
	s = strings.ToLower(s)
	if s[0:2] != "0x" {
		panic("string must start with 0x")
	}
	s = s[2:]

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

	time := block[uint64]{}
	for i := 0; i < len(base16); i += 4 {
		var x uint64 = 0
		x += uint64(base16[i])
		x += uint64(base16[i+1]) * 16
		x += uint64(base16[i+2]) * 16 * 16
		x += uint64(base16[i+3]) * 16 * 16 * 16
		time = time.set(i/4, x)
	}

	return fromTime(time)
}

func (a UintNTT) String() string {
	if base != 1<<16 {
		panic("not implemented")
	}
	// convert base 2^16 to base16 (2^4)
	var base16 []byte = nil
	for i := 0; i < a.time.len(); i++ {
		x := a.time.get(i)
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
	l := max(a.time.len(), b.time.len())
	cTime := makeBlock[uint64](l)
	for i := 0; i < l; i++ {
		cTime = cTime.set(i, add(a.time.get(i), b.time.get(i)))
	}
	return fromTime(cTime)
}

func (a UintNTT) Mul(b UintNTT) UintNTT {
	l := nextPowerOfTwo(uint64(a.time.len() + b.time.len()))

	aFreq, bFreq := time2freq(a.time, l), time2freq(b.time, l)
	freq := block[uint64]{}
	for i := 0; i < int(l); i++ {
		freq = freq.set(i, mul(aFreq.get(i), bFreq.get(i)))
	}
	time := freq2time(freq, l)
	return fromTime(time)
}

// Sub - subtract b from a using long subtraction
// if a < b, return 2nd complement and false
func (a UintNTT) Sub(b UintNTT) (UintNTT, bool) {
	l := max(a.time.len(), b.time.len())
	cTime := a.time.clone()
	var borrow uint64 = 0 // either zero or one
	for i := 0; i < l; i++ {
		x := sub(cTime.get(i)+base, b.time.get(i)+borrow) // x in [0, 2^{32}-1]
		cTime = cTime.set(i, x%base)
		borrow = 1 - x/base
	}
	return fromTime(cTime), borrow == 0
}

func (a UintNTT) IsZero() bool {
	for i := 0; i < a.time.len(); i++ {
		if a.time.get(i) != 0 {
			return false
		}
	}
	return true
}

func (a UintNTT) Cmp(b UintNTT) int {
	l := max(a.time.len(), b.time.len())
	for i := l - 1; i >= 0; i-- {
		if a.time.get(i) > b.time.get(i) {
			return +1
		}
		if a.time.get(i) < b.time.get(i) {
			return -1
		}
	}
	return 0
}

func (a UintNTT) shiftRight(n int) UintNTT {
	if n > a.time.len() {
		return UintNTT{}
	}
	cTime := a.time.clone()
	cTime = cTime.slice(n, cTime.len())
	return fromTime(cTime)
}

// inv : let m = 2^{16n}
// approx root of f(x) = m / x - a using Newton method
// error at most 1
func (a UintNTT) pinv(n int) UintNTT {
	if a.IsZero() {
		panic("division by zero")
	}
	x := FromUint64(1)
	// Newton iteration
	for {
		// x_{n+1} = x_n + x_n - (a x_n^2) / m
		left := x.Add(x)
		right := a.Mul(x).Mul(x).shiftRight(n)
		x1, ok := left.Sub(right)
		if !ok {
			// x is always on the left of the root - this will not happen
			panic("subtract overflow")
		}
		if x1.Cmp(x) == 0 {
			break
		}
		x = x1
	}
	return x
}

func (a UintNTT) Div(b UintNTT) UintNTT {
	n := max(a.time.len(), b.time.len()) + 1 // large enough
	x := b.pinv(n)
	return a.Mul(x).shiftRight(n)
}

func (a UintNTT) Mod(b UintNTT) UintNTT {
	x := a.Div(b)
	m, ok := a.Sub(b.Mul(x))
	if !ok {
		// this will not happen
		panic("subtraction overflow")
	}
	return m
}
