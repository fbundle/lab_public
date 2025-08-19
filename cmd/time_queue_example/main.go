package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fbundle/go_util/pkg/time_queue"
)

func schedule(tq *time_queue.Queue[func()], cancel func()) {
	now := time.Now()
	tq.Schedule(time_queue.Item[func()]{
		Time: now,
		Value: func() {
			fmt.Println("run at 0")
		},
	})

	tq.Schedule(time_queue.Item[func()]{
		Time: now.Add(time.Second),
		Value: func() {
			fmt.Println("run at 1")
		},
	})
	tq.Schedule(time_queue.Item[func()]{
		Time: now.Add(time.Second * 2),
		Value: func() {
			fmt.Println("run at 2")
		},
	})
	tq.Schedule(time_queue.Item[func()]{
		Time: now.Add(time.Second * 3),
		Value: func() {
			fmt.Println("stop at 3")
			cancel()
		},
	})
}

func main() {
	tq := time_queue.New[func()]()

	ctx, cancel := context.WithCancel(context.Background())
	go schedule(tq, cancel)
	for item := range tq.Dispatch(ctx) {
		item.Value()
	}
	for item := range tq.Flush(time.Now()) {
		item.Value()
	}
}
