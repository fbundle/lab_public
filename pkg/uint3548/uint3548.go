package uint3548

// Uint3548Block:  a block of N uint64s, each is in mod P
type Uint3584Block = [N]uint64

type Uint3584 struct {
	Time Uint3584Block
	Freq Uint3584Block
}

func (u Uint3584) Uint64() uint64 {
	return u.Time[0] + u.Time[1]*B + u.Time[2]*B*B
}

func New(time Uint3584Block) Uint3584 {
	// trim to [0, B-1] for easier conversion
	for i := 0; i < N; i++ {
		q, r := time[i]/B, time[i]%B
		time[i] = r
		if i+1 < N {
			time[i+1] = time[i+1] + q
		}
	}
	return Uint3584{
		Time: time,
		Freq: time2freq(time),
	}
}

func FromFreq(freq Uint3584Block) Uint3584 {
	return Uint3584{
		Time: freq2time(freq),
		Freq: freq,
	}
}

func FromUint64(u uint64) Uint3584 {
	return New(Uint3584Block{u})
}

func (a Uint3584) Add(b Uint3584) Uint3584 {
	time := Uint3584Block{}
	for i := 0; i < N; i++ {
		time[i] = add(a.Time[i], b.Time[i])
	}
	return New(time)
}

func (a Uint3584) Mul(b Uint3584) Uint3584 {
	freq := Uint3584Block{}
	for i := 0; i < N; i++ {
		freq[i] = mul(a.Freq[i], b.Freq[i])
	}
	return FromFreq(freq)
}

func time2freq(time Uint3584Block) Uint3584Block {
	return dft(time, N, R)
}

func freq2time(freq Uint3584Block) Uint3584Block {
	time := Uint3584Block{}
	for i, f := range dft(freq, N, invR) {
		time[i] = mul(f, invN)
	}
	return time
}
