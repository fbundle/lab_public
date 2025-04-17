package uint1792

import (
	"fmt"
	"os"
)

type DFT func(x Block, n int, omega uint64) (y Block)

var dft DFT = CooleyTukeyFFT

func SetDefaultDFT(f DFT) {
	dft = f
}

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
func CooleyTukeyFFT(x Block, n int, omega uint64) (y Block) {
	if n == 1 {
		return x
	}
	if n <= 0 || n%2 != 0 {
		panic("n must be power of 2")
	}
	var e, o Block // even and odd values of x
	for i := 0; i < n/2; i++ {
		e[i] = x[2*i]
		o[i] = x[2*i+1]
	}
	nextOmega := mul(omega, omega)
	eFFT := CooleyTukeyFFT(e, n/2, nextOmega)
	oFFT := CooleyTukeyFFT(o, n/2, nextOmega)

	var omegaPow uint64 = 1 // omega^0
	for i := 0; i < n/2; i++ {
		t := mul(omegaPow, oFFT[i])
		y[i] = add(eFFT[i], t)
		y[i+n/2] = sub(eFFT[i], t)
		omegaPow = mul(omegaPow, omega)
	}
	return y
}

func MixedRadixFFT192(x Block, omega uint64) (y Block) {
	// Factor n = r * s
	r, s := 3, 64

	// Reshape x into matrix x[r][s]
	matrix := make([]Block, r)
	for i := range matrix {
		for j := 0; j < s; j++ {
			matrix[i][j] = x[i*s+j]
		}
	}

	// Row-wise FFTs (length s)
	for i := 0; i < r; i++ {
		matrix[i] = CooleyTukeyFFT(matrix[i], s, pow(omega, uint64(r))) // ω^r
	}

	// Twiddle factors
	for i := 0; i < r; i++ {
		for j := 0; j < s; j++ {
			w := pow(omega, uint64(i*j)) // ω^{i*j}
			matrix[i][j] = mul(matrix[i][j], w)
		}
	}

	// Column-wise FFTs (length r)
	for j := 0; j < s; j++ {
		col := make([]uint64, r)
		for i := 0; i < r; i++ {
			col[i] = matrix[i][j]
		}
		col = NaiveDFT(col, pow(omega, uint64(s))) // ω^s
		for i := 0; i < r; i++ {
			matrix[i][j] = col[i]
		}
	}

	// Flatten output
	for i := 0; i < r; i++ {
		for j := 0; j < s; j++ {
			y[i*s+j] = matrix[i][j]
		}
	}
	return y
}
