package uint1792

import (
	"fmt"
	"os"
)

// NaiveDFT : naive implementation of DFT - for reference
// construct DFT matrix w of size (n, n) with omega as the root of unity
// return y = wx
func NaiveDFT(x []uint64, ω uint64) (y []uint64) {
	n := len(x)
	_, _ = fmt.Fprintf(os.Stderr, "WARNING : this implementation is for reference, use FFT instead")
	makeDftMat := func(n int, ω uint64) [][]uint64 {
		w := make([][]uint64, n)
		for i := 0; i < n; i++ {
			w[i] = make([]uint64, n)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				w[i][j] = pow(ω, uint64(i*j))
			}
		}
		return w
	}
	y = make([]uint64, n)
	w := makeDftMat(n, ω)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			y[i] = add(y[i], mul(w[i][j], x[j]))
		}
	}
	return y
}

// CooleyTukeyFFT :Cooley-Tukey algorithm
func CooleyTukeyFFT(x []uint64, ω uint64) (y []uint64) {
	n := len(x)
	if n == 1 {
		return x
	}
	if n <= 0 || n%2 != 0 {
		panic("n must be power of 2")
	}
	var e, o []uint64 // even and odd values of x
	for i := 0; i < n/2; i++ {
		e = append(e, x[2*i])
		o = append(o, x[2*i+1])
	}
	ω2 := mul(ω, ω)
	eFFT := CooleyTukeyFFT(e, ω2)
	oFFT := CooleyTukeyFFT(o, ω2)

	y = make([]uint64, n)
	var ωn uint64 = 1 // ω^0
	for i := 0; i < n/2; i++ {
		t := mul(ωn, oFFT[i])
		y[i] = add(eFFT[i], t)
		y[i+n/2] = sub(eFFT[i], t)
		ωn = mul(ωn, ω)
	}
	return y
}
