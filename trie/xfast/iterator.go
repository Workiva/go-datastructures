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

/*
Package err implements a threadsafe error interface.  In my places,
I found myself needing a lock to protect writing to a common error interface
from multiple go routines (channels are great but slow).  This just makes
that process more convenient.
*/

package xfast

// Entries is a typed list of Entry interfaces.
type Entries []Entry

// Iterator will iterate of the results of a query.
type Iterator struct {
	n     *node
	first bool
}

// Next will return a bool indicating if another value exists
// in the iterator.
func (iter *Iterator) Next() bool {
	if iter.first {
		iter.first = false
		return iter.n != nil
	}

	iter.n = iter.n.children[1]
	return iter.n != nil
}

// Value will return the Entry representing the iterator's current position.
// If no Entry exists at the present condition, the iterator is
// exhausted and this method will return nil.
func (iter *Iterator) Value() Entry {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
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
