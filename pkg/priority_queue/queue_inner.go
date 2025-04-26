package priority_queue

import (
	"container/heap"
)

type items []*Item

func (items items) Len() int {
	return len(items)
}

func (items items) Less(i int, j int) bool {
	return items[i].Priority < items[j].Priority
}

func (items items) Swap(i int, j int) {
	items[i], items[j] = items[j], items[i]
	items[i].index = i // update index after swap
	items[j].index = j // update index after swap
}

func (items *items) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*items)
	*items = append(*items, item)
}

func (items *items) Pop() interface{} {
	out := (*items)[len(*items)-1]
	out.index = -1 // for safety

	(*items)[len(*items)-1] = nil // avoid memory leak
	*items = (*items)[:len(*items)-1]

	return out
}

type queue items

func (q *queue) Push(item *Item) {
	heap.Push((*items)(q), item)
}

func (q *queue) Update(item *Item) {
	heap.Fix((*items)(q), item.index)
}

func (q *queue) Pop() *Item {
	if q.Len() == 0 {
		return nil
	}
	return heap.Pop((*items)(q)).(*Item)
}

func (q *queue) Len() int {
	return len(*(*items)(q))
}

func (q *queue) Peek() *Item {
	if q.Len() == 0 {
		return nil
	}
	return (*(*items)(q))[0]
}
