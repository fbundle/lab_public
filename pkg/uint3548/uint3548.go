package uint3548

const (
	P    = 18446744069414584321 // p = 2^64 - 2^32 + 1 and 2 is the 192-th primitive root of unity
	R, N = 8, 64                // 8 is 64-th primitive root of unity
	B    = 1 << 28              // choose base 2^d so that N * B * B < P
)

const (
	invR = 16140901060737761281 // inverse of R in mod P
	invN = 18158513693329981441 // inverse of N in mod P
)

type Uint3584Block = [N]uint64

type Uint3584 struct {
	Time Uint3584Block
	Freq Uint3584Block
}

func (u Uint3584) Uint64() uint64 {
	return u.Time[0] + u.Time[1]*B + u.Time[2]*B*B
}

func NewUint3548FromTime(time Uint3584Block) Uint3584 {
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] = time[i+1] + q
		}
	}
	return Uint3584{
		Time: time,
		Freq: dft(time, N, R),
	}
}

func NewUint3548FromFreq(freq Uint3584Block) Uint3584 {
	time := Uint3584Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = mul(f, invN)
	}
	return Uint3584{
		Time: time,
		Freq: freq,
	}
}

func NewUint3548FromUint64(u uint64) Uint3584 {
	return NewUint3548FromTime(Uint3584Block{u})
}

func (a Uint3584) Add(b Uint3584) Uint3584 {
	time := Uint3584Block{}
	for i := 0; i < N; i++ {
		time[i] = add(a.Time[i], b.Time[i])
	}
	return NewUint3548FromTime(time)
}

func (a Uint3584) Mul(b Uint3584) Uint3584 {
	freq := Uint3584Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(a.Freq[i], b.Freq[i])
	}
	return NewUint3548FromFreq(freq)
}
