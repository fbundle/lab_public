package transport

import (
	"github.com/khanh-nguyen-code/go_util/pkg/rlog/rpc"
)

type Transport interface {
	ListenAndServe(handler rpc.Handler) error
	Router(receiver rpc.Address) rpc.Handler
	Close() error
}
