type RPC interface {
	Broadcast(message interface{})
	Recv() <-chan interface{}
}

func NewRPC()