package rpc

import (
	"context"
	"sync/atomic"
)

type Address = string

type Router func(receiver Address) Handler

type Handler func(call *Call)

type state = uint32

const (
	state_ready   state = 0
	state_written state = 1
)

type Call struct {
	Request  interface{}
	Response interface{}
	doneCh   atomic.Value // chan struct{}
	state    state
}

func (call *Call) initOnce() {
	call.doneCh.CompareAndSwap(nil, make(chan struct{}))
}

func (call *Call) Write(response interface{}) {
	call.initOnce()
	if atomic.CompareAndSwapUint32(&call.state, state_ready, state_written) {
		call.Response = response
		close(call.doneCh.Load().(chan struct{}))
	}
}

func (call *Call) Done() <-chan struct{} {
	call.initOnce()
	return call.doneCh.Load().(chan struct{})
}

func WaitThenCancel(call *Call, ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-call.Done():
		call.Write(nil)
	}
}
