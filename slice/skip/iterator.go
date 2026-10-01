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

import "github.com/Workiva/go-datastructures/common"

// iterator represents an object that can be iterated.  It will
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

// Value returns a Comparator representing the iterator's present
// position in the query.  Returns nil if no values remain to iterate.
func (iter *iterator) Value() common.Comparator {
	if iter.n == nil {
		return nil
	}

	return iter.n.entry
}

// exhaust is a helper method to exhaust this iterator and return
// all remaining entries.
func (iter *iterator) exhaust() common.Comparators {
	entries := make(common.Comparators, 0, 10)
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
