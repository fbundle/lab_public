package priority_queue

import "container/heap"

type Item[T any] struct {
	Value    T
	Priority int
	index    int // The index this item in the heap.
}

type itemList[T any] []*Item[T] // implement heap.Interface

func (l itemList[T]) Len() int {
	return len(l)
}

func (l itemList[T]) Less(i int, j int) bool {
	return l[i].Priority < l[j].Priority
}

func (l itemList[T]) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
	l[i].index = i // update index after swap
	l[j].index = j // update index after swap
}

func (l *itemList[T]) Push(x interface{}) {
	item := x.(*Item[T])
	item.index = len(*l)
	*l = append(*l, item)
}

func (l *itemList[T]) Pop() interface{} {
	out := (*l)[len(*l)-1]
	out.index = -1 // for safety

	(*l)[len(*l)-1] = nil // avoid memory leak
	*l = (*l)[:len(*l)-1]

	return out
}

type Queue[T any] struct {
	itemList itemList[T]
}

func (q *Queue[T]) Push(item *Item[T]) {
	heap.Push(&q.itemList, item)
}

func (q *Queue[T]) Update(item *Item[T]) {
	heap.Fix(&q.itemList, item.index)
}

func (q *Queue[T]) Pop() *Item[T] {
	if q.Len() == 0 {
		return nil
	}
	return heap.Pop(&q.itemList).(*Item[T])
}

func (q *Queue[T]) Len() int {
	return len(q.itemList)
}

func (q *Queue[T]) Peek() *Item[T] {
	if q.Len() == 0 {
		return nil
	}
	return q.itemList[0]
}

func Empty[T any]() *Queue[T] {
	return &Queue[T]{
		itemList: nil,
	}
}
