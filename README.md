# ca
computational algebra

## pardic.go

implementation of $p$-adic integers

## uint_ntt.go and int_ntt.go

- my greatest appreciation to [apgoucher](https://cp4space.hatsya.com/2021/09/01/an-efficient-prime-for-number-theoretic-transforms/) to the prime $p = 2^{64} - 2^{32} + 1$ with $g=7$ being the generator of $(\mathbb{Z}/p)^\times$
- (almost) arbitrary precision integer ($[0, 2^{4294967294 \times 16})$ $\sim$ 8GB)  
- multiplication using FFT
- division using Newton iteration
- support negative integers
