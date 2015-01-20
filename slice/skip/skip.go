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
Package skip defines a skiplist datastructure.  That is, a data structure
that probabilistically determines relationships between keys.  By doing
so, it becomes easier to program than a binary search tree but maintains
similar speeds.

Performance characteristics:
Insert: O(log n)
Search: O(log n)
Delete: O(log n)
Space: O(n)

More information here: http://cglab.ca/~morin/teaching/5408/refs/p90b.pdf
*/

package skip

import (
	"math/rand"
	"time"
)

const p = .5 // the p level defines the probability that a node
// with a value at level i also has a value at i+1.  This number
// is also important in determining max level.  Max level will
// be defined as L(N) where L = log base (1/p) of n where n
// is the number of items in the list and N is the number of possible
// items in the universe.  If p = .5 then maxlevel = 32 is appropriate
// for uint32.

// generator will be the common generator to create random numbers. It
// is seeded with unix nanosecond when this line is executed at runtime,
// and only executed once ensuring all random numbers come from the same
// randomly seeded generator.
var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateLevel(maxLevel uint8) uint8 {
	var level uint8
	for level = uint8(1); level < maxLevel-1; level++ {
		if generator.Float64() >= p {
			return level
		}
	}

	return level
}

// Skip list is a datastructure that probabalistically determines
// relationships between nodes.  This results in a structure
// that performs similarly to a BST but is much easier to build
// from a programmatic perspective (no rotations).
type SkipList struct {
	maxLevel, level uint8
	head            *node
	num             uint64
	// a list of nodes that can be reused, should reduce
	// the number of allocations in the insert/delete case.
	cache nodes
}

// init will initialize this skiplist.  The parameter is expected
// to be of some uint type which will set this skiplist's maximum
// level.
func (sl *SkipList) init(ifc interface{}) {
	switch ifc.(type) {
	case uint8:
		sl.maxLevel = 8
	case uint16:
		sl.maxLevel = 16
	case uint32:
		sl.maxLevel = 32
	case uint64, uint:
		sl.maxLevel = 64
	}
	sl.cache = make(nodes, sl.maxLevel)
	sl.head = newNode(nil, sl.maxLevel)
}

func (sl *SkipList) search(key uint64, update []*node) *node {
	if sl.num == 0 { // nothing in the list
		return nil
	}

	var offset uint8
	n := sl.head
	for i := uint8(0); i <= sl.level; i++ {
		offset = sl.level - i
		for n.forward[offset] != nil && n.forward[offset].key() < key {
			n = n.forward[offset]
		}

		if update != nil {
			update[offset] = n
		}
	}

	return n.forward[0]
}

// Get will retrieve values associated with the keys provided.  If an
// associated value could not be found, a nil is returned in its place.
// This is an O(log n) operation.
func (sl *SkipList) Get(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))

	var n *node
	for _, key := range keys {
		n = sl.search(key, nil)
		if n == nil {
			entries = append(entries, nil)
		} else {
			entries = append(entries, n.entry)
		}
	}

	return entries
}

func (sl *SkipList) insert(entry Entry) Entry {
	sl.cache.reset()
	n := sl.search(entry.Key(), sl.cache)
	if n != nil && n.key() == entry.Key() { // a simple update in this case
		oldEntry := n.entry
		n.entry = entry
		return oldEntry
	}
	sl.num++

	nodeLevel := generateLevel(sl.maxLevel)
	if nodeLevel > sl.level {
		for i := sl.level; i <= nodeLevel; i++ {
			sl.cache[i] = sl.head
		}
		sl.level = nodeLevel
	}

	nn := newNode(entry, sl.maxLevel)
	for i := uint8(0); i <= nodeLevel; i++ {
		nn.forward[i] = sl.cache[i].forward[i]
		sl.cache[i].forward[i] = nn
	}

	return nil
}

// Insert will insert the provided entries into the list.  Returned
// is a list of entries that were overwritten.  This is expected to
// be an O(log n) operation.
func (sl *SkipList) Insert(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, sl.insert(e))
	}

	return overwritten
}

func (sl *SkipList) delete(key uint64) Entry {
	sl.cache.reset()
	n := sl.search(key, sl.cache)

	if n == nil || n.entry.Key() != key {
		return nil
	}

	sl.num--

	for i := uint8(0); i <= sl.level; i++ {
		if sl.cache[i].forward[i] != n {
			break
		}

		sl.cache[i].forward[i] = n.forward[i]
	}

	for sl.level > 0 && sl.head.forward[sl.level] == nil {
		sl.level = sl.level - 1
	}

	return n.entry
}

// Delete will remove the provided keys from the skiplist and return
// a list of in-order entries that were deleted.  This is a no-op if
// an associated key could not be found.  This is an O(log n) operation.
func (sl *SkipList) Delete(keys ...uint64) Entries {
	deleted := make(Entries, 0, len(keys))

	for _, key := range keys {
		deleted = append(deleted, sl.delete(key))
	}

	return deleted
}

// Len returns the number of items in this skiplist.
func (sl *SkipList) Len() uint64 {
	return sl.num
}

func (sl *SkipList) iter(key uint64) *iterator {
	n := sl.search(key, nil)
	if n == nil {
		return nilIterator()
	}

	return &iterator{
		first: true,
		n:     n,
	}
}

// Iter will return an iterator that can be used to iterate
// over all the values with a key equal to or greater than
// the key provided.
func (sl *SkipList) Iter(key uint64) Iterator {
	return sl.iter(key)
}

// New will allocate, initialize, and return a new skiplist.
// The provided parameter should be of type uint and will determine
// the maximum possible level that will be created to ensure
// a random and quick distribution of levels.  Parameter must
// be a uint type.
func New(ifc interface{}) *SkipList {
	sl := &SkipList{}
	sl.init(ifc)
	return sl
}
