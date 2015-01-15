package yfast

type mockEntry struct {
	// not going to use mock here as it skews benchmarks
	key uint64
}

func (me *mockEntry) Key() uint64 {
	return me.key
}

func newMockEntry(key uint64) *mockEntry {
	return &mockEntry{key}
}
