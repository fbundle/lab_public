package sat

import (
	"context"
	"time"
)

type ContextKey int

const (
	ContextKeySatisfiable ContextKey = 0
	ContextKeyAssignment  ContextKey = 1
)

type Value = int

const (
	ValueTrue    Value = 1
	ValueUnknown Value = 0
	ValueFalse   Value = -1
)

type Assignment []Value

func NewAssignment(numVariable int) Assignment {
	return make(Assignment, numVariable+1)
}

type solverCtx struct {
	ctx context.Context
	r   Value
	a   Assignment
}

func (c *solverCtx) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *solverCtx) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *solverCtx) Err() error {
	return c.ctx.Err()
}

func (c *solverCtx) Value(key interface{}) interface{} {
	if ctxKey, ok := key.(ContextKey); ok {
		switch ctxKey {
		case ContextKeySatisfiable:
			return c.r
		case ContextKeyAssignment:
			return c.a
		default:
			return nil
		}
	}
	return c.ctx.Value(key)
}
