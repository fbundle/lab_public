package time_queue

import (
	"context"

	"sync"
	"time"

	pq "github.com/fbundle/go_util/pkg/priority_queue"
)

type Item struct {
	Schedule time.Time
	Value    interface{}
}

type Queue interface {
	// Put : put item into the time queue, return update function
	Put(item *Item) (dispatchUpdater func(updater func(*Item)))
	// DispatchLoop : run handle(item) to items in the time queue
	// at scheduled time or put time whichever is sooner
	// @Note : can be run concurrently
	// @Note : Loop stops when context is cancelled, need to call Flush to dispatch all items before the current time
	DispatchLoop(ctx context.Context, handle func(item *Item))
	// Flush : Flush is the last call made to Queue to flush all items before the current time
	Flush(now time.Time, handle func(item *Item))
}

func New() Queue {
	return &queue{
		mu:   &sync.Mutex{},
		pq:   pq.New(),
		wait: nil,
	}
}

type queue struct {
	mu   *sync.Mutex
	pq   pq.Queue // protected by mu
	wait *cv      // protected by mu
}

func (q *queue) Put(item *Item) (dispatchUpdater func(updater func(*Item))) {
	if item == nil {
		return nil
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	i := &pq.Item{
		Value:    item,
		Priority: int(item.Schedule.UnixNano()), // priority always positive
	}
	q.pq.Push(i)
	if q.wait != nil {
		// if DispatchLoop is waiting for new timeValue signal
		q.wait.dispose()
	}
	return func(updater func(*Item)) {
		q.mu.Lock()
		defer q.mu.Unlock()
		updater(item)
		// set priority to zero
		if newPriority := int(item.Schedule.UnixNano()); i.Priority != newPriority {
			i.Priority = newPriority
			q.pq.Update(i)
		}
	}
}

func (q *queue) Flush(now time.Time, handle func(item *Item)) {
	for {
		item := func() *Item {
			q.mu.Lock()
			defer q.mu.Unlock()
			i := q.pq.Peek()
			if i == nil {
				// pq empty
				return nil
			}
			item := i.Value.(*Item)
			if item.Schedule.After(now) {
				// exceeded
				return nil
			}
			q.pq.Pop()
			return item
		}()
		if item == nil {
			break
		}
		handle(item)
	}
}

func (q *queue) DispatchLoop(ctx context.Context, handle func(item *Item)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// peek the priority queue
		// either
		// pop the item (nil, item)
		// wait (wait, nil)
		q.mu.Lock()
		wait, item := func() (wait *cv, item *Item) {
			i := q.pq.Peek()
			if i == nil {
				// no item in queue
				// wait until new item
				ctx, cancel := context.WithCancel(ctx)
				return &cv{ctx, cancel}, nil
			}

			item = i.Value.(*Item)
			if item.Schedule.Before(time.Now()) {
				// item starts right now
				q.pq.Pop()
				return nil, item
			} else {
				// item not time right now
				// wait until time or new item
				ctx, cancel := context.WithDeadline(ctx, item.Schedule)
				return &cv{ctx, cancel}, nil
			}
		}()
		q.wait = wait // set wait for putTime to send signal
		q.mu.Unlock()

		if wait == nil {
			handle(item)
		} else {
			// child context of ctx
			// will stop if ctx is cancelled
			<-wait.done()
		}
	}
}

type cv struct {
	ctx    context.Context
	cancel func()
}

func (cv *cv) done() <-chan struct{} {
	return cv.ctx.Done()
}
func (cv *cv) dispose() {
	cv.cancel()
}
