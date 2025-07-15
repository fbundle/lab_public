package wbt

const (
	DELTA = 3
)

type Key[T any] interface {
	Cmp(T) int
}

type node[T Key[T]] struct {
	weight uint
	height uint
	key    T
	left   *node[T]
	right  *node[T]
}

func makeNode[T Key[T]](key T, left *node[T], right *node[T]) *node[T] {
	return &node[T]{
		weight: 1 + weight(left) + weight(right),
		height: 1 + max(height(left), height(right)),
		key:    key,
		left:   left,
		right:  right,
	}
}

func height[T Key[T]](n *node[T]) uint {
	if n == nil {
		return 0
	}
	return n.height
}

func weight[T Key[T]](n *node[T]) uint {
	if n == nil {
		return 0
	}
	return n.weight
}

func balance[T Key[T]](n *node[T]) *node[T] {
	if n == nil {
		return nil
	}
	if weight(n.left)+weight(n.right) <= 1 {
		return n
	}
	if weight(n.left) > DELTA*weight(n.right) { // left is guaranteed to be non-nil
		// left rotate
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
		n1 := makeNode(n.key, lr, r)
		l1 := makeNode(l.key, ll, n1)
		return l1
	} else if DELTA*weight(n.left) < weight(n.right) { // right is guaranteed to be non-nil
		// right rotate
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
		n1 := makeNode(n.key, l, rl)
		r1 := makeNode(r.key, n1, rr)
		return r1
	} else {
		return n
	}
}

func get[T Key[T]](n *node[T], keyIn T) (keyOut T, ok bool) {
	if n == nil {
		return keyOut, false
	}
	cmp := n.key.Cmp(keyIn)
	switch {
	case cmp < 0:
		return get(n.right, keyIn)
	case cmp > 0:
		return get(n.left, keyIn)
	default:
		return n.key, true
	}
}

func set[T Key[T]](n *node[T], key T) *node[T] {
	if n == nil {
		return makeNode(key, nil, nil)
	}
	cmp := n.key.Cmp(key)
	switch {
	case cmp < 0:
		r1 := set(n.right, key)
		n1 := makeNode(n.key, n.left, r1)
		return balance(n1)
	case cmp > 0:
		l1 := set(n.left, key)
		n1 := makeNode(n.key, l1, n.right)
		return balance(n1)
	default:
		return makeNode(key, n.left, n.right)
	}
}

func iter[T Key[T]](n *node[T], f func(k T) bool) {
	if n == nil {
		return
	}
	iter(n.left, f)
	if !f(n.key) {
		return
	}
	iter(n.right, f)
}
