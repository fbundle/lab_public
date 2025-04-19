# ca
computational algebra

## pardic.go

implementation of $p$-adic integers

## uint1792.go

## uint_ntt.go and int_ntt.go

- my greatest appreciation to [apgoucher](https://cp4space.hatsya.com/2021/09/01/an-efficient-prime-for-number-theoretic-transforms/) to the prime $p = 2^{64} - 2^{32} + 1$ with $8$ being the $64$-th primitive root of unity in mod $p$
- (almost) arbitrary precision integer ($4294967294 \times 16$ bits $\sim$ 8GB per integer)  
- multiplication using FFT
- division using Newton iteration
- support negative integers
