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
		root: &Node[string]{
			leaf:    false,
			entries: []*Entry[string]{{key: "Q", value: 2}},
			childs: []*Node[string]{{
				leaf:    false,
				entries: []*Entry[string]{{key: "F", value: 0}, {key: "K", value: 3}},
				childs: []*Node[string]{{
					leaf:    true,
					entries: []*Entry[string]{{key: "C", value: 4}},
					childs:  []*Node[string]{},
				}, {
					leaf:    true,
					entries: []*Entry[string]{{key: "H", value: 6}},
					childs:  []*Node[string]{},
				}, {
					leaf:    true,
					entries: []*Entry[string]{{key: "L", value: 5}, {key: "M", value: 10}, {key: "N", value: 12}},
					childs:  []*Node[string]{},
				}},
			}, {
				leaf:    false,
				entries: []*Entry[string]{{key: "T", value: 7}},
				childs: []*Node[string]{
					{
						leaf:    true,
						entries: []*Entry[string]{{key: "R", value: 11}, {key: "S", value: 1}},
						childs:  []*Node[string]{},
					},
					{
						leaf:    true,
						entries: []*Entry[string]{{key: "V", value: 8}, {key: "W", value: 9}},
						childs:  []*Node[string]{},
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

	var check func(got *Node[string], expected *Node[string])
	check = func(got *Node[string], expected *Node[string]) {
		if got.leaf != expected.leaf {
			t.Fatalf(
				"expected leaf=%v: got=%v, expected=%v",
				expected.leaf, got, expected)
		}

		if !slices.EqualFunc(
			got.entries, expected.entries,
			func(g *Entry[string], e *Entry[string]) bool {
				return g.key == e.key && e.value == e.value
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
