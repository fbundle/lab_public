package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/fbundle/lab_public/lab/go_util/pkg/adt"
)

func main() {
	x := adt.Pipe[string, int]{
		Input: adt.Ok[string]("45+12"),
	}.Push(adt.Bind(func(s string) adt.Except[adt.Prod2[string, string]] {
		parts := strings.SplitN(s, "+", 2)
		if len(parts) != 2 {
			return adt.Err[adt.Prod2[string, string]](errors.New("invalid input"))
		}
		return adt.Ok(adt.NewProd2[string, string](parts[0], parts[1]))
	})).Push(adt.Bind(func(p adt.Prod2[string, string]) adt.Except[adt.Prod2[int, int]] {
		v1, v2 := p.Unwrap()
		i1, err := strconv.Atoi(v1)
		if err != nil {
			return adt.Err[adt.Prod2[int, int]](err)
		}
		i2, err := strconv.Atoi(v2)
		if err != nil {
			return adt.Err[adt.Prod2[int, int]](err)
		}
		return adt.Ok(adt.NewProd2[int, int](i1, i2))
	})).Push(adt.Map(func(p adt.Prod2[int, int]) adt.Except[int] {
		i1, i2 := p.Unwrap()
		return adt.Ok(i1 + i2)
	})).Finalize()
	var v int
	fmt.Println(x.Unwrap(&v))
}
