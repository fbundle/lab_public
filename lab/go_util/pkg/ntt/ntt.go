package ntt

import (
	"github.com/fbundle/lab_public/lab/go_util/pkg/vec"
)

type Block = vec.Vec[uint64]

func Mul(aTime Block, bTime Block) Block {
	l := nextPowerOfTwo(uint64(aTime.Len() + bTime.Len()))
	aFreq, bFreq := time2freq(aTime, l), time2freq(bTime, l)
	freq := Block{}
	for i := 0; i < int(l); i++ {
		freq = freq.Set(i, mul(aFreq.Get(i), bFreq.Get(i)))
	}
	time := freq2time(freq, l)
	return time
}
