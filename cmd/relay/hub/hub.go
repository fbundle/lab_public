package main

import (
	"flag"
	"fmt"
	"github.com/fbundle/go_util/pkg/relay"
)

var listenAddr *string

func init() {
	listenAddr = flag.String("listen", ":5010", "listen address")
	flag.Parse()
}

func main() {
	hub, err := relay.NewHub(*listenAddr)
	if err != nil {
		panic(err)
	}
	for {
		if err := hub.ListenAndServe(); err != nil {
			fmt.Printf("[hub] error: %v\n", err)
		}
		if err := hub.Close(); err != nil {
			fmt.Printf("[hub] error: %v\n", err)
		}
	}
}
