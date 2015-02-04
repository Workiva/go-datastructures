package palm

type mockKey int

func (mk mockKey) Compare(other Key) int {
	otherKey := other.(mockKey)

	if mk > otherKey {
		return 1
	}

	if mk < otherKey {
		return -1
	}

	return 0
}
