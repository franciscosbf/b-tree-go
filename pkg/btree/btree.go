package btree

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"sync"
)

type Entry[K cmp.Ordered] struct {
	key   K
	value any
}

func (e *Entry[K]) String() string {
	return fmt.Sprintf(
		"&Entry[string]{key:\"%v\",value:%v}",
		e.key, e.value)
}

type Node[K cmp.Ordered] struct {
	leaf    bool
	entries []*Entry[K]
	childs  []*Node[K]
}

func (n *Node[K]) String() string {
	return fmt.Sprintf(
		"&Node[string]{leaf:%v,entries:[]*Entry[string]%v,childs:[]*Node[string]%v}",
		n.leaf,
		strings.NewReplacer("[", "{", "]", "}", " ", ",").Replace(fmt.Sprint(n.entries)),
		strings.NewReplacer("[", "{", "]", "}", " ", ",").Replace(fmt.Sprint(n.childs)))
}

type BTree[K cmp.Ordered] struct {
	mutex sync.RWMutex
	t     int
	root  *Node[K]
}

func (bt *BTree[K]) isFull(node *Node[K]) bool {
	return len(node.entries) == (2*bt.t)-1
}

func (bt *BTree[K]) search(node *Node[K], key K) any {
	entries := node.entries
	i := 0

	for ; i < len(entries) && key > entries[i].key; i++ {
	}

	if i < len(entries) && key == entries[i].key {
		return entries[i].value
	}

	if node.leaf {
		return nil
	}

	return bt.search(node.childs[i], key)
}

func (bt *BTree[K]) Search(key K) any {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	return bt.search(bt.root, key)
}

func (bt *BTree[K]) splitChild(node *Node[K], i int) {
	left := node.childs[i]
	right := &Node[K]{leaf: left.leaf}

	median := left.entries[bt.t-1]

	right.entries = slices.Clone(left.entries[bt.t:])
	left.entries = left.entries[:bt.t-1]
	if !left.leaf {
		right.childs = slices.Clone(left.childs[bt.t:])
		left.childs = left.childs[:bt.t]
	}

	node.entries = slices.Insert(
		node.entries, i, median)
	node.childs = slices.Insert(node.childs, i+1, right)
}

func (bt *BTree[K]) splitRoot() {
	bt.root = &Node[K]{
		childs: []*Node[K]{bt.root},
	}

	bt.splitChild(bt.root, 0)
}

func (bt *BTree[K]) insertNonNull(
	node *Node[K],
	key K, value any,
) {
	i := len(node.entries) - 1
	for ; i >= 0 && key < node.entries[i].key; i-- {
	}
	i++

	if node.leaf {
		node.entries = slices.Insert(
			node.entries, i, &Entry[K]{key: key, value: value})

		return
	}

	if bt.isFull(node.childs[i]) {
		bt.splitChild(node, i)
		if key > node.entries[i].key {
			i++
		}
	}
	bt.insertNonNull(node.childs[i], key, value)
}

func (bt *BTree[K]) Insert(key K, value any) {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	if bt.isFull(bt.root) {
		bt.splitRoot()
	}

	bt.insertNonNull(bt.root, key, value)
}

func (bt *BTree[K]) Delete(key K) any {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	return nil // TODO: implement
}

func (bt *BTree[K]) String() string {
	return fmt.Sprintf("&BTree[string]{root:%v}", bt.root)
}

func New[K cmp.Ordered](
	minimumDegree int,
) *BTree[K] {
	if minimumDegree < 2 {
		panic("minimumDegree must be greater than 1")
	}

	return &BTree[K]{
		t:    minimumDegree,
		root: &Node[K]{leaf: true},
	}
}
