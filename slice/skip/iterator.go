package skip

// Iterator represents an object that can be iterated.  It will
// return false on Next and nil on Value if there are no further
// values to be iterated.
type Iterator struct {
	first bool
	n     *node
}

// Next returns a bool indicating if there are any further values
// in this iterator.
func (iter *Iterator) Next() bool {
	if iter.first {
		iter.first = false
		return iter.n != nil
	}

	if iter.n == nil {
		return false
	}

	iter.n = iter.n.forward[0]
	return iter.n != nil
}

// Value returns an Entry representing the iterator's present
// position in the query.  Returns nil if no values remain to iterate.
func (iter *Iterator) Value() Entry {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
}

// exhaust is a helper method to exhaust this iterator and return
// all remaining entries.
func (iter *Iterator) exhaust() Entries {
	entries := make(Entries, 0, 10)
	for i := iter; i.Next(); {
		entries = append(entries, i.Value())
	}

	return entries
}

// nilIterator returns an iterator that will always return false
// for Next and nil for Value.
func nilIterator() *Iterator {
	return &Iterator{}
}
