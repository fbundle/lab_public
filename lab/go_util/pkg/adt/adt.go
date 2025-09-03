package adt

import "errors"

type Prod2[T1 any, T2 any] struct {
	t1 T1
	t2 T2
}

func (p Prod2[T1, T2]) Unwrap() (t1 T1, t2 T2) {
	return p.t1, p.t2
}

func NewProd2[T1 any, T2 any](t1 T1, t2 T2) Prod2[T1, T2] {
	return Prod2[T1, T2]{
		t1: t1,
		t2: t2,
	}
}

var ErrType = errors.New("type_error")

type Sum2[T1 any, T2 any] struct {
	val any
}

func (s Sum2[T1, T2]) Unwrap1() Result[T1] {
	if v, ok := s.val.(T1); ok {
		return Some(v)
	} else {
		return Error[T1](ErrType)
	}
}
func (s Sum2[T1, T2]) Unwrap2() Result[T2] {
	if v, ok := s.val.(T2); ok {
		return Some(v)
	} else {
		return Error[T2](ErrType)
	}
}

func NewSum2[T1 any, T2 any](val any) Sum2[T1, T2] {
	okCount := 0
	if _, ok := val.(T1); ok {
		okCount++
	}
	if _, ok := val.(T2); ok {
		okCount++
	}
	if okCount == 0 {
		panic(ErrType)
	}
	return Sum2[T1, T2]{val: val}
}
