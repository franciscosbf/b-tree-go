package btree

import (
	"bytes"
	"cmp"
	"encoding/gob"
	"fmt"
	"slices"
	"sync"
)

type entry[K cmp.Ordered] struct {
	K K
	V any
}

func (e *entry[K]) String() string {
	return fmt.Sprintf(
		"entry{K: %v, V: %v}",
		e.K, e.V)
}

type node[K cmp.Ordered] struct {
	Leaf    bool
	Entries []*entry[K]
	Childs  []*node[K]
}

func (n *node[K]) String() string {
	return fmt.Sprintf(
		"node{Leaf: %v, Entries: %v, Childs: %v}",
		n.Leaf, n.Entries, n.Childs)
}

type tree[K cmp.Ordered] struct {
	T    int
	Root *node[K]
}

type BTree[K cmp.Ordered] struct {
	mutex sync.RWMutex
	tree  tree[K]
}

func (bt *BTree[K]) isFull(n *node[K]) bool {
	return len(n.Entries) == (2*bt.tree.T)-1
}

func (bt *BTree[K]) search(n *node[K], k K) any {
	entries := n.Entries
	i := 0

	for ; i < len(entries) && k > entries[i].K; i++ {
	}

	if i < len(entries) && k == entries[i].K {
		return entries[i].V
	}

	if n.Leaf {
		return nil
	}

	return bt.search(n.Childs[i], k)
}

func (bt *BTree[K]) splitChild(n *node[K], i int) {
	left := n.Childs[i]
	right := &node[K]{Leaf: left.Leaf}

	median := left.Entries[bt.tree.T-1]

	right.Entries = append(
		right.Entries,
		left.Entries[bt.tree.T:]...)
	left.Entries = left.Entries[:bt.tree.T-1]
	if !left.Leaf {
		right.Childs = append(
			right.Childs,
			left.Childs[bt.tree.T:]...)
		left.Childs = left.Childs[:bt.tree.T]
	}

	n.Entries = append(
		n.Entries[:i],
		append([]*entry[K]{median}, n.Entries[i:]...)...)
	n.Childs = append(
		n.Childs[:i+1],
		append([]*node[K]{right}, n.Childs[i+1:]...)...)
}

func (bt *BTree[K]) splitRoot() {
	bt.tree.Root = &node[K]{
		Childs: []*node[K]{bt.tree.Root},
	}

	bt.splitChild(bt.tree.Root, 0)
}

func (bt *BTree[K]) findRawPos(n *node[K], k K) int {
	i := len(n.Entries) - 1
	for ; i >= 0 && k < n.Entries[i].K; i-- {
	}

	return i
}

func (bt *BTree[K]) findPos(n *node[K], k K) int {
	i := bt.findRawPos(n, k)
	i++

	return i
}

func (bt *BTree[K]) insertNonNull(n *node[K], k K, v any) {
	i := bt.findRawPos(n, k)

	if i >= 0 && k == n.Entries[i].K {
		n.Entries[i].V = v

		return
	}

	i++

	if n.Leaf {
		n.Entries = append(
			n.Entries[:i],
			append([]*entry[K]{{K: k, V: v}}, n.Entries[i:]...)...)

		return
	}

	if bt.isFull(n.Childs[i]) {
		bt.splitChild(n, i)

		if k > n.Entries[i].K {
			i++
		}
	}

	bt.insertNonNull(n.Childs[i], k, v)
}

func (bt *BTree[K]) deleteAtLeafNode(n *node[K], k K) any {
	for i, entry := range n.Entries {
		if k == entry.K {
			v := entry.V

			n.Entries = slices.Delete(n.Entries, i, i+1)

			return v
		}
	}

	return nil
}

func (bt *BTree[K]) deleteAtInternalNode(n *node[K], i int) any {
	v := n.Entries[i].V

	pc := n.Childs[i]
	fc := n.Childs[i+1]
	switch {
	case len(pc.Entries) >= bt.tree.T:
		pe := pc.Entries[len(pc.Entries)-1]

		bt.delete(pc, pe.K)

		n.Entries[i] = pe
	case len(pc.Entries) == bt.tree.T-1 && len(fc.Entries) >= bt.tree.T:
		fe := fc.Entries[len(fc.Entries)-1]

		bt.delete(fc, fe.K)

		n.Entries[i] = fe
	default:
		pc.Entries = append(pc.Entries, fc.Entries...)
		pc.Childs = append(pc.Childs, fc.Childs...)

		n.Entries = slices.Delete(n.Entries, i, i+1)
		n.Childs = slices.Delete(n.Childs, i+1, i+1+1)

		if len(n.Entries) == 0 && bt.tree.Root == n {
			bt.tree.Root = pc
		}
	}

	return v
}

func (bt *BTree[K]) deleteBalance(n *node[K], i int, k K) any {
	if len(n.Childs[i].Entries) == bt.tree.T-1 {
		ki := max(i-1, 0)

		im1, ip1 := i-1, i+1

		if im1 >= 0 && len(n.Childs[im1].Entries) >= bt.tree.T {
			n.Childs[i].Entries = append(
				[]*entry[K]{n.Entries[ki]},
				n.Childs[i].Entries...)
			n.Entries[ki] = n.Childs[im1].Entries[len(n.Childs[im1].Entries)-1]
			n.Childs[im1].Entries = n.Childs[im1].Entries[:len(n.Childs[im1].Entries)-1]

			if !n.Childs[im1].Leaf {
				n.Childs[i].Childs = append(
					[]*node[K]{n.Childs[im1].Childs[len(n.Childs[im1].Childs)-1]},
					n.Childs[ip1].Childs...)
				n.Childs[im1].Childs = n.Childs[im1].Childs[:len(n.Childs[im1].Childs)-1]
			}
		} else if ip1 < len(n.Childs) && len(n.Childs[ip1].Entries) >= bt.tree.T {
			if i >= 1 && i <= len(n.Entries) {
				ki++
			}

			n.Childs[i].Entries = append(n.Childs[i].Entries, n.Entries[ki])
			n.Entries[ki] = n.Childs[ip1].Entries[0]
			n.Childs[ip1].Entries = n.Childs[ip1].Entries[1:]

			if !n.Childs[ip1].Leaf {
				n.Childs[i].Childs = append(n.Childs[i].Childs, n.Childs[ip1].Childs[0])
				n.Childs[ip1].Childs = n.Childs[ip1].Childs[1:]
			}
		} else {
			var nn *node[K]

			if im1 >= 0 && len(n.Childs[im1].Entries) == bt.tree.T-1 {
				pc := n.Childs[im1]
				nn = pc
				median := n.Entries[ki]

				pc.Entries = append(
					pc.Entries,
					append([]*entry[K]{median}, n.Childs[i].Entries...)...)
				pc.Childs = append(pc.Childs, n.Childs[i].Childs...)

				n.Entries = slices.Delete(n.Entries, ki, ki+1)
				n.Childs = slices.Delete(n.Childs, i, i+1)
			} else if ip1 < len(n.Childs) && len(n.Childs[ip1].Entries) == bt.tree.T-1 {
				fc := n.Childs[ip1]
				nn = fc
				median := n.Entries[ki]

				fc.Entries = append(
					append(n.Childs[i].Entries, median),
					fc.Entries...)
				fc.Childs = append(n.Childs[i].Childs, fc.Childs...)

				n.Entries = slices.Delete(n.Entries, ki, ki+1)
				n.Childs = slices.Delete(n.Childs, i, i+1)
			} else {
				return bt.delete(n.Childs[i], k)
			}

			if len(n.Entries) == 0 && bt.tree.Root == n {
				bt.tree.Root = nn
				n = bt.tree.Root
			}
		}

		i = bt.findPos(n, k)
	}

	return bt.delete(n.Childs[i], k)
}

func (bt *BTree[K]) deleteTraverse(n *node[K], k K) any {
	i := bt.findRawPos(n, k)

	if i >= 0 && n.Entries[i].K == k {
		return bt.deleteAtInternalNode(n, i)
	}

	return bt.deleteBalance(n, i+1, k)
}

func (bt *BTree[K]) delete(n *node[K], k K) any {
	if n.Leaf {
		return bt.deleteAtLeafNode(n, k)
	}

	return bt.deleteTraverse(n, k)
}

func (bt *BTree[K]) Search(k K) any {
	bt.mutex.RLock()
	defer bt.mutex.RUnlock()

	return bt.search(bt.tree.Root, k)
}

func (bt *BTree[K]) Insert(k K, v any) {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	if bt.isFull(bt.tree.Root) {
		bt.splitRoot()
	}

	bt.insertNonNull(bt.tree.Root, k, v)
}

func (bt *BTree[K]) Delete(k K) any {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	return bt.delete(bt.tree.Root, k)
}

func EncodeGob[K cmp.Ordered](bt *BTree[K]) []byte {
	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	var buf bytes.Buffer

	gob.NewEncoder(&buf).Encode(bt.tree)

	return buf.Bytes()
}

func DecodeGob[K cmp.Ordered](raw []byte) (*BTree[K], error) {
	var buf bytes.Buffer
	buf.Write(raw)

	var bt BTree[K]
	if err := gob.NewDecoder(&buf).Decode(&bt); err != nil {
		return nil, err
	}

	return &bt, nil
}

func (bt *BTree[K]) String() string {
	return fmt.Sprintf("BTree{Root: %v}", bt.tree.Root)
}

func New[K cmp.Ordered](minimumDegree int) *BTree[K] {
	if minimumDegree < 2 {
		panic("minimumDegree must be at least 2")
	}

	return &BTree[K]{
		tree: tree[K]{
			T:    minimumDegree,
			Root: &node[K]{Leaf: true},
		},
	}
}
