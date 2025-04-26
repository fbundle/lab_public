package priority_queue

type Item struct {
	Value    interface{}
	Priority int
	index    int // The index this item in the heap.
}

type Queue interface {
	Push(item *Item)
	Update(item *Item)
	Pop() *Item
	Len() int
	Peek() *Item
}

func New(items ...*Item) Queue {
	q := make([]*Item, 0, len(items))
	copy(q, items)
	// update index
	for i := range q {
		q[i].index = i
	}
	return (*queue)(&q)
}
