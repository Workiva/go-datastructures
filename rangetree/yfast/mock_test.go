package yfast

type mockEntry struct {
	values []uint64
}

func (me *mockEntry) ValueAtDimension(dimension uint64) uint64 {
	return me.values[dimension]
}

func newMockEntry(values ...uint64) *mockEntry {
	return &mockEntry{
		values: values,
	}
}
