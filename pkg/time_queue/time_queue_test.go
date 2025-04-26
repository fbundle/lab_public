package time_queue_test

import (
	"ca/pkg/time_queue"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type record struct {
	scheduleTime time.Time // time scheduled to collect
	putTime      time.Time // time put into the priority_queue
	collectTime  time.Time // actual time collected
}

const eps = 10 * time.Millisecond

const interval = 20 * time.Millisecond

const start = time.Second

const slots = 1000

var recordList []*record

var zero time.Time

func init() {
	rand.Seed(1234)
	zero = time.Now().Add(start)
	recordList = make([]*record, 0, slots)
	for i := 0; i < slots; i++ {
		recordList = append(recordList, &record{
			scheduleTime: zero.Add(interval * time.Duration(i)),
		})
	}
	rand.Shuffle(len(recordList), func(i, j int) {
		recordList[i], recordList[j] = recordList[j], recordList[i]
	})
}

func TestTimer_Put(t *testing.T) {
	q := time_queue.New()
	go func() {
		time.Sleep(time.Until(zero))
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for _, record := range recordList {
			q.Put(&time_queue.Item{
				Schedule: record.scheduleTime,
				Value:    record,
			})
			record.putTime = time.Now()
			<-ticker.C // i-th record is scheduled at i x interval
		}
	}()
	counter := 0
	ctx, cancel := context.WithCancel(context.Background())
	q.DispatchLoop(ctx, func(item *time_queue.Item) {
		record := item.Value.(*record)
		record.collectTime = time.Now()
		counter++
		if counter >= slots {
			cancel()
		}
	})

	for i, r := range recordList {
		fmt.Printf("[%d] put %v schedule %v collect %v\n", i, r.putTime.Sub(zero), r.scheduleTime.Sub(zero), r.collectTime.Sub(zero))

		if r.putTime.Before(r.scheduleTime) {
			if r.collectTime.Sub(r.scheduleTime) > eps {
				t.Error("put before schedule, but collect too late")
				return
			}
			if r.scheduleTime.Sub(r.collectTime) > eps {
				t.Error("put before schedule, but collect too early")
				return
			}
		}
		if r.putTime.After(r.scheduleTime) {
			if r.collectTime.Sub(r.putTime) > eps {
				t.Error("put after schedule, but collect to late")
				return
			}
		}
	}
}
