package ntt

import (
	"go_util/pkg/vec"
	"math/bits"
	"sync"
)

var dft = iterativeCooleyTukeyFFT

func time2freq(time Block, length uint64) Block {
	// extend  into powers of 2
	n := nextPowerOfTwo(length)
	time = time.Slice(0, int(n)) // extend to length n

	omega := getPrimitiveRoot(n)
	freq := dft(time, omega)
	return freq
}

func freq2time(freq Block, length uint64) Block {
	// extend  into powers of 2
	n := nextPowerOfTwo(length)
	freq = freq.Slice(0, int(n)) // extend to length n
	omega := getPrimitiveRoot(n)
	il := inv(n)

	time := dft(freq, inv(omega))
	for i := 0; i < time.Len(); i++ {
		f := time.Get(i)
		time = time.Set(i, mul(f, il))
	}

	return time
}

// cooleyTukeyFFT :Cooley-Tukey algorithm
func cooleyTukeyFFT(x Block, omega uint64) Block {
	n := x.Len()
	if n == 1 {
		return x
	}
	if n <= 0 || n%2 != 0 {
		panic("n must be power of 2")
	}
	// even and odd values of x
	e, o := vec.MakeVec[uint64](n/2), vec.MakeVec[uint64](n/2)
	for i := 0; i < n/2; i++ {
		e = e.Set(i, x.Get(2*i))
		o = o.Set(i, x.Get(2*i+1))
	}
	omega_2 := mul(omega, omega)

	eFFT := cooleyTukeyFFT(e, omega_2)
	oFFT := cooleyTukeyFFT(o, omega_2)

	y := vec.MakeVec[uint64](n)
	for i := 0; i < n/2; i++ {
		j := i + n/2
		t := mul(pow(omega, uint64(i)), oFFT.Get(i))
		y = y.Set(i, add(eFFT.Get(i), t))
		y = y.Set(j, sub(eFFT.Get(i), t))
	}
	return y
}

// iterativeCooleyTukeyFFT : from deepseek
func iterativeCooleyTukeyFFT(x Block, omega uint64) Block {
	n := x.Len()
	if n&(n-1) != 0 {
		panic("n must be power of two")
	}
	logN := bits.TrailingZeros64(uint64(n))

	// Bit-reversal permutation (creates new vector)
	// Reverse bits helper (unchanged)
	reverseBits := func(num uint32, bits int) uint32 {
		reversed := uint32(0)
		for i := 0; i < bits; i++ {
			reversed = (reversed << 1) | (num & 1)
			num >>= 1
		}
		return reversed
	}
	reversed := vec.MakeVec[uint64](n)
	for i := 0; i < n; i++ {
		rev := reverseBits(uint32(i), logN)
		reversed = reversed.Set(i, x.Get(int(rev)))
	}

	// Main computation (build new vector at each stage)
	for stage := 1; stage <= logN; stage++ {
		m := 1 << stage
		wm := pow(omega, uint64(n>>stage))
		newVec := reversed.Clone()

		for k := 0; k < n; k += m {
			w := uint64(1)
			for j := 0; j < m/2; j++ {
				idx1 := k + j
				idx2 := k + j + m/2

				u := reversed.Get(idx1)
				t := mul(reversed.Get(idx2), w)

				newVec = newVec.Set(idx1, add(u, t))
				newVec = newVec.Set(idx2, sub(u, t))

				w = mul(w, wm)
			}
		}
		reversed = newVec
	}

	return reversed
}
func iterativeParallelCooleyTukeyFFT(x Block, omega uint64) Block {
	n := x.Len()
	if n&(n-1) != 0 {
		panic("n must be power of two")
	}
	logN := bits.TrailingZeros64(uint64(n))

	reverseBits := func(num uint32, bits int) uint32 {
		reversed := uint32(0)
		for i := 0; i < bits; i++ {
			reversed = (reversed << 1) | (num & 1)
			num >>= 1
		}
		return reversed
	}
	reversed := vec.MakeVec[uint64](n)
	for i := 0; i < n; i++ {
		rev := reverseBits(uint32(i), logN)
		reversed = reversed.Set(i, x.Get(int(rev)))
	}

	for stage := 1; stage <= logN; stage++ {
		m := 1 << stage
		wm := pow(omega, uint64(n>>stage))
		newVec := reversed.Clone()

		var wg sync.WaitGroup

		for k := 0; k < n; k += m {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()

				w := uint64(1)
				for j := 0; j < m/2; j++ {
					idx1 := k + j
					idx2 := k + j + m/2

					u := reversed.Get(idx1)
					t := mul(reversed.Get(idx2), w)

					newVec = newVec.Set(idx1, add(u, t))
					newVec = newVec.Set(idx2, sub(u, t))

					w = mul(w, wm)
				}
			}(k)
		}

		wg.Wait()
		reversed = newVec
	}

	return reversed
}
