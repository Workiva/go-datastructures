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

// Entry defines items that can be inserted into the skip list.
// This will also be the type returned from a query.
type Entry interface {
	// Key defines this entry's place in the skip list.
	Key() uint64
}

// Entries is a typed list of interface Entry.
type Entries []Entry

// Iterator defines an interface that allows a consumer to iterate
// all results of a query.  All values will be visited in-order.
type Iterator interface {
	// Next returns a bool indicating if there is future value
	// in the iterator and moves the iterator to that value.
	Next() bool
	// Value returns an Entry representing the iterator's current
	// position.  If there is no value, this returns nil.
	Value() Entry
	// exhaust is a helper method that will iterate this iterator
	// to completion and return a list of resulting Entries
	// in order.
	exhaust() Entries
}
