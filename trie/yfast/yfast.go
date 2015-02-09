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
Package yfast implements a y-fast trie.  Instead of a red-black BBST
for the leaves, this implementation uses a simple ordered list.  This
package should have searches that are as performant as the x-fast
trie while having faster inserts/deletes and linear space consumption.

Performance characteristics:
Space: O(n)
Get: O(log log M)
Search: O(log log M)
Insert: O(log log M)
Delete: O(log log M)

where n is the number of items in the trie and M is the size of the
universe, ie, 2^m where m is the number of bits in the specified key
size.

This particular implementation also uses fixed bucket sizes as this should
aid in multithreading these functions for performance optimization.
*/
package yfast

import "github.com/Workiva/go-datastructures/trie/xfast"

// YFastTrie implements all the methods available to the y-fast
// trie datastructure.  The top half is composed of an x-fast trie
// while the leaves are composed of ordered lists of type Entry,
// an interface found in this package.
type YFastTrie struct {
	num   uint64
	xfast *xfast.XFastTrie
	bits  uint8
}

func (yfast *YFastTrie) init(intType interface{}) {
	switch intType.(type) {
	case uint8:
		yfast.bits = 8
	case uint16:
		yfast.bits = 16
	case uint32:
		yfast.bits = 32
	case uint, uint64:
		yfast.bits = 64
	default:
		// we'll panic with a bad value to the constructor.
		panic(`Invalid universe size provided.`)
	}

	yfast.xfast = xfast.New(intType)
}

// getBucketKey finds the largest possible value in this key's bucket.
// This is the representative value for the entry in the x-fast trie.
func (yfast *YFastTrie) getBucketKey(key uint64) uint64 {
	i := key/uint64(yfast.bits) + 1
	return uint64(yfast.bits)*i - 1
}

func (yfast *YFastTrie) insert(entry Entry) Entry {
	// first, we need to determine if we have a node in the x-trie
	// that already matches for the key
	bundleKey := yfast.getBucketKey(entry.Key())
	bundle := yfast.xfast.Get(bundleKey)

	if bundle != nil {
		overwritten := bundle.(*entriesWrapper).entries.insert(entry)
		if overwritten == nil {
			yfast.num++
			return nil
		}

		return overwritten
	}

	yfast.num++
	entries := make(Entries, 0, yfast.bits)
	entries.insert(entry)

	ew := &entriesWrapper{
		key:     bundleKey,
		entries: entries,
	}

	yfast.xfast.Insert(ew)
	return nil
}

// Insert will insert the provided entries into the y-fast trie
// and return a list of entries that were overwritten.
func (yfast *YFastTrie) Insert(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, yfast.insert(e))
	}

	return overwritten
}

func (yfast *YFastTrie) delete(key uint64) Entry {
	bundleKey := yfast.getBucketKey(key)

	bundle := yfast.xfast.Get(bundleKey)
	if bundle == nil {
		return nil
	}

	ew := bundle.(*entriesWrapper)
	entry := ew.entries.delete(key)
	if entry == nil {
		return nil
	}

	yfast.num--

	if len(ew.entries) == 0 {
		yfast.xfast.Delete(bundleKey)
	}

	return entry
}

// Delete will delete the provided keys from the y-fast trie
// and return a list of entries that were deleted.
func (yfast *YFastTrie) Delete(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, yfast.delete(key))
	}

	return entries
}

func (yfast *YFastTrie) get(key uint64) Entry {
	bundleKey := yfast.getBucketKey(key)
	bundle := yfast.xfast.Get(bundleKey)
	if bundle == nil {
		return nil
	}

	entry := bundle.(*entriesWrapper).entries.get(key)
	if entry == nil { // go interfaces :(
		return nil
	}

	return entry
}

// Get will look for the provided key in the y-fast trie and return
// the associated value if it is found.  If it is not found, this
// method returns nil.
func (yfast *YFastTrie) Get(key uint64) Entry {
	entry := yfast.get(key)
	if entry == nil {
		return nil
	}

	return entry
}

// Len returns the number of items in the y-fast trie.
func (yfast *YFastTrie) Len() uint64 {
	return yfast.num
}

func (yfast *YFastTrie) successor(key uint64) Entry {
	bundle := yfast.xfast.Successor(key)
	if bundle == nil {
		return nil
	}

	entry, _ := bundle.(*entriesWrapper).entries.successor(key)
	if entry == nil {
		return nil
	}

	return entry
}

// Successor returns an Entry with a key equal to or immediately
// greater than the provided key.  If such an Entry does not exist
// this returns nil.
func (yfast *YFastTrie) Successor(key uint64) Entry {
	entry := yfast.successor(key)
	if entry == nil {
		return nil
	}

	return entry
}

func (yfast *YFastTrie) predecessor(key uint64) Entry {
	// harder case because our representative value in the
	// x-fast trie is the a max
	bundleKey := yfast.getBucketKey(key)
	bundle := yfast.xfast.Predecessor(bundleKey)
	if bundle == nil {
		return nil
	}

	ew := bundle.(*entriesWrapper)
	entry, _ := ew.entries.predecessor(key)
	if entry != nil {
		return entry
	}

	// it's possible we do exist somewhere earlier in the x-fast trie
	bundle = yfast.xfast.Predecessor(bundleKey - 1)
	if bundle == nil {
		return nil
	}

	ew = bundle.(*entriesWrapper)

	entry, _ = ew.entries.predecessor(key)
	if entry == nil {
		return nil
	}

	return entry
}

// Predecessor returns an Entry with a key equal to or immediately
// preceeding than the provided key.  If such an Entry does not exist
// this returns nil.
func (yfast *YFastTrie) Predecessor(key uint64) Entry {
	entry := yfast.predecessor(key)
	if entry == nil {
		return nil
	}

	return entry
}

func (yfast *YFastTrie) iter(key uint64) *Iterator {
	xfastIter := yfast.xfast.Iter(key)
	xfastIter.Next()
	bundle := xfastIter.Value()
	if bundle == nil {
		return nilIterator()
	}

	i := bundle.(*entriesWrapper).entries.search(key)
	return &Iterator{
		index:         i - 1,
		xfastIterator: xfastIter,
		entries:       bundle.(*entriesWrapper),
	}
}

// Iter will return an iterator that will iterate across all values
// that start or immediately proceed the provided key.  Iteration
// happens in ascending order.
func (yfast *YFastTrie) Iter(key uint64) *Iterator {
	return yfast.iter(key)
}

// New constructs, initializes, and returns a new y-fast trie.
// Provided should be a uint type that specifies the number
// of bits in the desired universe.  This will affect the time
// complexity of all lookup and mutate operations.
func New(ifc interface{}) *YFastTrie {
	yfast := &YFastTrie{}
	yfast.init(ifc)
	return yfast
}
