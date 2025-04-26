package transport

import (
	"ca/pkg/relay"
	"ca/pkg/relay/proto/gen/relay_pb"
	"ca/pkg/rlog/rpc"
	rlog_codec "ca/pkg/rlog_util/codec"
	"ca/pkg/uuid"
	"encoding/json"
	"sync"
)

type message struct {
	Uuid    string `json:"uuid"`
	Payload []byte `json:"payload"`
}

func NewRelay(agent relay.Peer) Transport {
	return &relayTransport{
		agent:      agent,
		watchChMap: sync.Map{},
	}
}

type relayTransport struct {
	agent      relay.Peer
	watchChMap sync.Map // map[uuid]chan []byte
}

func (r *relayTransport) Close() error {
	return r.agent.Close()
}
func (r *relayTransport) ListenAndServe(handler rpc.Handler) error {
	return r.agent.DialAndServe(func(m *relay_pb.Message) {
		go func() {
			reqMsg := &message{}
			err := json.Unmarshal(m.Payload, reqMsg)
			if err != nil {
				return
			}
			if val, loaded := r.watchChMap.LoadAndDelete(reqMsg.Uuid); loaded {
				watchCh := val.(chan []byte)
				watchCh <- reqMsg.Payload
				close(watchCh)
				return
			}
			if isRequest(reqMsg.Uuid) {
				resBytes := func() []byte {
					req, err := rlog_codec.Unmarshal(reqMsg.Payload)
					if err != nil {
						return nil
					}
					call := &rpc.Call{
						Request:  req,
						Response: nil,
					}
					handler(call)
					<-call.Done()
					b, err := rlog_codec.Marshal(call.Response)
					if err != nil {
						return nil
					}
					return b
				}()
				resMsg := &message{
					Uuid:    toResponseUuid(reqMsg.Uuid),
					Payload: resBytes,
				}
				b, err := json.Marshal(resMsg)
				if err != nil {
					return
				}
				err = r.agent.Write(&relay_pb.Message{
					Receiver: m.Sender,
					Payload:  b,
				})
				if err != nil {
					return
				}
			}
		}()
	})
}

func (r *relayTransport) Router(receiver rpc.Address) rpc.Handler {
	return func(call *rpc.Call) {
		response := func() interface{} {
			reqBytes, err := rlog_codec.Marshal(call.Request)
			if err != nil {
				return nil
			}
			reqId := makeRequestUuid()
			resId := toResponseUuid(reqId)
			watchCh := make(chan []byte, 1)
			r.watchChMap.Store(resId, watchCh)
			defer func() {
				r.watchChMap.Delete(resId)
			}()
			reqMsg := &message{
				Uuid:    reqId,
				Payload: reqBytes,
			}
			b, err := json.Marshal(reqMsg)
			if err != nil {
				return nil
			}
			err = r.agent.Write(&relay_pb.Message{
				Receiver: receiver,
				Payload:  b,
			})
			if err != nil {
				return nil
			}
			select {
			case <-call.Done():
				return nil
			case resBytes := <-watchCh:
				response, err := rlog_codec.Unmarshal(resBytes)
				if err != nil {
					response = nil
					return nil
				}
				return response
			}
		}()
		call.Write(response)
	}
}

func isRequest(id string) bool {
	return id[:3] == "req"
}

func makeRequestUuid() string {
	return "req" + uuid.New()
}

func toResponseUuid(reqId string) string {
	return "res" + reqId[3:]
}
