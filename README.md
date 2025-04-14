# ca
computational algebra

## pardic.go

implementation of $p$-adic integers

## uint1792.go

- `uint1792` with division $\lfloor a / b \rfloor$ for $1 < b < 2^{896}$ (can increase the bound close to $2^{1792}$)
- multiplication using FFT
- TODO : use mixed-radix Cooley-Tukey FFT
- TODO : use different primes then use CRT to construct output in larger base
