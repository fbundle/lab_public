package rlog

import (
	"context"
	"github.com/khanh-nguyen-code/go_util/pkg/rlog/rpc"
	"sync/atomic"
)

type requestConstructor func(receiver NodeId) func() (request interface{})

type responseHandler func(receiver NodeId) func(response interface{})

func batchRPC(router rpc.Router, cluster ClusterState, requestConstructor requestConstructor, responseHandler responseHandler) <-chan struct{} {
	done := make(chan struct{})
	doneCount := uint32(0)
	var ctx = context.Background()
	if cluster.RpcTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), cluster.RpcTimeout)
		defer cancel()
	}
	for receiver, address := range cluster.AddressBook {
		go func(receiver NodeId, address rpc.Address) {
			call := &rpc.Call{
				Request:  requestConstructor(receiver)(),
				Response: nil,
			}
			go rpc.WaitThenCancel(call, ctx)
			router(address)(call)
			<-call.Done()
			if responseHandler != nil {
				responseHandler(receiver)(call.Response)
			}
			if atomic.AddUint32(&doneCount, 1) == uint32(len(cluster.AddressBook)) {
				close(done)
			}
		}(receiver, address)
	}
	return done
}
