package adt

import (
	"errors"

	"github.com/fbundle/lab_public/lab/go_util/pkg/compose"
)

type Except[T any] struct {
	val T
	err error
}

func (e Except[T]) Unwrap(val *T) error {
	if val != nil {
		*val = e.val
	}
	return e.err
}

func Err[T any](err error) Except[T] {
	return Except[T]{
		err: err,
	}
}

func Ok[T any](val T) Except[T] {
	return Except[T]{
		val: val,
		err: nil,
	}
}

func Eval[T1 any, T2 any](v1 T1, f func(T1) T2) T2 {
	return f(v1)
}

func Bind[T1 any, T2 any](f func(T1) Except[T2]) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return f(e1.val)
	}
}

func Map[T1 any, T2 any](f func(T1) T2) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return Ok[T2](f(e1.val))
	}
}

func Seq[T1 any, T2 any](f Except[func(T1) T2]) func(e1 Except[T1]) Except[T2] {
	return func(e1 Except[T1]) Except[T2] {
		if f.err != nil {
			return Err[T2](f.err)
		}
		if e1.err != nil {
			return Err[T2](e1.err)
		}
		return Ok[T2](f.val(e1.val))
	}
}

type Pipe[T1 any, T2 any] struct {
	Input Except[T1]
	pipe  []any
}

func (p Pipe[T1, T2]) Push(f any) Pipe[T1, T2] {
	p.pipe = append(p.pipe, f)
	return p
}

func (p Pipe[T1, T2]) Finalize() Except[T2] {
	f := compose.Compose(p.pipe...).(func(Except[T1]) Except[T2])
	if f == nil {
		return Err[T2](errors.New("type_error"))
	}
	return f(p.Input)
}
