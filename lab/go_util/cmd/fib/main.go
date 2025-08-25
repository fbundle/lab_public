package main

import (
	"fmt"
	"github.com/fbundle/lab_public/lab/go_util/pkg/fib"
	"github.com/fbundle/lab_public/lab/go_util/pkg/int_ntt"
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

	x := fib.Fib(int_ntt.Nat{}, uint64(n))
	_ = x
	fmt.Println(x)
}
