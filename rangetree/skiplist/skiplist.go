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
Package skiplist implements an n-dimensional rangetree based on a skip
list.  This should be faster than a straight slice implementation as
memcopy is avoided.

Time complexities revolve around the ability to quickly find items
in the n-dimensional skiplist.  That time can be defined by the
number of items in any dimension.  Let N1, N2,... Nn define the
number of dimensions.

Performance characteristics:
Space: O(n)
Search: O(log N1 + log N2 + ...log Nn) = O(log N1*N2*...Nn)
Insert: O(log N1 + log N2 + ...log Nn) = O(log N1*N2*...Nn)
Delete: O(log N1 + log N2 + ...log Nn) = O(log N1*N2*...Nn)
*/

package skiplist

import (
	"github.com/Workiva/go-datastructures/rangetree"
	"github.com/Workiva/go-datastructures/slice/skip"
)

// keyed is required as in the rangetree code we often want to compare
// two different types of bundles and this allows us to do so without
// checking for each one.
type keyed interface {
	key() uint64
}

type skipEntry uint64

// Compare is required by the skip.Entry interface.
func (se skipEntry) Compare(other skip.Entry) int {
	otherSe := other.(skipEntry)
	if se == otherSe {
		return 0
	}

	if se > otherSe {
		return 1
	}

	return -1
}

func (se skipEntry) key() uint64 {
	return uint64(se)
}

// isLastDimension simply returns dimension == lastDimension-1.
// This panics if dimension >= lastDimension.
func isLastDimension(dimension, lastDimension uint64) bool {
	if dimension >= lastDimension { // useful in testing and denotes a serious problem
		panic(`Dimension is greater than possible dimensions.`)
	}

	return dimension == lastDimension-1
}

// needsDeletion returns a bool indicating if the provided value
// needs to be deleted based on the provided index and number.
func needsDeletion(value, index, number int64) bool {
	if number > 0 {
		return false
	}

	number = -number // get the magnitude
	offset := value - index

	return offset >= 0 && offset < number
}

// dimensionalBundle is an intermediate holder up to the last
// dimension and represents a wrapper around a skiplist.
type dimensionalBundle struct {
	id uint64
	sl *skip.SkipList
}

// Compare returns a value indicating the relative relationship and the
// provided bundle.
func (db *dimensionalBundle) Compare(e skip.Entry) int {
	keyed := e.(keyed)
	if db.id == keyed.key() {
		return 0
	}

	if db.id > keyed.key() {
		return 1
	}

	return -1
}

// key returns the key for this bundle.
func (db *dimensionalBundle) key() uint64 {
	return db.id
}

// lastBundle represents a bundle living at the last dimension
// of the tree.
type lastBundle struct {
	id    uint64
	entry rangetree.Entry
}

// Compare returns a value indicating the relative relationship and the
// provided bundle.
func (lb *lastBundle) Compare(e skip.Entry) int {
	keyed := e.(keyed)
	if lb.id == keyed.key() {
		return 0
	}

	if lb.id > keyed.key() {
		return 1
	}

	return -1
}

// Key returns the key for this bundle.
func (lb *lastBundle) key() uint64 {
	return lb.id
}

type skipListRT struct {
	top                *skip.SkipList
	dimensions, number uint64
}

func (rt *skipListRT) init(dimensions uint64) {
	rt.dimensions = dimensions
	rt.top = skip.New(uint64(0))
}

func (rt *skipListRT) add(entry rangetree.Entry) rangetree.Entry {
	var (
		value int64
		e     skip.Entry
		sl    = rt.top
		db    *dimensionalBundle
		lb    *lastBundle
	)

	for i := uint64(0); i < rt.dimensions; i++ {
		value = entry.ValueAtDimension(i)
		e = sl.Get(skipEntry(value))[0]
		if isLastDimension(i, rt.dimensions) {
			if e != nil { // this is an overwrite
				lb = e.(*lastBundle)
				oldEntry := lb.entry
				lb.entry = entry
				return oldEntry
			}

			// need to add new sl entry
			lb = &lastBundle{id: uint64(value), entry: entry}
			rt.number++
			sl.Insert(lb)
			return nil
		}

		if e == nil { // we need the intermediate dimension
			db = &dimensionalBundle{id: uint64(value), sl: skip.New(uint64(0))}
			sl.Insert(db)
		} else {
			db = e.(*dimensionalBundle)
		}

		sl = db.sl
	}

	panic(`Ran out of dimensions before for loop completed.`)
}

// Add will add the provided entries to the tree.  Any entries that
// were overwritten will be returned in the order in which they
// were overwritten.  If an entry's addition does not overwrite, a nil
// is returned for that cell for its index in the provided entries.
func (rt *skipListRT) Add(entries ...rangetree.Entry) rangetree.Entries {
	overwritten := make(rangetree.Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, rt.add(e))
	}

	return overwritten
}

func (rt *skipListRT) get(entry rangetree.Entry) rangetree.Entry {
	var (
		sl    = rt.top
		e     skip.Entry
		value uint64
	)
	for i := uint64(0); i < rt.dimensions; i++ {
		value = uint64(entry.ValueAtDimension(i))
		e = sl.Get(skipEntry(value))[0]
		if e == nil {
			return nil
		}

		if isLastDimension(i, rt.dimensions) {
			return e.(*lastBundle).entry
		}

		sl = e.(*dimensionalBundle).sl
	}

	panic(`Reached past for loop without finding last dimension.`)
}

// Get will return any rangetree.Entries matching the provided entries.
// Similar in functionality to a key lookup, this returns nil for any
// entry that could not be found.
func (rt *skipListRT) Get(entries ...rangetree.Entry) rangetree.Entries {
	results := make(rangetree.Entries, 0, len(entries))
	for _, e := range entries {
		results = append(results, rt.get(e))
	}

	return results
}

// Len returns the number of entries in the tree.
func (rt *skipListRT) Len() uint64 {
	return rt.number
}

// deleteRecursive is used by the delete logic.  The recursion depth
// only goes as far as the number of dimensions, so this shouldn't be an
// issue.
func (rt *skipListRT) deleteRecursive(sl *skip.SkipList, dimension uint64,
	entry rangetree.Entry) rangetree.Entry {

	value := entry.ValueAtDimension(dimension)
	if isLastDimension(dimension, rt.dimensions) {
		entries := sl.Delete(skipEntry(value))
		if entries[0] == nil {
			return nil
		}

		rt.number--
		return entries[0].(*lastBundle).entry
	}

	db, ok := sl.Get(skipEntry(value))[0].(*dimensionalBundle)
	if !ok { // value was not found
		return nil
	}

	result := rt.deleteRecursive(db.sl, dimension+1, entry)
	if result == nil {
		return nil
	}

	if db.sl.Len() == 0 {
		sl.Delete(db)
	}

	return result
}

func (rt *skipListRT) delete(entry rangetree.Entry) rangetree.Entry {
	return rt.deleteRecursive(rt.top, 0, entry)
}

// Delete will remove the provided entries from the tree.
func (rt *skipListRT) Delete(entries ...rangetree.Entry) {
	for _, e := range entries {
		rt.delete(e)
	}
}

func (rt *skipListRT) apply(sl *skip.SkipList, dimension uint64,
	interval rangetree.Interval, fn func(rangetree.Entry) bool) bool {

	lowValue, highValue := interval.LowAtDimension(dimension), interval.HighAtDimension(dimension)

	var e skip.Entry

	for iter := sl.Iter(skipEntry(lowValue)); iter.Next(); {
		e = iter.Value()
		if int64(e.(keyed).key()) >= highValue {
			break
		}

		if isLastDimension(dimension, rt.dimensions) {
			if !fn(e.(*lastBundle).entry) {
				return false
			}
		} else {

			if !rt.apply(e.(*dimensionalBundle).sl, dimension+1, interval, fn) {
				return false
			}
		}
	}

	return true
}

// Apply will call the provided function with each entry that exists
// within the provided range, in order.  Return false at any time to
// cancel iteration.  Altering the entry in such a way that its location
// changes will result in undefined behavior.
func (rt *skipListRT) Apply(interval rangetree.Interval, fn func(rangetree.Entry) bool) {
	rt.apply(rt.top, 0, interval, fn)
}

// Query will return a list of entries that fall within
// the provided interval.
func (rt *skipListRT) Query(interval rangetree.Interval) rangetree.Entries {
	entries := make(rangetree.Entries, 0, 100)
	rt.apply(rt.top, 0, interval, func(e rangetree.Entry) bool {
		entries = append(entries, e)
		return true
	})

	return entries
}

func (rt *skipListRT) flatten(sl *skip.SkipList, dimension uint64, entries *rangetree.Entries) {
	lastDimension := isLastDimension(dimension, rt.dimensions)
	for iter := sl.Iter(skipEntry(0)); iter.Next(); {
		if lastDimension {
			*entries = append(*entries, iter.Value().(*lastBundle).entry)
		} else {
			rt.flatten(iter.Value().(*dimensionalBundle).sl, dimension+1, entries)
		}
	}
}

func (rt *skipListRT) insert(sl *skip.SkipList, dimension, insertDimension uint64,
	index, number int64, deleted, affected *rangetree.Entries) {

	var e skip.Entry
	lastDimension := isLastDimension(dimension, rt.dimensions)
	affectedDimension := dimension == insertDimension
	var iter skip.Iterator
	if dimension == insertDimension {
		iter = sl.Iter(skipEntry(index))
	} else {
		iter = sl.Iter(skipEntry(0))
	}

	var toDelete skip.Entries
	if number < 0 {
		toDelete = make(skip.Entries, 0, 100)
	}

	for iter.Next() {
		e = iter.Value()
		if !affectedDimension {
			rt.insert(e.(*dimensionalBundle).sl, dimension+1,
				insertDimension, index, number, deleted, affected,
			)
			continue
		}
		if needsDeletion(int64(e.(keyed).key()), index, number) {
			toDelete = append(toDelete, e)
			continue
		}

		if lastDimension {
			e.(*lastBundle).id += uint64(number)
			*affected = append(*affected, e.(*lastBundle).entry)
		} else {
			e.(*dimensionalBundle).id += uint64(number)
			rt.flatten(e.(*dimensionalBundle).sl, dimension+1, affected)
		}
	}

	if len(toDelete) > 0 {
		for _, e := range toDelete {
			if lastDimension {
				*deleted = append(*deleted, e.(*lastBundle).entry)
			} else {
				rt.flatten(e.(*dimensionalBundle).sl, dimension+1, deleted)
			}
		}

		sl.Delete(toDelete...)
	}
}

// InsertAtDimension will increment items at and above the given index
// by the number provided.  Provide a negative number to to decrement.
// Returned are two lists.  The first list is a list of entries that
// were moved.  The second is a list entries that were deleted.  These
// lists are exclusive.
func (rt *skipListRT) InsertAtDimension(dimension uint64,
	index, number int64) (rangetree.Entries, rangetree.Entries) {

	if dimension >= rt.dimensions || number == 0 {
		return rangetree.Entries{}, rangetree.Entries{}
	}

	affected := make(rangetree.Entries, 0, 100)
	var deleted rangetree.Entries
	if number < 0 {
		deleted = make(rangetree.Entries, 0, 100)
	}

	rt.insert(rt.top, 0, dimension, index, number, &deleted, &affected)
	rt.number -= uint64(len(deleted))
	return affected, deleted
}

func new(dimensions uint64) *skipListRT {
	sl := &skipListRT{}
	sl.init(dimensions)
	return sl
}

// New will allocate, initialize, and return a new rangetree.RangeTree
// with the provided number of dimensions.
func New(dimensions uint64) rangetree.RangeTree {
	return new(dimensions)
}
