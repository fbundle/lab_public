package uint_ntt

import "ca/pkg/vec"

type Vec[T any] = vec.Vec[T]

func MakeVec[T any](n int) Vec[T] {
	return Vec[T](vec.MakeVec[T](n))
}
