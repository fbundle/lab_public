package transport

import (
	"ca/pkg/rlog/rpc"
)

type Transport interface {
	ListenAndServe(handler rpc.Handler) error
	Router(receiver rpc.Address) rpc.Handler
	Close() error
}
