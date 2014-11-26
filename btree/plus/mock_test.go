package plus

type mockKey struct {
	value, id uint64
}

func (mk *mockKey) ID() uint64 {
	return mk.id
}

func (mk *mockKey) Compare(other Key) int {
	key := other.(*mockKey)
	if key.value == mk.value {
		return 0
	}
	if key.value > mk.value {
		return 1
	}

	return -1
}

func newMockKey(value, id uint64) *mockKey {
	return &mockKey{value, id}
}
