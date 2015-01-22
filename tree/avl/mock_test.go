package avl

type mockEntry int

func (me mockEntry) Compare(other Entry) int {
	otherMe := other.(mockEntry)
	if me > otherMe {
		return 1
	}

	if me < otherMe {
		return -1
	}

	return 0
}
