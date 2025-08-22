package priority_queue_test

import (
	"fmt"
	"testing"

	pq "github.com/fbundle/go_util/pkg/priority_queue"
)

func TestQueue(t *testing.T) {
	q := pq.Empty[int]()
	q.Push(&pq.Item[int]{
		Value:    2,
		Priority: 1,
	})
	q.Push(&pq.Item[int]{
		Value:    3,
		Priority: 3,
	})
	i := &pq.Item[int]{
		Value:    4,
		Priority: 2,
	}
	q.Push(i)

	fmt.Println(q.Pop())
	i.Priority = 5
	q.Update(i)

	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
}
