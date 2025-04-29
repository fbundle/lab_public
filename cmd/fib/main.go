package main

import (
	"ca/pkg/fib"
	"ca/pkg/uint_ntt"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: go run main.go <integer>")
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Error: argument must be an integer")
		return
	}

	x := fib.Fib(uint_ntt.UintNTT{}, uint64(n))
	_ = x
	fmt.Println(x)
}
