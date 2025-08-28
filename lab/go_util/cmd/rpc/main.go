package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fbundle/lab_public/lab/go_util/pkg/rpc"
)

func mustMarshalJSON(o any) string {
	b, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func main() {
	type AddReq struct {
		Values []int
	}

	type AddRes struct {
		Sum int
	}

	type SubReq struct {
		A int
		B int
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	addr := "localhost:14001"
	s, err := rpc.NewTCPServer(addr)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	d := rpc.NewDispatcher().Register("add", func(req *AddReq) (res *AddRes) {
		sum := 0
		for _, v := range req.Values {
			sum += v
		}
		return &AddRes{
			Sum: sum,
		}
	}).Register("sub", func(req *SubReq) (res *int) {
		diff := req.A - req.B
		return &diff
	})

	go s.ListenAndServe(ctx, d, rpc.NewMessageIO())

	localTransport := d.Handle
	remoteTransport := rpc.TCPTransport(ctx, addr, rpc.NewMessageIO())

	res1, err := rpc.RPC[AddReq, AddRes](localTransport, "add", &AddReq{Values: []int{1, 2, 3}})
	if err != nil {
		panic(err)
	}
	fmt.Println(mustMarshalJSON(res1))

	res2, err := rpc.RPC[SubReq, int](remoteTransport, "sub", &SubReq{A: 7, B: 5})
	if err != nil {
		panic(err)
	}
	fmt.Println(mustMarshalJSON(res2))
}
