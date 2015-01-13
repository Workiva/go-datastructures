package xfast

type Entries []Entry

type Iterator struct {
	n     *node
	first bool
}

func (iter *Iterator) Next() bool {
	if iter.first {
		iter.first = false
		return iter.n != nil
	}

	iter.n = iter.n.children[1]
	return iter.n != nil
}

func (iter *Iterator) Value() Entry {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
}

func (iter *Iterator) exhaust() Entries {
	entries := make(Entries, 0, 100)
	for it := iter; it.Next(); {
		entries = append(entries, it.Value())
	}

	return entries
}
