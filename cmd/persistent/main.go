package main

import (
	"fmt"
	"github.com/fbundle/go_util/pkg/persistent"
)

type Int int

func (i Int) Less(j Int) bool {
	return i < j
}

func testPersistent() {
	l := persistent.EmptyList[int]()
	l = l.Push(1).Push(2).Push(3)
	for i, elem := range l.Iter {
		fmt.Println(i, elem)
	}

	m := persistent.EmptyMap[Int, string]()
	m = m.Set(1, "one").Set(2, "two").Set(3, "three")
	for k, v := range m.Iter {
		fmt.Println(k, v)
	}
}

func main() {
	testPersistent()
}
