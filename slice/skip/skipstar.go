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
SkipList* is a data structure combining elements of both a skiplist
and the bottom half of a Y-fast trie.  The idea is that you can quickly
search for a sub branch of the skip list and that this branch can
fit entirely within cache, thereby improving the performance characteristics
over a standard skip list.  This also keeps any memcopy operation limited
to O(log M) where M is the size of the desired universe.

Performance vs standard skip list.
BenchmarkInsert-8	 2000000	       949 ns/op
BenchmarkGet-8	 3000000	       516 ns/op
BenchmarkDelete-8	 3000000	       499 ns/op
BenchmarkStarInsert-8	 3000000	       453 ns/op
BenchmarkStarGet-8	 3000000	       524 ns/op
BenchmarkStarDelete-8	 3000000	       469 ns/op
*/
package skip

import "log"

func init() {
	log.Printf(`I HATE THIS.`)
}

// SkipList* implements all methods of a standard skip list but attempts
// to improve performance by ensuring cache locality.
type SkipStarList struct {
	ary uint8
	num uint64
	sl  *SkipList
}

// entryBundle is a helper struct used to define the nodes that
// can be inserted into the SkipList*.
type entryBundle struct {
	key     uint64
	entries Entries
}

// Key will return the key associated with this entity bundle.
// This is required by the Entry interface.
func (eb *entryBundle) Key() uint64 {
	return eb.key
}

func newEntryBundle(key uint64, size uint8) *entryBundle {
	return &entryBundle{
		key:     key,
		entries: make(Entries, 0, size),
	}
}

func (ssl *SkipStarList) init(ifc interface{}) {
	switch ifc.(type) {
	case uint8:
		ssl.ary = 8
	case uint16:
		ssl.ary = 16
	case uint32:
		ssl.ary = 32
	case uint64, uint:
		ssl.ary = 64
	}
	ssl.sl = New(ifc)
}

func (ssl *SkipStarList) getNormalizedKey(key uint64) uint64 {
	key = key/uint64(ssl.ary) + 1
	return key * uint64(ssl.ary)
}

func (ssl *SkipStarList) insert(entry Entry) Entry {
	key := ssl.getNormalizedKey(entry.Key())
	eb, ok := ssl.sl.Get(entry.Key())[0].(*entryBundle)
	if !ok {
		eb = newEntryBundle(key, ssl.ary)
		ssl.sl.Insert(eb)
	}

	e := eb.entries.insert(entry)
	if e == nil {
		ssl.num++
	}
	return e
}

// Insert will insert the provded entries into the SkipList*.  Any
// existing entry with a matching key will be overwritten.  The returned
// list of a list of entries that were overwritten, in order.  A nil
// will be in the in-order position for any non-overwritten entries.
func (ssl *SkipStarList) Insert(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, ssl.insert(e))
	}

	return overwritten
}

func (ssl *SkipStarList) get(key uint64) Entry {
	normalizedKey := ssl.getNormalizedKey(key)
	eb, ok := ssl.sl.Get(normalizedKey)[0].(*entryBundle)
	if ok {
		return eb.entries.get(key)
	}
	return nil
}

// Get will return a list of entries associated with the provided keys.
// A nil will be returned for any key not found.
func (ssl *SkipStarList) Get(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, ssl.get(key))
	}

	return entries
}

func (ssl *SkipStarList) delete(key uint64) Entry {
	normalizedKey := ssl.getNormalizedKey(key)
	eb, ok := ssl.sl.Get(normalizedKey)[0].(*entryBundle)
	if !ok {
		return nil
	}

	deleted := eb.entries.delete(key)
	if deleted != nil {
		ssl.num--
		if len(eb.entries) == 0 {
			ssl.sl.Delete(eb.key)
		}
	}

	return deleted
}

// Delete will remove the provided keys from the SkipList* and
// return a list of entries that were deleted.
func (ssl *SkipStarList) Delete(keys ...uint64) Entries {
	deleted := make(Entries, 0, len(keys))
	for _, key := range keys {
		deleted = append(deleted, ssl.delete(key))
	}

	return deleted
}

func (ssl *SkipStarList) iter(key uint64) *starIterator {
	normalizedKey := ssl.getNormalizedKey(key)
	iter := ssl.sl.Iter(normalizedKey)
	if !iter.Next() {
		return &starIterator{
			index: iteratorExhausted,
		}
	}

	eb := iter.Value().(*entryBundle)
	return &starIterator{
		index:   eb.entries.search(key) - 1,
		entries: eb.entries,
		iter:    iter,
	}
}

// Iter will return an iterator that will visit every value
// equal to or greater than the provided key.
func (ssl *SkipStarList) Iter(key uint64) Iterator {
	return ssl.iter(key)
}

// Len returns the number of items in the SkipList*.
func (ssl *SkipStarList) Len() uint64 {
	return ssl.num
}

// NewStar will allocate, initialize, and return a new SkipStarList.
// The Skip* list has an node size defined by the provided interface
// parameter.  This parameter must be a uint type (uint8, uint16, etc).
func NewStar(ifc interface{}) *SkipStarList {
	ssl := &SkipStarList{}
	ssl.init(ifc)
	return ssl
}
