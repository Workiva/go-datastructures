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

package yfast

import "sort"

type entriesWrapper struct {
	key     uint64
	entries Entries
}

// Key will return the key of the highest entry in this list.
// This is required by the x-fast trie Entry interface.  This
// returns 0 if this list is empty.
func (ew *entriesWrapper) Key() uint64 {
	return ew.key
}

// Entries is a typed list of Entry.  The size of entries
// will be limited to 1/2log M to 2log M where M is the size
// of the universe.
type Entries []Entry

// search will perform a sort package search on this list
// of entries and return an index indicating position.
// If the returned index is >= len(entries) then a suitable
// position could not be found.  The index does not guarantee
// equality, just indicates where the key would be inserted.
func (entries Entries) search(key uint64) int {
	return sort.Search(len(entries), func(i int) bool {
		return entries[i].Key() >= key
	})
}

// insert will insert the provided entry into this list of
// entries.  Returned is an entry if an entry already exists
// for the provided key.  If nothing is overwritten, Entry
// will be nil.
func (entries *Entries) insert(entry Entry) Entry {
	i := entries.search(entry.Key())

	if i == len(*entries) {
		*entries = append(*entries, entry)
		return nil
	}

	if (*entries)[i].Key() == entry.Key() {
		oldEntry := (*entries)[i]
		(*entries)[i] = entry
		return oldEntry
	}

	(*entries) = append(*entries, nil)
	copy((*entries)[i+1:], (*entries)[i:])
	(*entries)[i] = entry
	return nil
}

// delete will remove the provided key from this list of entries.
// Returned is a deleted Entry.  This will be nil if the key
// cannot be found.
func (entries *Entries) delete(key uint64) Entry {
	i := entries.search(key)
	if i == len(*entries) { // key not found
		return nil
	}

	if (*entries)[i].Key() != key {
		return nil
	}

	oldEntry := (*entries)[i]
	copy((*entries)[i:], (*entries)[i+1:])
	(*entries)[len(*entries)-1] = nil // GC
	*entries = (*entries)[:len(*entries)-1]
	return oldEntry
}

// max returns the value of the highest key in this list
// of entries.  The bool indicates if it's a valid key, that
// is if there is more than zero entries in this list.
func (entries Entries) max() (uint64, bool) {
	if len(entries) == 0 {
		return 0, false
	}

	return entries[len(entries)-1].Key(), true
}

// get will perform a lookup over this list of entries
// and return an Entry if it exists.  Returns nil if the
// entry does not exist.
func (entries Entries) get(key uint64) Entry {
	i := entries.search(key)
	if i == len(entries) {
		return nil
	}

	if entries[i].Key() == key {
		return entries[i]
	}

	return nil
}

// successor will return the first entry that has a key
// greater than or equal to provided key.  Also returned
// is the index of the find.  Returns nil, -1 if a successor does
// not exist.
func (entries Entries) successor(key uint64) (Entry, int) {
	i := entries.search(key)
	if i == len(entries) {
		return nil, -1
	}

	return entries[i], i
}

// predecessor will return the first entry that has a key
// less than or equal to the provided key.  Also returned
// is the index of the find.  Returns nil, -1 if a predecessor
// does not exist.
func (entries Entries) predecessor(key uint64) (Entry, int) {
	if len(entries) == 0 {
		return nil, -1
	}

	i := entries.search(key)
	if i == len(entries) {
		return entries[i-1], i - 1
	}

	if entries[i].Key() == key {
		return entries[i], i
	}

	i--

	if i < 0 {
		return nil, -1
	}

	return entries[i], i
}
