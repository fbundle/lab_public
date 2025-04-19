package uint_ntt

const (
	base  = 1 << 16 // pick base = 2^d, max_n * base * base < p so that multiplication won't overflow
	max_n = 4294967294
)

// block : polynomial in F_p[X]
type block struct {
	data []uint64 // TODO - optimize this
}

func makeBlock(n int) block {
	return block{make([]uint64, n)}
}

func (b block) clone() block {
	c := makeBlock(b.len())
	copy(c.data, b.data)
	return c
}

func (b block) len() int {
	return len(b.data)
}

func (b block) get(i int) uint64 {
	if i >= b.len() {
		return 0
	}
	return b.data[i]
}

func (b block) set(i int, v uint64) block {
	for i >= b.len() {
		b.data = append(b.data, 0)
		if b.len() > max_n {
			panic("too many blocks")
		}
	}
	b.data[i] = v
	return b
}

func (b block) slice(beg int, end int) block {
	for end > b.len()-1 {
		b.data = append(b.data, 0)
	}
	return block{b.data[beg:end]}
}
