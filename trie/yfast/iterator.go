package yfast

import "github.com/Workiva/go-datastructures/trie/xfast"

const iteratorExhausted = -2

func iterExhausted(iter *Iterator) bool {
	return iter.index == iteratorExhausted
}

type iterator interface {
	Next() bool
	Value() xfast.Entry
}

// Iterator will iterate of the results of a query.
type Iterator struct {
	xfastIterator iterator
	index         int
	entries       *entriesWrapper
}

// Next will return a bool indicating if another value exists
// in the iterator.
func (iter *Iterator) Next() bool {
	if iterExhausted(iter) {
		return false
	}
	iter.index++
	if iter.index >= len(iter.entries.entries) {
		next := iter.xfastIterator.Next()
		if !next {
			iter.index = iteratorExhausted
			return false
		}
		var ok bool
		iter.entries, ok = iter.xfastIterator.Value().(*entriesWrapper)
		if !ok {
			iter.index = iteratorExhausted
			return false
		}
		iter.index = 0
	}

	return true
}

// Value will return the Entry representing the iterator's current position.
// If no Entry exists at the present condition, the iterator is
// exhausted and this method will return nil.
func (iter *Iterator) Value() Entry {
	if iterExhausted(iter) {
		return nil
	}

	if iter.entries == nil || iter.index < 0 || iter.index >= len(iter.entries.entries) {
		return nil
	}

	return iter.entries.entries[iter.index]
}

// exhaust is a helper function that will exhaust this iterator
// and return a list of entries.  This is for internal use only.
func (iter *Iterator) exhaust() Entries {
	entries := make(Entries, 0, 100)
	for it := iter; it.Next(); {
		entries = append(entries, it.Value())
	}

	return entries
}

func nilIterator() *Iterator {
	return &Iterator{
		index: iteratorExhausted,
	}
}
