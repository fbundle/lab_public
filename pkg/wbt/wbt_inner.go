package wbt

const (
	DELTA = 3
)

type Comparable[T any] interface {
	Cmp(T) int
}

type node[T Comparable[T]] struct {
	weight uint
	height uint
	entry  T
	left   *node[T]
	right  *node[T]
}

func makeNode[T Comparable[T]](entry T, left *node[T], right *node[T]) *node[T] {
	return &node[T]{
		weight: 1 + weight(left) + weight(right),
		height: 1 + max(height(left), height(right)),
		entry:  entry,
		left:   left,
		right:  right,
	}
}

func height[T Comparable[T]](n *node[T]) uint {
	if n == nil {
		return 0
	}
	return n.height
}

func weight[T Comparable[T]](n *node[T]) uint {
	if n == nil {
		return 0
	}
	return n.weight
}

func get[T Comparable[T]](n *node[T], entryIn T) (entryOut T, ok bool) {
	if n == nil {
		return entryOut, false
	}
	cmp := n.entry.Cmp(entryIn)
	switch {
	case cmp < 0:
		return get(n.right, entryIn)
	case cmp > 0:
		return get(n.left, entryIn)
	default:
		return n.entry, true
	}
}

func iter[T Comparable[T]](n *node[T], f func(k T) bool) {
	if n == nil {
		return
	}
	iter(n.left, f)
	if !f(n.entry) {
		return
	}
	iter(n.right, f)
}

func balance[T Comparable[T]](n *node[T]) *node[T] {
	if n == nil {
		return nil
	}
	if weight(n.left)+weight(n.right) <= 1 {
		return n
	}
	if weight(n.left) > DELTA*weight(n.right) { // left is guaranteed to be non-nil
		// right rotate
		//         n
		//   l           r
		// ll lr
		//
		//      becomes
		//
		//         l
		//   ll          n
		//             lr r

		l, r := n.left, n.right
		ll, lr := l.left, l.right
		n1 := makeNode(n.entry, lr, r)
		l1 := makeNode(l.entry, ll, n1)
		return l1
	} else if DELTA*weight(n.left) < weight(n.right) { // right is guaranteed to be non-nil
		// left rotate
		//         n
		//   l           r
		//             rl rr
		//
		//      becomes
		//
		//         r
		//   n          rr
		//  l rl

		l, r := n.left, n.right
		rl, rr := r.left, r.right
		n1 := makeNode(n.entry, l, rl)
		r1 := makeNode(r.entry, n1, rr)
		return r1
	} else {
		return n
	}
}

func set[T Comparable[T]](n *node[T], entryIn T) *node[T] {
	if n == nil {
		return makeNode(entryIn, nil, nil)
	}
	cmp := n.entry.Cmp(entryIn)
	switch {
	case cmp < 0:
		r1 := set(n.right, entryIn)
		n1 := makeNode(n.entry, n.left, r1)
		return balance(n1)
	case cmp > 0:
		l1 := set(n.left, entryIn)
		n1 := makeNode(n.entry, l1, n.right)
		return balance(n1)
	default:
		return makeNode(entryIn, n.left, n.right)
	}
}

func del[T Comparable[T]](n *node[T], entryIn T) *node[T] {
	if n == nil {
		return nil
	}
	cmp := n.entry.Cmp(entryIn)
	switch {
	case cmp < 0:
		r1 := del(n.right, entryIn)
		n1 := makeNode(n.entry, n.left, r1)
		return balance(n1)
	case cmp > 0:
		l1 := del(n.left, entryIn)
		n1 := makeNode(n.entry, l1, n.right)
		return balance(n1)
	default:
		return merge(n.left, n.right)
	}
}

func getMinEntry[T Comparable[T]](n *node[T]) T {
	if n == nil {
		panic("min of nil tree")
	}
	if n.left == nil {
		return n.entry
	} else {
		return getMinEntry(n.left)
	}
}

func merge[T Comparable[T]](l *node[T], r *node[T]) *node[T] {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	// merge chooses minimum from the right side
	// or maximum from left side but this is just a small optimization
	entry := getMinEntry(r)
	r1 := del(r, entry)
	n1 := makeNode(entry, l, r1)
	return balance(n1)
}

// split - ([1, 2, 3, 4], 3) -> [1, 2, 3] , [4]
func split[T Comparable[T]](n *node[T], entry T) (*node[T], *node[T]) {
	if n == nil {
		return nil, nil
	}
	cmp := n.entry.Cmp(entry)
	switch {
	case cmp < 0:
		rl1, rr1 := split(n.right, entry)
		n1 := makeNode(n.entry, n.left, rl1)
		n2 := balance(n1)
		return n2, rr1
	case cmp > 0:
		ll1, lr1 := split(n.left, entry)
		n1 := makeNode(n.entry, lr1, n.right)
		n2 := balance(n1)
		return ll1, n2
	default:
		return n.left, set(n.right, n.entry)
	}
}
