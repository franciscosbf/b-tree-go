package btree

import (
	"slices"
	"strings"
	"testing"
)

var staticSample = struct {
	keys []string
	bt   *BTree[string]
}{
	keys: strings.Split("F S Q K C L H T V W M R N", " "),
	bt: &BTree[string]{
		t: 2,
		root: &node[string]{
			entries: []*entry[string]{{k: "Q", v: 2}},
			childs: []*node[string]{
				{
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
				},
			},
		},
	},
}

func checkTree(t *testing.T, got *node[string], expected *node[string]) {
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
		checkTree(t, got.childs[i], expectedChild)
	}
}

func TestSearch(t *testing.T) {
	btSample := staticSample.bt

	for key, expectedValue := range map[string]any{"Q": 2, "K": 3, "S": 1} {
		value := btSample.Search(key)

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

	for i, key := range staticSample.keys {
		bt.Insert(key, i)
	}

	btSample := staticSample.bt

	checkTree(t, bt.root, btSample.root)
}

func TestDeletion(t *testing.T) {
	samples := struct {
		bt       *BTree[string]
		outcomes []struct {
			keyToRemove   string
			expectedValue int
			expectedBt    *BTree[string]
		}
	}{
		bt: &BTree[string]{
			t: 3,
			root: &node[string]{
				entries: []*entry[string]{{"P", 1}},
				childs: []*node[string]{
					{
						entries: []*entry[string]{{"C", 2}, {"G", 3}, {"M", 4}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"A", 5}, {"B", 6}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"D", 7}, {"E", 8}, {"F", 9}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"J", 10}, {"K", 11}, {"L", 12}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"N", 13}, {"O", 14}},
							},
						},
					},
					{
						entries: []*entry[string]{{"T", 15}, {"X", 16}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"U", 20}, {"V", 21}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
							},
						},
					},
				},
			},
		},
		outcomes: []struct {
			keyToRemove   string
			expectedValue int
			expectedBt    *BTree[string]
		}{
			{
				keyToRemove:   "F",
				expectedValue: 9,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"P", 1}},
						childs: []*node[string]{
							{
								entries: []*entry[string]{{"C", 2}, {"G", 3}, {"M", 4}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"A", 5}, {"B", 6}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"D", 7}, {"E", 8}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"J", 10}, {"K", 11}, {"L", 12}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"N", 13}, {"O", 14}},
									},
								},
							},
							{
								entries: []*entry[string]{{"T", 15}, {"X", 16}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"U", 20}, {"V", 21}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
									},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "M",
				expectedValue: 4,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"P", 1}},
						childs: []*node[string]{
							{
								entries: []*entry[string]{{"C", 2}, {"G", 3}, {"L", 12}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"A", 5}, {"B", 6}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"D", 7}, {"E", 8}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"J", 10}, {"K", 11}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"N", 13}, {"O", 14}},
									},
								},
							},
							{
								entries: []*entry[string]{{"T", 15}, {"X", 16}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"U", 20}, {"V", 21}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
									},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "G",
				expectedValue: 3,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"P", 1}},
						childs: []*node[string]{
							{
								entries: []*entry[string]{{"C", 2}, {"L", 12}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"A", 5}, {"B", 6}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"D", 7}, {"E", 8}, {"J", 10}, {"K", 11}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"N", 13}, {"O", 14}},
									},
								},
							},
							{
								entries: []*entry[string]{{"T", 15}, {"X", 16}},
								childs: []*node[string]{
									{
										leaf:    true,
										entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"U", 20}, {"V", 21}},
									},
									{
										leaf:    true,
										entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
									},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "D",
				expectedValue: 7,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"C", 2}, {"L", 12}, {"P", 1}, {"T", 15}, {"X", 16}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"A", 5}, {"B", 6}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"E", 8}, {"J", 10}, {"K", 11}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"N", 13}, {"O", 14}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"U", 20}, {"V", 21}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "B",
				expectedValue: 6,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"E", 8}, {"L", 12}, {"P", 1}, {"T", 15}, {"X", 16}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"A", 5}, {"C", 2}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"J", 10}, {"K", 11}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"N", 13}, {"O", 14}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"U", 20}, {"V", 21}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "O",
				expectedValue: 14,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"E", 8}, {"L", 12}, {"Q", 17}, {"T", 15}, {"X", 16}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"A", 5}, {"C", 2}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"J", 10}, {"K", 11}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"N", 13}, {"P", 1}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"R", 18}, {"S", 19}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"U", 20}, {"V", 21}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "L",
				expectedValue: 12,
				expectedBt: &BTree[string]{
					t: 3,
					root: &node[string]{
						entries: []*entry[string]{{"E", 8}, {"Q", 17}, {"T", 15}, {"X", 16}},
						childs: []*node[string]{
							{
								leaf:    true,
								entries: []*entry[string]{{"A", 5}, {"C", 2}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"J", 10}, {"K", 11}, {"N", 13}, {"P", 1}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"R", 18}, {"S", 19}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"U", 20}, {"V", 21}},
							},
							{
								leaf:    true,
								entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
							},
						},
					},
				},
			},
		},
	}

	for _, outcome := range samples.outcomes {
		t.Logf("Testing deletion of key %v...", outcome.keyToRemove)

		value := samples.bt.Delete(outcome.keyToRemove)

		if value == nil {
			t.Fatalf("didn't find key \"%v\"", outcome.keyToRemove)
		}

		if value != outcome.expectedValue {
			t.Fatalf(
				"got different value for key \"%v\": got=%v expected=%v",
				outcome.keyToRemove, value, outcome.expectedValue)
		}

		checkTree(t, samples.bt.root, outcome.expectedBt.root)
	}
}
