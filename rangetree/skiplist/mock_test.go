package skiplist

type mockEntry struct {
	values []int64
}

func (me *mockEntry) ValueAtDimension(dimension uint64) int64 {
	return me.values[dimension]
}

func newMockEntry(values ...int64) *mockEntry {
	return &mockEntry{values: values}
}
