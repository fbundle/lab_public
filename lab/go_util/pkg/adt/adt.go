package adt

import "errors"

type Option[T any] struct {
	Val T
	Err error
}

func (o Option[T]) Unwrap(val *T) error {
	if val != nil {
		*val = o.Val
	}
	return o.Err
}

func Wrap[T any](f func(...any) (T, error)) func(...any) Option[T] {
	return func(args ...any) Option[T] {
		val, err := f(args...)
		if err != nil {
			return Error[T](err)
		}
		return Some(val)
	}
}

func Error[T any](err error) Option[T] {
	return Option[T]{
		Err: err,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Val: val,
		Err: nil,
	}
}

func Map[T1 any, T2 any](f func(T1) (T2, error)) func(Option[T1]) Option[T2] {
	return func(o Option[T1]) Option[T2] {
		if o.Err != nil {
			return Error[T2](o.Err)
		}
		val2, err := f(o.Val)
		if err != nil {
			return Error[T2](err)
		}
		return Some(val2)
	}
}

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

func (s Sum2[T1, T2]) Unwrap1() Option[T1] {
	if v, ok := s.val.(T1); ok {
		return Some(v)
	} else {
		return Error[T1](ErrType)
	}
}
func (s Sum2[T1, T2]) Unwrap2() Option[T2] {
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

type Option2[T1 any, T2 any] struct {
	val1 T1
	val2 T2
	err  error
}

func (o Option2[T1, T2]) Unwrap(val1 *T1, val2 *T2) error {
	if val1 != nil {
		*val1 = o.val1
	}
	if val2 != nil {
		*val2 = o.val2
	}
	return o.err
}

func Some2[T1 any, T2 any](val1 T1, val2 T2) Option2[T1, T2] {
	return Option2[T1, T2]{
		val1: val1,
		val2: val2,
		err:  nil,
	}
}

func Error2[T1 any, T2 any](err error) Option2[T1, T2] {
	return Option2[T1, T2]{
		err: err,
	}
}
