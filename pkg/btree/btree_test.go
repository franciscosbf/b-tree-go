package btree

import (
	"slices"
	"strings"
	"testing"
)

var sampleKeys = strings.Split("F S Q K C L H T V W M R N", " ")

func newSampleWithMinimumDegree2() *BTree[string] {
	return &BTree[string]{
		t: 2,
		root: &node[string]{
			leaf:    false,
			entries: []*entry[string]{{k: "Q", v: 2}},
			childs: []*node[string]{{
				leaf:    false,
				entries: []*entry[string]{{k: "F", v: 0}, {k: "K", v: 3}},
				childs: []*node[string]{{
					leaf:    true,
					entries: []*entry[string]{{k: "C", v: 4}},
					childs:  []*node[string]{},
				}, {
					leaf:    true,
					entries: []*entry[string]{{k: "H", v: 6}},
					childs:  []*node[string]{},
				}, {
					leaf:    true,
					entries: []*entry[string]{{k: "L", v: 5}, {k: "M", v: 10}, {k: "N", v: 12}},
					childs:  []*node[string]{},
				}},
			}, {
				leaf:    false,
				entries: []*entry[string]{{k: "T", v: 7}},
				childs: []*node[string]{
					{
						leaf:    true,
						entries: []*entry[string]{{k: "R", v: 11}, {k: "S", v: 1}},
						childs:  []*node[string]{},
					},
					{
						leaf:    true,
						entries: []*entry[string]{{k: "V", v: 8}, {k: "W", v: 9}},
						childs:  []*node[string]{},
					},
				},
			}},
		},
	}
}

func TestSearch(t *testing.T) {
	bt := newSampleWithMinimumDegree2()

	for key, expectedValue := range map[string]any{"Q": 2, "K": 3, "S": 1} {
		value := bt.Search(key)

		if value == nil {
			t.Fatalf("didn't find key \"%v\"", key)
		}

		if value != expectedValue {
			t.Fatalf(
				"got different value for key \"%v\": got=%v expected=%v",
				key, value, expectedValue)
		}
	}
}

func TestInsertion(t *testing.T) {
	bt := New[string](2)

	for i, key := range sampleKeys {
		bt.Insert(key, i)
	}

	btSample := newSampleWithMinimumDegree2()

	var check func(got *node[string], expected *node[string])
	check = func(got *node[string], expected *node[string]) {
		if got.leaf != expected.leaf {
			t.Fatalf(
				"expected leaf=%v: got=%v, expected=%v",
				expected.leaf, got, expected)
		}

		if !slices.EqualFunc(
			got.entries, expected.entries,
			func(g *entry[string], e *entry[string]) bool {
				return g.k == e.k && e.v == e.v
			},
		) {
			t.Fatalf(
				"entries aren't equal: got=%v, expected=%v",
				got.entries, expected.entries)
		}

		if len(got.childs) != len(expected.childs) {
			t.Fatalf(
				"different number of childs: got=%v, expected=%v",
				got.childs, expected.childs)
		}

		for i, expectedChild := range expected.childs {
			check(got.childs[i], expectedChild)
		}
	}
	check(bt.root, btSample.root)
}
