package uint_ntt

import (
	v "ca/pkg/vec"
)

type Vec[T any] = v.Vec[T]

func makeVec[T any](n int) Vec[T] {
	return Vec[T](v.MakeVec[T](n))
}
