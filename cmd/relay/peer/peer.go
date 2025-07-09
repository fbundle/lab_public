package main

import (
	"bufio"
	"flag"
	"fmt"
	"go_util/pkg/relay"
	"go_util/pkg/relay/proto/gen/relay_pb"
	"os"
	"strings"
)

var name *string
var relayAddr *string

func init() {
	name = flag.String("name", "", "name of client")
	relayAddr = flag.String("relay", "127.0.0.1:5010", "address of relay")
	flag.Parse()
}

func main() {
	if *name == "" {
		panic("name is required")
	}
	peer, err := relay.NewPeer(*name, "", *relayAddr)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := peer.DialAndServe(func(m *relay_pb.Message) {
				fmt.Printf("[%s] %s\n", m.Sender, string(m.Payload))
			}); err != nil {
				fmt.Printf("[peer] error: %v\n", err)
			}
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		slice := strings.SplitN(text, " ", 2)
		if len(slice) != 2 {
			continue
		}
		err = peer.Write(&relay_pb.Message{
			Receiver: slice[0],
			Payload:  []byte(slice[1]),
		})
		if err != nil {
			fmt.Printf("error %s\n", err.Error())
		}
	}
}
