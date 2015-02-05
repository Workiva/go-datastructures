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

package plus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchKeys(t *testing.T) {
	keys := keys{newMockKey(1), newMockKey(2), newMockKey(4)}

	testKey := newMockKey(5)
	assert.Equal(t, 3, keySearch(keys, testKey))

	testKey = newMockKey(2)
	assert.Equal(t, 1, keySearch(keys, testKey))

	testKey = newMockKey(0)
	assert.Equal(t, 0, keySearch(keys, testKey))

	testKey = newMockKey(3)
	assert.Equal(t, 2, keySearch(keys, testKey))

	assert.Equal(t, 0, keySearch(nil, testKey))
}

func TestTreeInsert2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)

	tree.Insert(keys...)

	assert.Len(t, tree.root.(*inode).keys, 2)
	assert.Len(t, tree.root.(*inode).nodes, 3)
	assert.IsType(t, &inode{}, tree.root)
}

func TestTreeInsert3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	assert.Len(t, tree.root.(*inode).keys, 1)
	assert.Len(t, tree.root.(*inode).nodes, 2)
	assert.IsType(t, &inode{}, tree.root)
}

func TestTreeInsertQuery2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsertQuery3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsertReverseOrder2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()
	keys.reverse() // we want to fetch things in the correct
	// ascending order

	assert.Equal(t, keys, result)
}

func TestTreeInsertReverseOrder3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()
	keys.reverse() // we want to fetch things in the correct
	// ascending order

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5_WithEndDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)
	duplicate := newMockKey(4)
	tree.Insert(duplicate)
	keys[4] = duplicate

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5_WithMiddleDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)
	duplicate := newMockKey(2)
	tree.Insert(duplicate)
	keys[2] = duplicate

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5WithEarlyDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)
	duplicate := newMockKey(0)
	tree.Insert(duplicate)
	keys[0] = duplicate

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5WithDuplicateID(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	key := newMockKey(2)
	tree.Insert(keys...)
	tree.Insert(key)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5MiddleQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(2))
	result := iter.exhaust()

	assert.Equal(t, keys[2:], result)
}

func TestTreeInsert3_4_5LateQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(4))
	result := iter.exhaust()

	assert.Equal(t, keys[4:], result)
}

func TestTreeInsert3_4_5AfterQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(5))
	result := iter.exhaust()

	assert.Len(t, result, 0)
}

func TestTreeInternalNodeSplit(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(10)

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInternalNodeSplitReverseOrder(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(10)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()
	keys.reverse()

	assert.Equal(t, keys, result)
}

func TestTreeInternalNodeSplitRandomOrder(t *testing.T) {
	ids := []int{6, 2, 9, 0, 3, 4, 7, 1, 8, 5}
	keys := make(keys, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, newMockKey(id))
	}

	tree := newBTree(4)
	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(0))
	result := iter.exhaust()

	assert.Len(t, result, 10)
	for i, key := range result {
		assert.Equal(t, newMockKey(i), key)
	}
}

func TestTreeRandomOrderQuery(t *testing.T) {
	ids := []int{6, 2, 9, 0, 3, 4, 7, 1, 8, 5}
	keys := make(keys, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, newMockKey(id))
	}

	tree := newBTree(4)
	tree.Insert(keys...)

	iter := tree.Iter(newMockKey(4))
	result := iter.exhaust()

	assert.Len(t, result, 6)
	for i, key := range result {
		assert.Equal(t, newMockKey(i+4), key)
	}
}

func TestTreeGet(t *testing.T) {
	keys := constructRandomMockKeys(100)
	tree := newBTree(64)
	tree.Insert(keys...)

	assert.Equal(t, uint64(100), tree.Len())
	assert.Equal(t, keys, tree.Get(keys...))
}

func TestTreeGetNotFound(t *testing.T) {
	keys := constructMockKeys(5)
	tree := newBTree(64)
	tree.Insert(keys...)

	assert.Equal(t, Keys{nil}, tree.Get(newMockKey(20)))
}

func TestGetExactMatchesOnly(t *testing.T) {
	k1 := newMockKey(0)
	k2 := newMockKey(5)
	tree := newBTree(64)
	tree.Insert(k1, k2)

	assert.Equal(t, Keys{nil}, tree.Get(newMockKey(3)))
}

func BenchmarkIteration(b *testing.B) {
	numItems := 1000
	ary := uint64(16)

	keys := constructMockKeys(numItems)
	tree := newBTree(ary)
	tree.Insert(keys...)

	searchKey := newMockKey(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		iter := tree.Iter(searchKey)
		iter.exhaust()
	}
}

func BenchmarkInsert(b *testing.B) {
	numItems := b.N
	ary := uint64(16)

	keys := constructMockKeys(numItems)
	tree := newBTree(ary)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i%numItems])
	}
}

func BenchmarkBulkAdd(b *testing.B) {
	numItems := 10000
	keys := constructRandomMockKeys(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newBTree(128)
		tree.Insert(keys...)
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := b.N
	ary := uint64(16)

	keys := constructMockKeys(numItems)
	tree := newBTree(ary)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys[i%numItems])
	}
}
