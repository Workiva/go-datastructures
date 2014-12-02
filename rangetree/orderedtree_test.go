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

package rangetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func constructMultiDimensionalOrderedTree(number uint64) (
	*orderedTree, Entries) {

	tree := newOrderedTree(2)

	entries := make(Entries, 0, number)
	for i := uint64(0); i < number; i++ {
		entries = append(entries, constructMockEntry(i, int64(i), int64(i)))
	}

	tree.Add(entries...)

	return tree, entries
}

func TestOTRootAddMultipleDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(1)

	assert.Equal(t, uint64(1), tree.Len())

	result := tree.Query(constructMockInterval(dimension{0, 1}, dimension{0, 1}))
	assert.Equal(t, Entries{entries[0]}, result)
}

func TestOTMultipleAddMultipleDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(4)

	assert.Equal(t, uint64(4), tree.Len())

	result := tree.Query(constructMockInterval(dimension{0, 1}, dimension{0, 1}))
	assert.Equal(t, Entries{entries[0]}, result)

	result = tree.Query(constructMockInterval(dimension{3, 4}, dimension{3, 4}))
	assert.Equal(t, Entries{entries[3]}, result)

	result = tree.Query(constructMockInterval(dimension{0, 4}, dimension{0, 4}))
	assert.Equal(t, entries, result)

	result = tree.Query(constructMockInterval(dimension{1, 3}, dimension{1, 3}))
	assert.Equal(t, Entries{entries[1], entries[2]}, result)

	result = tree.Query(constructMockInterval(dimension{0, 2}, dimension{10, 20}))
	assert.Len(t, result, 0)

	result = tree.Query(constructMockInterval(dimension{10, 20}, dimension{0, 2}))
	assert.Len(t, result, 0)

	result = tree.Query(constructMockInterval(dimension{0, 2}, dimension{0, 1}))
	assert.Equal(t, Entries{entries[0]}, result)

	result = tree.Query(constructMockInterval(dimension{0, 1}, dimension{0, 2}))
	assert.Equal(t, Entries{entries[0]}, result)
}

func TestOTAddInOrderMultiDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(10)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Equal(t, uint64(10), tree.Len())
	assert.Len(t, result, 10)
	assert.Equal(t, entries, result)
}

func TestOTAddReverseOrderMultiDimensions(t *testing.T) {
	tree := newOrderedTree(2)

	for i := uint64(10); i > 0; i-- {
		tree.Add(constructMockEntry(i, int64(i), int64(i)))
	}

	result := tree.Query(constructMockInterval(dimension{0, 11}, dimension{0, 11}))
	assert.Len(t, result, 10)
	assert.Equal(t, uint64(10), tree.Len())
}

func TestOTAddRandomOrderMultiDimensions(t *testing.T) {
	tree := newOrderedTree(2)

	starts := []uint64{0, 4, 2, 1, 3}

	for _, start := range starts {
		tree.Add(constructMockEntry(start, int64(start), int64(start)))
	}

	result := tree.Query(constructMockInterval(dimension{0, 5}, dimension{0, 5}))
	assert.Len(t, result, 5)
	assert.Equal(t, uint64(5), tree.Len())
}

func TestOTAddLargeNumbersMultiDimension(t *testing.T) {
	numItems := uint64(1000)
	tree := newOrderedTree(2)

	for i := uint64(0); i < numItems; i++ {
		tree.Add(constructMockEntry(i, int64(i), int64(i)))
	}

	result := tree.Query(
		constructMockInterval(
			dimension{0, int64(numItems)},
			dimension{0, int64(numItems)},
		),
	)
	assert.Equal(t, numItems, tree.Len())
	assert.Len(t, result, int(numItems))
}

func BenchmarkOTAddItemsMultiDimensions(b *testing.B) {
	numItems := uint64(1000)
	entries := make(Entries, 0, numItems)

	for i := uint64(0); i < numItems; i++ {
		entries = append(entries, constructMockEntry(i, int64(i), int64(i)))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newOrderedTree(2)
		tree.Add(entries...)
	}
}

func BenchmarkOTQueryItemsMultiDimensions(b *testing.B) {
	numItems := uint64(1000)
	entries := make(Entries, 0, numItems)

	for i := uint64(0); i < numItems; i++ {
		entries = append(entries, constructMockEntry(i, int64(i), int64(i)))
	}

	tree := newOrderedTree(2)
	tree.Add(entries...)
	iv := constructMockInterval(
		dimension{0, int64(numItems)},
		dimension{0, int64(numItems)},
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Query(iv)
	}
}

func TestOTRootDeleteMultiDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(1)
	tree.Delete(entries...)

	assert.Equal(t, uint64(0), tree.Len())

	result := tree.Query(constructMockInterval(dimension{0, 100}, dimension{0, 100}))
	assert.Len(t, result, 0)
}

func TestOTDeleteMultiDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(4)

	tree.Delete(entries[2])

	assert.Equal(t, uint64(3), tree.Len())

	result := tree.Query(constructMockInterval(dimension{0, 4}, dimension{0, 4}))
	assert.Equal(t, Entries{entries[0], entries[1], entries[3]}, result)

	result = tree.Query(constructMockInterval(dimension{3, 4}, dimension{3, 4}))
	assert.Equal(t, Entries{entries[3]}, result)

	result = tree.Query(constructMockInterval(dimension{0, 3}, dimension{0, 3}))
	assert.Equal(t, Entries{entries[0], entries[1]}, result)
}

func TestOTDeleteInOrderMultiDimensions(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(10)

	tree.Delete(entries[5])

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Len(t, result, 9)
	assert.Equal(t, uint64(9), tree.Len())

	assert.NotContains(t, result, entries[5])
}

func TestOTDeleteReverseOrderMultiDimensions(t *testing.T) {
	tree := newOrderedTree(2)

	entries := NewEntries()
	for i := uint64(10); i > 0; i-- {
		entries = append(entries, constructMockEntry(i, int64(i), int64(i)))
	}

	tree.Add(entries...)

	tree.Delete(entries[5])

	result := tree.Query(constructMockInterval(dimension{0, 11}, dimension{0, 11}))
	assert.Len(t, result, 9)
	assert.Equal(t, uint64(9), tree.Len())

	assert.NotContains(t, result, entries[5])
}

func TestOTDeleteRandomOrderMultiDimensions(t *testing.T) {
	tree := newOrderedTree(2)

	entries := NewEntries()
	starts := []uint64{0, 4, 2, 1, 3}
	for _, start := range starts {
		entries = append(entries, constructMockEntry(start, int64(start), int64(start)))
	}

	tree.Add(entries...)

	tree.Delete(entries[2])

	result := tree.Query(constructMockInterval(dimension{0, 11}, dimension{0, 11}))

	assert.Len(t, result, 4)
	assert.Equal(t, uint64(4), tree.Len())

	assert.NotContains(t, result, entries[2])
}

func TestOTDeleteEmptyTreeMultiDimensions(t *testing.T) {
	tree := newOrderedTree(2)

	tree.Delete(constructMockEntry(0, 0, 0))

	assert.Equal(t, uint64(0), tree.Len())
}

func BenchmarkOTDeleteItemsMultiDimensions(b *testing.B) {
	numItems := uint64(1000)
	entries := make(Entries, 0, numItems)

	for i := uint64(0); i < numItems; i++ {
		entries = append(entries, constructMockEntry(i, int64(i), int64(i)))
	}

	trees := make([]*orderedTree, 0, b.N)
	for i := 0; i < b.N; i++ {
		tree := newOrderedTree(2)
		tree.Add(entries...)
		trees = append(trees, tree)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trees[i].Delete(entries...)
	}
}

func TestOverwrites(t *testing.T) {
	tree, _ := constructMultiDimensionalOrderedTree(1)

	entry := constructMockEntry(0, 0, 0)

	tree.Add(entry)

	results := tree.Query(constructMockInterval(dimension{0, 100}, dimension{0, 100}))

	assert.Equal(t, Entries{entry}, results)
	assert.Equal(t, uint64(1), tree.Len())

	newEntry := constructMockEntry(0, 0, 0)

	overwritten := tree.Add(newEntry)
	assert.Equal(t, Entries{entry}, overwritten)
}

func TestTreeApply(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(2)

	result := make(Entries, 0, len(entries))

	tree.Apply(constructMockInterval(dimension{0, 100}, dimension{0, 100}),
		func(e Entry) bool {
			result = append(result, e)
			return true
		},
	)

	assert.Equal(t, entries, result)
}

func TestApplyWithBail(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(2)

	result := make(Entries, 0, 1)

	tree.Apply(constructMockInterval(dimension{0, 100}, dimension{0, 100}),
		func(e Entry) bool {
			result = append(result, e)
			return false
		},
	)

	assert.Equal(t, entries[:1], result)
}

func BenchmarkApply(b *testing.B) {
	numItems := 1000

	tree, _ := constructMultiDimensionalOrderedTree(uint64(numItems))

	iv := constructMockInterval(
		dimension{0, int64(numItems)}, dimension{0, int64(numItems)},
	)
	fn := func(Entry) bool { return true }

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Apply(iv, fn)
	}
}

func TestInsertPositiveIndexFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(2)

	modified, deleted := tree.InsertAtDimension(1, 1, 1)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree.Query(constructMockInterval(dimension{2, 10}, dimension{1, 10}))
	assert.Equal(t, entries[1:], result)
}

func TestInsertPositiveIndexSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 1, 1)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{2, 10}))
	assert.Equal(t, entries[1:], result)
}

func TestInsertPositiveIndexOutOfBoundsFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 4, 1)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))

	assert.Equal(t, entries, result)
}

func TestInsertPositiveIndexOutOfBoundsSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 4, 1)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))

	assert.Equal(t, entries, result)
}

func TestInsertMultiplePositiveIndexFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 1, 2)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree.Query(constructMockInterval(dimension{3, 10}, dimension{1, 10}))
	assert.Equal(t, entries[1:], result)
}

func TestInsertMultiplePositiveIndexSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 1, 2)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{3, 10}))
	assert.Equal(t, entries[1:], result)
}

func TestInsertNegativeIndexFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 1, -1)
	assert.Equal(t, entries[1:2], deleted)
	assert.Equal(t, entries[2:], modified)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{1, 10}))
	assert.Equal(t, entries[2:], result)

	result = tree.Query(constructMockInterval(dimension{2, 10}, dimension{1, 10}))
	assert.Len(t, result, 0)
	assert.Equal(t, uint64(2), tree.Len())
}

func TestInsertNegativeIndexSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 1, -1)
	assert.Equal(t, entries[1:2], deleted)
	assert.Equal(t, entries[2:], modified)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{1, 10}))
	assert.Equal(t, entries[2:], result)

	result = tree.Query(constructMockInterval(dimension{1, 10}, dimension{2, 10}))
	assert.Len(t, result, 0)
	assert.Equal(t, uint64(2), tree.Len())
}

func TestInsertNegativeIndexOutOfBoundsFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 4, -1)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))

	assert.Equal(t, entries, result)
	assert.Equal(t, uint64(3), tree.Len())
}

func TestInsertNegativeIndexOutOfBoundsSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 4, -1)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))

	assert.Equal(t, entries, result)
	assert.Equal(t, uint64(3), tree.Len())
}

func TestInsertMultipleNegativeIndexFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 1, -2)
	assert.Equal(t, entries[1:], deleted)
	assert.Len(t, modified, 0)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{1, 10}))
	assert.Len(t, result, 0)
	assert.Equal(t, uint64(1), tree.Len())
}

func TestInsertMultipleNegativeIndexSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(2, 1, -2)
	assert.Equal(t, entries[1:], deleted)
	assert.Len(t, modified, 0)

	result := tree.Query(constructMockInterval(dimension{1, 10}, dimension{1, 10}))
	assert.Len(t, result, 0)
	assert.Equal(t, uint64(1), tree.Len())
}

func TestInsertInvalidDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(3, 1, -1)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Equal(t, entries, result)
}

func TestInsertInvalidNumber(t *testing.T) {
	tree, entries := constructMultiDimensionalOrderedTree(3)

	modified, deleted := tree.InsertAtDimension(1, 1, 0)
	assert.Len(t, modified, 0)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Equal(t, entries, result)
}

func BenchmarkInsertFirstDimension(b *testing.B) {
	numItems := uint64(100000)

	tree, _ := constructMultiDimensionalOrderedTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.InsertAtDimension(1, 0, 1)
	}
}

func BenchmarkInsertSecondDimension(b *testing.B) {
	numItems := uint64(100000)

	tree, _ := constructMultiDimensionalOrderedTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.InsertAtDimension(2, 0, 1)
	}
}

func BenchmarkDeleteFirstDimension(b *testing.B) {
	numItems := uint64(100000)

	tree, _ := constructMultiDimensionalOrderedTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.InsertAtDimension(1, 0, -1)
	}
}

func BenchmarkDeleteSecondDimension(b *testing.B) {
	numItems := uint64(100000)

	tree, _ := constructMultiDimensionalOrderedTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.InsertAtDimension(2, 0, -1)
	}
}
