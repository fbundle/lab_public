package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var a atomic.Value
	swapped := a.CompareAndSwap(nil, 1)
	fmt.Println(swapped, a.Load())
}
