package skip

type mockEntry struct {
	key uint64
}

func (me *mockEntry) Key() uint64 {
	return me.key
}

func newMockEntry(key uint64) *mockEntry {
	return &mockEntry{key}
}
