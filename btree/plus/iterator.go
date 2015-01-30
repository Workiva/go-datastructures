package plus

import "log"

func init() {
	log.Printf(`test`)
}

type payloadIterator struct {
	keys  sortedByIDKeys
	index int
}

func (pi *payloadIterator) next() bool {
	pi.index++
	if pi.index >= len(pi.keys) {
		return false
	}

	return true
}

func (pi *payloadIterator) value() Key {
	return pi.keys[pi.index]
}

type iterator struct {
	node  *lnode
	index int
	pi    *payloadIterator
}

func (iter *iterator) Next() bool {
	if iter.index == -2 {
		return false
	}

	if iter.pi != nil {
		if iter.pi.next() {
			return true
		}
		iter.pi = nil
	}

	iter.index++
	if iter.index >= len(iter.node.keys) {
		iter.node = iter.node.pointer
		if iter.node == nil {
			iter.index = -2
			return false
		}
		iter.index = 0
	}

	iter.pi = &payloadIterator{
		keys:  iter.node.keys[iter.index].(*payload).keys,
		index: -1,
	}
	return iter.pi.next()
}

func (iter *iterator) Value() Key {
	return iter.pi.value()
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
		index: -2,
	}
}
