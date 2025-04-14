package mp

import (
	"fmt"
	"os"
)

const (
	P    = 18446744069414584321 // p = 2^64 - 2^32 + 1 and 2 is the 192-th primitive root of unity
	R, N = 8, 64                // 8 is 64-th primitive root of unity
	B    = 1 << 16              // choose base 2^d so that N * B * B < max_uint64
)

const (
	invR = 16140901060737761281 // inverse of R in mod P
	invN = 18158513693329981441 // inverse of N in mod P
)

type Uint1024Block = [N]uint64

type Uint1024 struct {
	Time Uint1024Block
	Freq Uint1024Block
}

func NewUint1024FromTime(time Uint1024Block) Uint1024 {
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] += q
		}
	}
	return Uint1024{
		Time: time,
		Freq: dft(time, N, R),
	}
}

func NewUint1024FromFreq(freq Uint1024Block) Uint1024 {
	time := Uint1024Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = (invN * f) % P
	}
	return Uint1024{
		Time: time,
		Freq: freq,
	}
}

func NewUint1024FromUint64(u uint64) Uint1024 {
	return NewUint1024FromTime(Uint1024Block{u})
}

func (a Uint1024) Add(b Uint1024) Uint1024 {
	time := Uint1024Block{}
	for i := 0; i < N; i++ {
		time[i] = a.Time[i] + b.Time[i]
	}
	return NewUint1024FromTime(time)
}

func (a Uint1024) Mul(b Uint1024) Uint1024 {
	freq := Uint1024Block{}
	for i := 0; i < N; i++ {
		freq[i] = a.Freq[i] * b.Freq[i]
	}
	return NewUint1024FromFreq(freq)
}

// TODO - implement dft and dftCT using uint128 because multiplying two uint64s cause overflow

func TestDft() {
	invR := invmod(R, P)
	fmt.Println((R * invR) % P)
	os.Exit(0)
	// invN := invmod(N, P)

	m1 := makeDftMat(N, R)
	m2 := makeDftMat(N, invR)
	w := make([][]uint64, N)
	for i := 0; i < N; i++ {
		w[i] = make([]uint64, N)
	}
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			for k := 0; k < N; k++ {
				w[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}
	fmt.Println(w)

}

func makeDftMat(n int, omega uint64) [][]uint64 {
	w := make([][]uint64, n)
	for i := 0; i < n; i++ {
		w[i] = make([]uint64, n)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			w[i][j] = powmod(omega, uint64(i*j), P)
		}
	}
	return w
}

func dft(block Uint1024Block, n int, omega uint64) Uint1024Block {
	w := makeDftMat(n, omega)
	out := Uint1024Block{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			out[i] += w[i][j] * block[j]
		}
	}
	return out
}

// dftCT :Cooley-Tukey algorithm
func dftCT(block Uint1024Block, n int, omega uint64) Uint1024Block {
	if n == 1 {
		return block
	}
	if n <= 0 || n > N || n%2 != 0 {
		panic("n must be 1 or even in [2, 64]")
	}
	var even, odd Uint1024Block
	for i := 0; i < n/2; i++ {
		even[i] = block[2*i]
		odd[i] = block[2*i+1]
	}
	evenFFT := dftCT(even, n/2, omega*omega)
	oddFFT := dftCT(odd, n/2, omega*omega)

	result := Uint1024Block{}
	var omegaPow uint64 = 1 // omega^0
	for i := 0; i < n/2; i++ {
		t := (omegaPow * oddFFT[i]) % P
		result[i] = (evenFFT[i] + t) % P
		result[i+n/2] = (evenFFT[i] - t) % P
		omegaPow = (omegaPow * omega) % P
	}
	return result
}
