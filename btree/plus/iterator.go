package plus

const iteratorExhausted = -2

type iterator struct {
	node  *lnode
	index int
}

func (iter *iterator) Next() bool {
	if iter.index == iteratorExhausted {
		return false
	}

	iter.index++
	if iter.index >= len(iter.node.keys) {
		iter.node = iter.node.pointer
		if iter.node == nil {
			iter.index = iteratorExhausted
			return false
		}
		iter.index = 0
	}

	return true
}

func (iter *iterator) Value() Key {
	if iter.index == iteratorExhausted ||
		iter.index < 0 || iter.index >= len(iter.node.keys) {

		return nil
	}

	return iter.node.keys[iter.index]
}

// exhaust is a test function that's not exported
func (iter *iterator) exhaust() keys {
	keys := make(keys, 0, 10)
	for iter := iter; iter.Next(); {
		keys = append(keys, iter.Value())
	}

	return keys
}

func nilIterator() *iterator {
	return &iterator{
		index: iteratorExhausted,
	}
}
