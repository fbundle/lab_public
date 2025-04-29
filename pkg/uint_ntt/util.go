package uint_ntt

// canonicalize : rewrite so that all coefficients in [0, base)
func canonicalize(time Block) Block {
	originalLen := time.Len()
	for i := 0; i < originalLen; i++ {
		q, r := time.Get(i)/base, time.Get(i)%base
		time = time.Set(i, r)
		time = time.Set(i+1, time.Get(i+1)+q)
	}
	if time.Len() > 0 {
		for time.Get(time.Len()-1) >= base {
			q, r := time.Get(time.Len()-1)/base, time.Get(time.Len()-1)%base
			time = time.Set(time.Len()-1, r)
			time = time.Set(time.Len(), q)
		}
	}
	return time
}

// trim : trim unused zeros at high degree
func trim(block Block) Block {
	for block.Len() > 0 && block.Get(block.Len()-1) == 0 {
		block = block.Slice(0, block.Len()-1)
	}
	return block
}
