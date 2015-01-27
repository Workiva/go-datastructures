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

type mockInterval struct {
	lows, highs []int64
}

func (mi *mockInterval) LowAtDimension(dimension uint64) int64 {
	return mi.lows[dimension]
}

func (mi *mockInterval) HighAtDimension(dimension uint64) int64 {
	return mi.highs[dimension]
}

func newMockInterval(lows, highs []int64) *mockInterval {
	return &mockInterval{
		lows:  lows,
		highs: highs,
	}
}
