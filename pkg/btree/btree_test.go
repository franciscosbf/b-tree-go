package btree

import (
	"slices"
	"testing"
)

func checkTree(t *testing.T, got *node[string], expected *node[string]) {
	if got.Leaf != expected.Leaf {
		t.Fatalf(
			"expected leaf=%v: got=%v, expected=%v",
			expected.Leaf, got, expected)
	}

	if !slices.EqualFunc(
		got.Entries, expected.Entries,
		func(g *entry[string], e *entry[string]) bool {
			return g.K == e.K && e.V == e.V
		},
	) {
		t.Fatalf(
			"entries aren't equal: got=%v, expected=%v",
			got.Entries, expected.Entries)
	}

	if len(got.Childs) != len(expected.Childs) {
		t.Fatalf(
			"different number of childs: got=%v, expected=%v",
			got.Childs, expected.Childs)
	}

	for i, expectedChild := range expected.Childs {
		checkTree(t, got.Childs[i], expectedChild)
	}
}

func TestSearch(t *testing.T) {
	bt := &BTree[string]{
		tree: tree[string]{
			T: 2,
			Root: &node[string]{
				Entries: []*entry[string]{{K: "Q", V: 2}},
				Childs: []*node[string]{
					{
						Entries: []*entry[string]{{K: "F", V: 0}, {K: "K", V: 3}},
						Childs: []*node[string]{{
							Leaf:    true,
							Entries: []*entry[string]{{K: "C", V: 4}},
							Childs:  []*node[string]{},
						}, {
							Leaf:    true,
							Entries: []*entry[string]{{K: "H", V: 6}},
							Childs:  []*node[string]{},
						}, {
							Leaf:    true,
							Entries: []*entry[string]{{K: "L", V: 5}, {K: "M", V: 10}, {K: "N", V: 12}},
							Childs:  []*node[string]{},
						}},
					}, {
						Entries: []*entry[string]{{K: "T", V: 7}},
						Childs: []*node[string]{
							{
								Leaf:    true,
								Entries: []*entry[string]{{K: "R", V: 11}, {K: "S", V: 1}},
								Childs:  []*node[string]{},
							},
							{
								Leaf:    true,
								Entries: []*entry[string]{{K: "V", V: 8}, {K: "W", V: 9}},
								Childs:  []*node[string]{},
							},
						},
					},
				},
			},
		},
	}

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
	expectedBt := &BTree[string]{
		tree: tree[string]{
			T: 2,
			Root: &node[string]{
				Entries: []*entry[string]{{K: "Q", V: 2}},
				Childs: []*node[string]{
					{
						Entries: []*entry[string]{{K: "F", V: 0}, {K: "K", V: 3}},
						Childs: []*node[string]{{
							Leaf:    true,
							Entries: []*entry[string]{{K: "C", V: 4}},
							Childs:  []*node[string]{},
						}, {
							Leaf:    true,
							Entries: []*entry[string]{{K: "H", V: 6}},
							Childs:  []*node[string]{},
						}, {
							Leaf:    true,
							Entries: []*entry[string]{{K: "L", V: 5}, {K: "M", V: 10}, {K: "N", V: 12}},
							Childs:  []*node[string]{},
						}},
					}, {
						Entries: []*entry[string]{{K: "T", V: 7}},
						Childs: []*node[string]{
							{
								Leaf:    true,
								Entries: []*entry[string]{{K: "R", V: 11}, {K: "S", V: 1}},
								Childs:  []*node[string]{},
							},
							{
								Leaf:    true,
								Entries: []*entry[string]{{K: "V", V: 8}, {K: "W", V: 9}},
								Childs:  []*node[string]{},
							},
						},
					},
				},
			},
		},
	}

	for i, key := range []string{"F", "S", "Q", "K", "C", "L", "H", "T", "V", "W", "M", "R", "N"} {
		bt.Insert(key, i)
	}

	checkTree(t, bt.tree.Root, expectedBt.tree.Root)
}

func TestDuplicatedInsertion(t *testing.T) {
	type testCase struct {
		keyToInsert   string
		valueToInsert int
		expectedBt    *BTree[string]
	}

	type testSample struct {
		bt    *BTree[string]
		cases []*testCase
	}

	ts := testSample{
		bt: &BTree[string]{
			tree: tree[string]{
				T: 2,
				Root: &node[string]{
					Entries: []*entry[string]{{"B", 1}, {"D", 2}},
					Childs: []*node[string]{
						{
							Leaf:    true,
							Entries: []*entry[string]{{"A", 3}},
						},
						{
							Leaf:    true,
							Entries: []*entry[string]{{"C", 4}},
						},
						{
							Leaf:    true,
							Entries: []*entry[string]{{"E", 5}},
						},
					},
				},
			},
		},
		cases: []*testCase{
			{
				keyToInsert:   "D",
				valueToInsert: 12,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 2,
						Root: &node[string]{
							Entries: []*entry[string]{{"B", 1}, {"D", 12}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 3}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"C", 4}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"E", 5}},
								},
							},
						},
					},
				},
			},
			{
				keyToInsert:   "C",
				valueToInsert: 14,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 2,
						Root: &node[string]{
							Entries: []*entry[string]{{"B", 1}, {"D", 12}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 3}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"C", 14}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"E", 5}},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range ts.cases {
		t.Logf(
			"Testing duplicated insertion with key-value pair (%v, %v)...",
			tc.keyToInsert, tc.valueToInsert)

		ts.bt.Insert(tc.keyToInsert, tc.valueToInsert)

		checkTree(t, ts.bt.tree.Root, tc.expectedBt.tree.Root)
	}
}

func TestDeletion(t *testing.T) {
	type testCase struct {
		keyToRemove   string
		expectedValue int
		expectedBt    *BTree[string]
	}

	type testSample struct {
		bt    *BTree[string]
		cases []*testCase
	}

	ts1 := testSample{
		bt: &BTree[string]{
			tree: tree[string]{
				T: 3,
				Root: &node[string]{
					Entries: []*entry[string]{{"P", 1}},
					Childs: []*node[string]{
						{
							Entries: []*entry[string]{{"C", 2}, {"G", 3}, {"M", 4}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 5}, {"B", 6}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"D", 7}, {"E", 8}, {"F", 9}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"J", 10}, {"K", 11}, {"L", 12}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"N", 13}, {"O", 14}},
								},
							},
						},
						{
							Entries: []*entry[string]{{"T", 15}, {"X", 16}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"U", 20}, {"V", 21}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
								},
							},
						},
					},
				},
			},
		},
		cases: []*testCase{
			{
				keyToRemove:   "F",
				expectedValue: 9,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"P", 1}},
							Childs: []*node[string]{
								{
									Entries: []*entry[string]{{"C", 2}, {"G", 3}, {"M", 4}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"A", 5}, {"B", 6}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"D", 7}, {"E", 8}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"J", 10}, {"K", 11}, {"L", 12}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"N", 13}, {"O", 14}},
										},
									},
								},
								{
									Entries: []*entry[string]{{"T", 15}, {"X", 16}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"U", 20}, {"V", 21}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
										},
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
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"P", 1}},
							Childs: []*node[string]{
								{
									Entries: []*entry[string]{{"C", 2}, {"G", 3}, {"L", 12}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"A", 5}, {"B", 6}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"D", 7}, {"E", 8}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"J", 10}, {"K", 11}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"N", 13}, {"O", 14}},
										},
									},
								},
								{
									Entries: []*entry[string]{{"T", 15}, {"X", 16}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"U", 20}, {"V", 21}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
										},
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
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"P", 1}},
							Childs: []*node[string]{
								{
									Entries: []*entry[string]{{"C", 2}, {"L", 12}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"A", 5}, {"B", 6}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"D", 7}, {"E", 8}, {"J", 10}, {"K", 11}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"N", 13}, {"O", 14}},
										},
									},
								},
								{
									Entries: []*entry[string]{{"T", 15}, {"X", 16}},
									Childs: []*node[string]{
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"U", 20}, {"V", 21}},
										},
										{
											Leaf:    true,
											Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
										},
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
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"C", 2}, {"L", 12}, {"P", 1}, {"T", 15}, {"X", 16}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 5}, {"B", 6}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"E", 8}, {"J", 10}, {"K", 11}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"N", 13}, {"O", 14}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"U", 20}, {"V", 21}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "B",
				expectedValue: 6,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"E", 8}, {"L", 12}, {"P", 1}, {"T", 15}, {"X", 16}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 5}, {"C", 2}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"J", 10}, {"K", 11}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"N", 13}, {"O", 14}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Q", 17}, {"R", 18}, {"S", 19}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"U", 20}, {"V", 21}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "O",
				expectedValue: 14,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"E", 8}, {"L", 12}, {"Q", 17}, {"T", 15}, {"X", 16}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 5}, {"C", 2}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"J", 10}, {"K", 11}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"N", 13}, {"P", 1}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"R", 18}, {"S", 19}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"U", 20}, {"V", 21}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
								},
							},
						},
					},
				},
			},
			{
				keyToRemove:   "L",
				expectedValue: 12,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Entries: []*entry[string]{{"E", 8}, {"Q", 17}, {"T", 15}, {"X", 16}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 5}, {"C", 2}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"J", 10}, {"K", 11}, {"N", 13}, {"P", 1}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"R", 18}, {"S", 19}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"U", 20}, {"V", 21}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"Y", 22}, {"Z", 23}},
								},
							},
						},
					},
				},
			},
		},
	}

	ts2 := testSample{
		bt: &BTree[string]{
			tree: tree[string]{
				T: 3,
				Root: &node[string]{
					Entries: []*entry[string]{{"L", 1}},
					Childs: []*node[string]{
						{
							Leaf:    true,
							Entries: []*entry[string]{{"A", 2}, {"B", 3}},
						},
						{
							Leaf:    true,
							Entries: []*entry[string]{{"E", 4}, {"J", 5}},
						},
					},
				},
			},
		},
		cases: []*testCase{
			{
				keyToRemove:   "L",
				expectedValue: 1,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Leaf:    true,
							Entries: []*entry[string]{{"A", 2}, {"B", 3}, {"E", 4}, {"J", 5}},
						},
					},
				},
			},
		},
	}

	ts3 := testSample{
		bt: &BTree[string]{
			tree: tree[string]{
				T: 3,
				Root: &node[string]{
					Leaf:    true,
					Entries: []*entry[string]{{"W", 1}},
				},
			},
		},
		cases: []*testCase{
			{
				keyToRemove:   "W",
				expectedValue: 1,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 3,
						Root: &node[string]{
							Leaf: true,
						},
					},
				},
			},
		},
	}

	ts4 := testSample{
		bt: &BTree[string]{
			tree: tree[string]{
				T: 2,
				Root: &node[string]{
					Entries: []*entry[string]{{"B", 1}, {"D", 2}},
					Childs: []*node[string]{
						{
							Leaf:    true,
							Entries: []*entry[string]{{"A", 3}},
						},
						{
							Leaf:    true,
							Entries: []*entry[string]{{"C", 4}},
						},
						{
							Leaf:    true,
							Entries: []*entry[string]{{"E", 5}},
						},
					},
				},
			},
		},
		cases: []*testCase{
			{
				keyToRemove:   "C",
				expectedValue: 4,
				expectedBt: &BTree[string]{
					tree: tree[string]{
						T: 2,
						Root: &node[string]{
							Entries: []*entry[string]{{"D", 2}},
							Childs: []*node[string]{
								{
									Leaf:    true,
									Entries: []*entry[string]{{"A", 3}, {"B", 1}},
								},
								{
									Leaf:    true,
									Entries: []*entry[string]{{"E", 5}},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, ts := range []testSample{ts1, ts2, ts3, ts4} {
		for _, tc := range ts.cases {
			t.Logf("Testing deletion of key %v (sample %v)...", tc.keyToRemove, i+1)

			value := ts.bt.Delete(tc.keyToRemove)

			if value == nil {
				t.Fatalf("didn't find key \"%v\"", tc.keyToRemove)
			}

			if value != tc.expectedValue {
				t.Fatalf(
					"got different value for key \"%v\": got=%v expected=%v",
					tc.keyToRemove, value, tc.expectedValue)
			}

			checkTree(t, ts.bt.tree.Root, tc.expectedBt.tree.Root)
		}
	}
}
