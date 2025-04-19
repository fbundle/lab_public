package uint_ntt

const (
	base  = 1 << 16 // pick base = 2^d, max_n * base * base < p so that multiplication won't overflow
	max_n = 4294967294
)

// Block : polynomial in F_p[X]
type Block struct {
	block []uint64
}

func makeBlock(n int) Block {
	return Block{make([]uint64, n)}
}

func (b Block) clone() Block {
	c := makeBlock(b.len())
	copy(c.block, b.block)
	return c
}

func (b Block) append(v uint64) Block {
	b.block = append(b.block, v)
	return b
}

func (b Block) len() int {
	return len(b.block)
}

func (b Block) get(i int) uint64 {
	if i >= b.len() {
		return 0
	}
	return b.block[i]
}

func (b Block) set(i int, v uint64) Block {
	for i >= b.len() {
		b.block = append(b.block, 0)
		if b.len() > max_n {
			panic("too many blocks")
		}
	}
	b.block[i] = v
	return b
}

func (b Block) slice(beg int, end int) Block {
	return Block{b.block[beg:end]}
}
