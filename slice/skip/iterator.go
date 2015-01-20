/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package skip

const iteratorExhausted = -2

// Iterator represents an object that can be iterated.  It will
// return false on Next and nil on Value if there are no further
// values to be iterated.
type iterator struct {
	first bool
	n     *node
}

// Next returns a bool indicating if there are any further values
// in this iterator.
func (iter *iterator) Next() bool {
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
func (iter *iterator) Value() Entry {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
}

// exhaust is a helper method to exhaust this iterator and return
// all remaining entries.
func (iter *iterator) exhaust() Entries {
	entries := make(Entries, 0, 10)
	for i := iter; i.Next(); {
		entries = append(entries, i.Value())
	}

	return entries
}

// nilIterator returns an iterator that will always return false
// for Next and nil for Value.
func nilIterator() *iterator {
	return &iterator{}
}

type starIterator struct {
	entries Entries
	iter    Iterator
	index   int
}

func (si *starIterator) isExhausted() bool {
	return si.index == iteratorExhausted
}

func (si *starIterator) Next() bool {
	if si.isExhausted() {
		return false
	}

	si.index++
	if si.index >= len(si.entries) {
		canNext := si.iter.Next()
		if !canNext {
			si.index = iteratorExhausted
			return false
		}

		si.entries = si.iter.Value().(*entryBundle).entries
		si.index = 0
	}

	return true
}

func (si *starIterator) Value() Entry {
	if si.isExhausted() {
		return nil
	}

	if si.entries == nil || si.index < 0 || si.index >= len(si.entries) {
		return nil
	}

	return si.entries[si.index]
}

func (si *starIterator) exhaust() Entries {
	entries := make(Entries, 0, 20)
	for i := si; i.Next(); {
		entries = append(entries, i.Value())
	}

	return entries
}
