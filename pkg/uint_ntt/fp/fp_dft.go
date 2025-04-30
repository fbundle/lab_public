package fp

import (
	"ca/pkg/uint_ntt/util"
	"ca/pkg/vec"
	"math/bits"
	"sync"
)

var dft = IterativeCooleyTukeyFFT

func Mul(aTime util.Block, bTime util.Block) util.Block {
	l := util.GetNextPowerOfTwo(uint64(aTime.Len() + bTime.Len()))
	aFreq, bFreq := time2freq(aTime, l), time2freq(bTime, l)
	freq := util.Block{}
	for i := 0; i < int(l); i++ {
		freq = freq.Set(i, mul(aFreq.Get(i), bFreq.Get(i)))
	}
	time := freq2time(freq, l)
	return time
}

func time2freq(time util.Block, length uint64) util.Block {
	// extend  into powers of 2
	n := util.GetNextPowerOfTwo(length)
	time = time.Slice(0, int(n)) // extend to length n

	omega := getPrimitiveRoot(n)
	freq := dft(time, omega)
	return util.TrimBlock(freq)
}

func freq2time(freq util.Block, length uint64) util.Block {
	// extend  into powers of 2
	n := util.GetNextPowerOfTwo(length)
	freq = freq.Slice(0, int(n)) // extend to length n
	omega := getPrimitiveRoot(n)
	il := inv(n)

	time := dft(freq, inv(omega))
	for i := 0; i < time.Len(); i++ {
		f := time.Get(i)
		time = time.Set(i, mul(f, il))
	}

	return util.TrimBlock(time)
}

// CooleyTukeyFFT :Cooley-Tukey algorithm
func CooleyTukeyFFT(x util.Block, omega uint64) util.Block {
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

	eFFT := CooleyTukeyFFT(e, omega_2)
	oFFT := CooleyTukeyFFT(o, omega_2)

	y := vec.MakeVec[uint64](n)
	for i := 0; i < n/2; i++ {
		j := i + n/2
		t := mul(pow(omega, uint64(i)), oFFT.Get(i))
		y = y.Set(i, add(eFFT.Get(i), t))
		y = y.Set(j, sub(eFFT.Get(i), t))
	}
	return y
}

// IterativeCooleyTukeyFFT : from deepseek
func IterativeCooleyTukeyFFT(x util.Block, omega uint64) util.Block {
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
func IterativeParallelCooleyTukeyFFT(x util.Block, omega uint64) util.Block {
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
