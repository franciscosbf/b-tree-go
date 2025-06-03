package btree

import (
	"cmp"
	"fmt"
	"sync"
)

type entry[K cmp.Ordered] struct {
	k K
	v any
}

func (e *entry[K]) String() string {
	return fmt.Sprintf(
		"entry{key: %v, value: %v}",
		e.k, e.v)
}

type node[K cmp.Ordered] struct {
	leaf    bool
	entries []*entry[K]
	childs  []*node[K]
}

func (n *node[K]) String() string {
	return fmt.Sprintf(
		"node{leaf: %v, entries: %v, childs: %v}",
		n.leaf, n.entries, n.childs)
}

type BTree[K cmp.Ordered] struct {
	mutex sync.RWMutex
	t     int
	root  *node[K]
}

func (bt *BTree[K]) isFull(n *node[K]) bool {
	return len(n.entries) == (2*bt.t)-1
}

func (bt *BTree[K]) search(n *node[K], k K) any {
	entries := n.entries
	i := 0

	for ; i < len(entries) && k > entries[i].k; i++ {
	}

	if i < len(entries) && k == entries[i].k {
		return entries[i].v
	}

	if n.leaf {
		return nil
	}

	return bt.search(n.childs[i], k)
}

func (bt *BTree[K]) Search(k K) any {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	return bt.search(bt.root, k)
}

func (bt *BTree[K]) splitChild(n *node[K], i int) {
	left := n.childs[i]
	right := &node[K]{leaf: left.leaf}

	median := left.entries[bt.t-1]

	right.entries = append(
		right.entries,
		left.entries[bt.t:]...)
	left.entries = left.entries[:bt.t-1]
	if !left.leaf {
		right.childs = append(
			right.childs,
			left.childs[bt.t:]...)
		left.childs = left.childs[:bt.t]
	}

	n.entries = append(
		n.entries[:i],
		append([]*entry[K]{median}, n.entries[i:]...)...)
	n.childs = append(
		n.childs[:i+1],
		append([]*node[K]{right}, n.childs[i+1:]...)...)
}

func (bt *BTree[K]) splitRoot() {
	bt.root = &node[K]{
		childs: []*node[K]{bt.root},
	}

	bt.splitChild(bt.root, 0)
}

func (bt *BTree[K]) insertNonNull(
	n *node[K],
	k K, v any,
) {
	i := len(n.entries) - 1
	for ; i >= 0 && k < n.entries[i].k; i-- {
	}
	i++

	if n.leaf {
		n.entries = append(
			n.entries[:i],
			append([]*entry[K]{{k: k, v: v}}, n.entries[i:]...)...)

		return
	}

	if bt.isFull(n.childs[i]) {
		bt.splitChild(n, i)
		if k > n.entries[i].k {
			i++
		}
	}
	bt.insertNonNull(n.childs[i], k, v)
}

func (bt *BTree[K]) Insert(k K, v any) {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	if bt.isFull(bt.root) {
		bt.splitRoot()
	}

	bt.insertNonNull(bt.root, k, v)
}

func (bt *BTree[K]) Delete(k K) any {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	return nil // TODO: implement
}

func (bt *BTree[K]) String() string {
	return fmt.Sprintf("BTree{root: %v}", bt.root)
}

func New[K cmp.Ordered](
	minimumDegree int,
) *BTree[K] {
	if minimumDegree < 2 {
		panic("minimumDegree must be greater than 1")
	}

	return &BTree[K]{
		t:    minimumDegree,
		root: &node[K]{leaf: true},
	}
}
