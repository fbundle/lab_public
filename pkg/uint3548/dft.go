package uint3548

import "fmt"

func TestFFT() {
	x := New(Uint3584Block{131414234243, 12314555, 123131, 5345777, 7646456})
	z := FromFreq(x.Freq)
	fmt.Println(x)
	fmt.Println(z)
}

type DFT func(block Uint3584Block, n int, omega uint64) Uint3584Block

var dft DFT = CooleyTukeyFFT

func SetDefaultDFT(f DFT) {
	dft = f
}

func NaiveDFT(block Uint3584Block, n int, omega uint64) Uint3584Block {
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
	w := makeDftMat(n, omega)
	out := Uint3584Block{}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			out[i] = add(out[i], mul(w[i][j], block[j]))
		}
	}
	return out
}

// CooleyTukeyFFT :Cooley-Tukey algorithm
func CooleyTukeyFFT(block Uint3584Block, n int, omega uint64) Uint3584Block {
	if n == 1 {
		return block
	}
	if n <= 0 || n%2 != 0 {
		panic("n must be power of 2")
	}
	var even, odd Uint3584Block
	for i := 0; i < n/2; i++ {
		even[i] = block[2*i]
		odd[i] = block[2*i+1]
	}
	nextOmega := mul(omega, omega)
	evenFFT := CooleyTukeyFFT(even, n/2, nextOmega)
	oddFFT := CooleyTukeyFFT(odd, n/2, nextOmega)

	result := Uint3584Block{}
	var omegaPow uint64 = 1 // omega^0
	for i := 0; i < n/2; i++ {
		t := mul(omegaPow, oddFFT[i])
		result[i] = add(evenFFT[i], t)
		result[i+n/2] = sub(evenFFT[i], t)
		omegaPow = mul(omegaPow, omega)
	}
	return result
}
