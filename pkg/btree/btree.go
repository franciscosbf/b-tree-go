package btree

import (
	"cmp"
	"fmt"
	"slices"
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

func (bt *BTree[K]) insertNonNull(n *node[K], k K, v any) {
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

func (bt *BTree[K]) deleteAtLeafNode(n *node[K], k K) any {
	for i, entry := range n.entries {
		if k == entry.k {
			v := entry.v

			n.entries = slices.Delete(n.entries, i, i+1)

			return v
		}
	}

	return nil
}

func (bt *BTree[K]) deleteAtInternalNode(n *node[K], i int) any {
	v := n.entries[i].v

	pc := n.childs[i]
	fc := n.childs[i+1]
	switch {
	case len(pc.entries) >= bt.t:
		pe := pc.entries[len(pc.entries)-1]
		bt.delete(pc, pe.k)
		n.entries[i] = pe
	case len(pc.entries) == bt.t-1 && len(fc.entries) >= bt.t:
		fe := fc.entries[len(fc.entries)-1]
		bt.delete(fc, fe.k)
		n.entries[i] = fe
	default:
		pc.entries = append(pc.entries, fc.entries...)
		pc.childs = append(pc.childs, fc.childs...)
		n.entries = slices.Delete(n.entries, i, i+1)
		n.childs = slices.Delete(n.childs, i+1, i+1+1)

		if len(n.entries) == 0 && bt.root == n {
			bt.root = pc
		}
	}

	return v
}

func (bt *BTree[K]) deleteBalance(n *node[K], k K) any {
	i := len(n.entries) - 1
	for ; i >= 0 && k < n.entries[i].k; i-- {
	}
	i++

	if len(n.childs[i].entries) == bt.t-1 {
		switch {
		case i-1 > 0 && len(n.childs[i-1].entries) >= bt.t:
			n.childs[i].entries = append(
				[]*entry[K]{n.entries[i-1]},
				n.childs[i].entries...)
			n.entries[i] = n.childs[i-1].entries[len(n.childs[i-1].entries)-1]
			n.childs[i-1].entries = n.childs[i-1].entries[:len(n.childs[i-1].entries)-1]
			if !n.childs[i-1].leaf {
				n.childs[i].childs = append(
					[]*node[K]{n.childs[i-1].childs[len(n.childs[i-1].childs)-1]},
					n.childs[i+1].childs...)
				n.childs[i-1].childs = n.childs[i-1].childs[:len(n.childs[i-1].childs)-1]
			}
		case i+1 < len(n.childs) && len(n.childs[i+1].entries) >= bt.t:
			n.childs[i].entries = append(n.childs[i].entries, n.entries[i])
			n.entries[i] = n.childs[i+1].entries[0]
			n.childs[i+1].entries = n.childs[i+1].entries[1:]
			if !n.childs[i+1].leaf {
				n.childs[i].childs = append(n.childs[i].childs, n.childs[i+1].childs[0])
				n.childs[i+1].childs = n.childs[i+1].childs[1:]
			}
		case (i-1 > 0) && len(n.childs[i-1].entries) == bt.t-1 && (i+1 < len(n.childs)) && len(n.childs[i+1].entries) == bt.t-1:
			pc := n.childs[i-1]
			median := n.entries[i-1]
			pc.entries = append(
				pc.entries,
				append([]*entry[K]{median}, n.childs[i].entries...)...)
			pc.childs = append(pc.childs, n.childs[i].childs...)
			n.entries = slices.Delete(n.entries, i-1, i-1+1)
			n.childs = slices.Delete(n.childs, i, i+1)

			if len(n.entries) == 0 && bt.root == n {
				bt.root = pc
				n = bt.root
			}
		}

		i = len(n.entries) - 1
		for ; i >= 0 && k < n.entries[i].k; i-- {
		}
		i++
	}

	return bt.delete(n.childs[i], k)
}

func (bt *BTree[K]) deleteTraverse(n *node[K], k K) any {
	for i, entry := range n.entries {
		if k == entry.k {
			return bt.deleteAtInternalNode(n, i)
		}
	}

	return bt.deleteBalance(n, k)
}

func (bt *BTree[K]) delete(n *node[K], k K) any {
	if n.leaf {
		return bt.deleteAtLeafNode(n, k)
	}

	return bt.deleteTraverse(n, k)
}

func (bt *BTree[K]) Search(k K) any {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	return bt.search(bt.root, k)
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

	return bt.delete(bt.root, k)
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
