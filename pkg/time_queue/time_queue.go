package time_queue

import (
	"context"
	"iter"

	"sync"
	"time"

	pq "github.com/fbundle/go_util/pkg/priority_queue"
)

type Item[T any] struct {
	Time  time.Time
	Value T
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		mu:   sync.Mutex{},
		pq:   pq.Empty[Item[T]](),
		wait: nil,
	}
}

type Queue[T any] struct {
	mu   sync.Mutex
	pq   *pq.Queue[Item[T]] // protected by mu
	wait *cv                // protected by mu
}

func (q *Queue[T]) Schedule(item Item[T]) (dispatchUpdater func(updater func(Item[T]) Item[T])) {
	q.mu.Lock()
	defer q.mu.Unlock()
	i := &pq.Item[Item[T]]{
		Value:    item,
		Priority: int(item.Time.UnixNano()), // priority always positive
	}
	q.pq.Push(i)
	if q.wait != nil {
		// if DispatchLoop is waiting for new timeValue signal
		q.wait.dispose()
	}
	return func(updater func(Item[T]) Item[T]) {
		q.mu.Lock()
		defer q.mu.Unlock()
		newItem := updater(item)

		i.Value = newItem

		newPriority := int(newItem.Time.UnixNano())
		if newPriority != i.Priority {
			i.Priority = newPriority
			q.pq.Update(i) // update priority if it changes
		}
	}
}

func (q *Queue[T]) Flush(now time.Time) iter.Seq[Item[T]] {
	return func(yield func(Item[T]) bool) {
		for {
			item, ok := func() (Item[T], bool) {
				q.mu.Lock()
				defer q.mu.Unlock()
				i := q.pq.Peek()
				if i == nil { // pq empty
					return zero[Item[T]](), false
				}
				item := i.Value
				if item.Time.After(now) { // all items after after now
					return zero[Item[T]](), false
				}
				q.pq.Pop()
				return item, true
			}()
			if !ok {
				break
			}
			ok = yield(item)
			if !ok {
				break
			}
		}
	}

}

func (q *Queue[T]) Dispatch(ctx context.Context) iter.Seq[Item[T]] {
	return func(yield func(Item[T]) bool) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			q.mu.Lock()
			wait, item, ok := func() (wait *cv, item Item[T], ok bool) {
				i := q.pq.Peek()
				if i == nil { // no item in queue, wait until new item
					return cvWithCancel(ctx), zero[Item[T]](), false
				}

				item = i.Value
				if item.Time.Before(time.Now()) { // item starts right now
					q.pq.Pop()
					return nil, item, true
				} else { // item not time right now, wait until time or new item
					return cvWithDeadline(ctx, item.Time), zero[Item[T]](), false
				}
			}()
			q.wait = wait // set wait for putTime to send signal
			q.mu.Unlock()

			if ok {
				ok = yield(item)
				if !ok {
					return
				}
			} else {
				<-q.wait.done() // wait for next item
			}
		}
	}
}

func zero[T any]() T {
	var z T
	return z
}

func cvWithDeadline(ctx context.Context, deadline time.Time) *cv {
	ctx, cancel := context.WithDeadline(ctx, deadline)
	return &cv{ctx: ctx, cancel: cancel}
}

func cvWithCancel(ctx context.Context) *cv {
	ctx, cancel := context.WithCancel(ctx)
	return &cv{ctx: ctx, cancel: cancel}
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
