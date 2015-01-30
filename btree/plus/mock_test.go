package plus

type mockKey struct {
	value int
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

func newMockKey(value int) *mockKey {
	return &mockKey{value}
}
