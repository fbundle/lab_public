package uint1792

import (
	"fmt"
	"os"
)

// NaiveDFT : naive implementation of DFT - for reference
// construct DFT matrix w of size (n, n) with omega as the root of unity
// return y = wx
func NaiveDFT(x []uint64, omega uint64) (y []uint64) {
	n := len(x)
	_, _ = fmt.Fprintf(os.Stderr, "WARNING : this implementation is for reference, use FFT instead")
	makeDftMat := func(n int, omega uint64) [][]uint64 {
		w := make([][]uint64, n)
		for i := 0; i < n; i++ {
			w[i] = make([]uint64, n)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				w[i][j] = pow(omega, uint64(i*j))
			}
		}
		return w
	}
	y = make([]uint64, n)
	w := makeDftMat(n, omega)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			y[i] = add(y[i], mul(w[i][j], x[j]))
		}
	}
	return y
}

// CooleyTukeyFFT :Cooley-Tukey algorithm
func CooleyTukeyFFT(x []uint64, omega uint64) (y []uint64) {
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
	nextOmega := mul(omega, omega)
	eFFT := CooleyTukeyFFT(e, nextOmega)
	oFFT := CooleyTukeyFFT(o, nextOmega)

	y = make([]uint64, n)
	var omegaPow uint64 = 1 // omega^0
	for i := 0; i < n/2; i++ {
		t := mul(omegaPow, oFFT[i])
		y[i] = add(eFFT[i], t)
		y[i+n/2] = sub(eFFT[i], t)
		omegaPow = mul(omegaPow, omega)
	}
	return y
}
