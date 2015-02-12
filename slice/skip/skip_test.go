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

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/common"
)

func generateMockEntries(num int) common.Comparators {
	entries := make(common.Comparators, 0, num)
	for i := uint64(0); i < uint64(num); i++ {
		entries = append(entries, newMockEntry(i))
	}

	return entries
}

func generateRandomMockEntries(num int) common.Comparators {
	entries := make(common.Comparators, 0, num)
	for i := 0; i < num; i++ {
		entries = append(entries, newMockEntry(uint64(rand.Int())))
	}

	return entries
}

func TestInsertByPosition(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	m3 := newMockEntry(2)
	sl := New(uint8(0))
	sl.InsertAtPosition(2, m1)
	sl.InsertAtPosition(0, m2)
	sl.InsertAtPosition(0, m3)

	assert.Equal(t, m3, sl.ByPosition(0))
	assert.Equal(t, m2, sl.ByPosition(1))
	assert.Equal(t, m1, sl.ByPosition(2))
	assert.Nil(t, sl.ByPosition(3))
}

func TestGetByPosition(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))
	sl.Insert(m1, m2)

	assert.Equal(t, m1, sl.ByPosition(0))
	assert.Equal(t, m2, sl.ByPosition(1))
	assert.Nil(t, sl.ByPosition(2))
}

func TestSplitAt(t *testing.T) {
	m1 := newMockEntry(3)
	m2 := newMockEntry(5)
	m3 := newMockEntry(7)
	sl := New(uint8(0))
	sl.Insert(m1, m2, m3)

	left, right := sl.SplitAt(1)
	assert.Equal(t, uint64(2), left.Len())
	assert.Equal(t, uint64(1), right.Len())
	assert.Equal(t, common.Comparators{m1, m2, nil}, left.Get(m1, m2, m3))
	assert.Equal(t, common.Comparators{nil, nil, m3}, right.Get(m1, m2, m3))
	assert.Equal(t, m1, left.ByPosition(0))
	assert.Equal(t, m2, left.ByPosition(1))
	assert.Equal(t, m3, right.ByPosition(0))
	assert.Equal(t, nil, left.ByPosition(2))
	assert.Equal(t, nil, right.ByPosition(1))
}

func TestSplitLargeSkipList(t *testing.T) {
	entries := generateMockEntries(100)
	leftEntries := entries[:50]
	rightEntries := entries[50:]
	sl := New(uint64(0))
	sl.Insert(entries...)

	left, right := sl.SplitAt(49)
	assert.Equal(t, uint64(50), left.Len())
	for _, le := range leftEntries {
		result, index := left.GetWithPosition(le)
		assert.Equal(t, le, result)
		assert.Equal(t, le, left.ByPosition(index))
	}

	assert.Equal(t, uint64(50), right.Len())
	for _, re := range rightEntries {
		result, index := right.GetWithPosition(re)
		assert.Equal(t, re, result)
		assert.Equal(t, re, right.ByPosition(index))
	}
}

func TestSplitLargeSkipListOddNumber(t *testing.T) {
	entries := generateMockEntries(99)
	leftEntries := entries[:50]
	rightEntries := entries[50:]
	sl := New(uint64(0))
	sl.Insert(entries...)

	left, right := sl.SplitAt(49)
	assert.Equal(t, uint64(50), left.Len())
	for _, le := range leftEntries {
		result, index := left.GetWithPosition(le)
		assert.Equal(t, le, result)
		assert.Equal(t, le, left.ByPosition(index))
	}

	assert.Equal(t, uint64(49), right.Len())
	for _, re := range rightEntries {
		result, index := right.GetWithPosition(re)
		assert.Equal(t, re, result)
		assert.Equal(t, re, right.ByPosition(index))
	}
}

func TestSplitAtSkipListLength(t *testing.T) {
	entries := generateMockEntries(5)
	sl := New(uint64(0))
	sl.Insert(entries...)

	left, right := sl.SplitAt(4)
	assert.Equal(t, sl, left)
	assert.Nil(t, right)
}

func TestGetWithPosition(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))
	sl.Insert(m1, m2)

	e, pos := sl.GetWithPosition(m1)
	assert.Equal(t, m1, e)
	assert.Equal(t, uint64(0), pos)

	e, pos = sl.GetWithPosition(m2)
	assert.Equal(t, m2, e)
	assert.Equal(t, uint64(1), pos)
}

func TestReplaceAtPosition(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))

	sl.Insert(m1, m2)
	m3 := newMockEntry(9)
	sl.ReplaceAtPosition(0, m3)
	assert.Equal(t, m3, sl.ByPosition(0))
	assert.Equal(t, m2, sl.ByPosition(1))
}

func TestInsertRandomGetByPosition(t *testing.T) {
	entries := generateRandomMockEntries(100)
	sl := New(uint64(0))
	sl.Insert(entries...)

	for _, e := range entries {
		_, pos := sl.GetWithPosition(e)
		assert.Equal(t, e, sl.ByPosition(pos))
	}
}

func TestGetManyByPosition(t *testing.T) {
	entries := generateMockEntries(10)
	sl := New(uint64(0))
	sl.Insert(entries...)

	for i, e := range entries {
		assert.Equal(t, e, sl.ByPosition(uint64(i)))
	}
}

func TestGetPositionAfterDelete(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))
	sl.Insert(m1, m2)

	sl.Delete(m1)
	assert.Equal(t, m2, sl.ByPosition(0))
	assert.Nil(t, sl.ByPosition(1))

	sl.Delete(m2)
	assert.Nil(t, sl.ByPosition(0))
	assert.Nil(t, sl.ByPosition(1))
}

func TestGetPositionBulkDelete(t *testing.T) {
	es := generateMockEntries(20)
	e1 := es[:10]
	e2 := es[10:]
	sl := New(uint64(0))
	sl.Insert(e1...)
	sl.Insert(e2...)

	for _, e := range e1 {
		sl.Delete(e)
	}
	for i, e := range e2 {
		assert.Equal(t, e, sl.ByPosition(uint64(i)))
	}
}

func TestSimpleInsert(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1)
	assert.Equal(t, common.Comparators{m1}, sl.Get(m1))
	assert.Equal(t, uint64(1), sl.Len())
	assert.Equal(t, common.Comparators{nil}, overwritten)
	assert.Equal(t, common.Comparators{nil}, sl.Get(mockEntry(1)))

	overwritten = sl.Insert(m2)
	assert.Equal(t, common.Comparators{m2}, sl.Get(m2))
	assert.Equal(t, common.Comparators{nil}, sl.Get(mockEntry(7)))
	assert.Equal(t, uint64(2), sl.Len())
	assert.Equal(t, common.Comparators{nil}, overwritten)
}

func TestSimpleOverwrite(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(5)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1)
	assert.Equal(t, common.Comparators{nil}, overwritten)
	assert.Equal(t, uint64(1), sl.Len())

	overwritten = sl.Insert(m2)
	assert.Equal(t, common.Comparators{m1}, overwritten)
	assert.Equal(t, uint64(1), sl.Len())
}

func TestInsertOutOfOrder(t *testing.T) {
	m1 := newMockEntry(6)
	m2 := newMockEntry(5)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1, m2)
	assert.Equal(t, common.Comparators{nil, nil}, overwritten)

	assert.Equal(t, common.Comparators{m1, m2}, sl.Get(m1, m2))
}

func TestSimpleDelete(t *testing.T) {
	m1 := newMockEntry(5)
	sl := New(uint8(0))
	sl.Insert(m1)

	deleted := sl.Delete(m1)
	assert.Equal(t, common.Comparators{m1}, deleted)
	assert.Equal(t, uint64(0), sl.Len())
	assert.Equal(t, common.Comparators{nil}, sl.Get(m1))

	deleted = sl.Delete(m1)
	assert.Equal(t, common.Comparators{nil}, deleted)
}

func TestDeleteAll(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))
	sl.Insert(m1, m2)

	deleted := sl.Delete(m1, m2)
	assert.Equal(t, common.Comparators{m1, m2}, deleted)
	assert.Equal(t, uint64(0), sl.Len())
	assert.Equal(t, common.Comparators{nil, nil}, sl.Get(m1, m2))
}

func TestIter(t *testing.T) {
	sl := New(uint8(0))
	m1 := newMockEntry(5)
	m2 := newMockEntry(10)

	sl.Insert(m1, m2)

	iter := sl.Iter(mockEntry(0))
	assert.Equal(t, common.Comparators{m1, m2}, iter.exhaust())

	iter = sl.Iter(mockEntry(5))
	assert.Equal(t, common.Comparators{m1, m2}, iter.exhaust())

	iter = sl.Iter(mockEntry(6))
	assert.Equal(t, common.Comparators{m2}, iter.exhaust())

	iter = sl.Iter(mockEntry(10))
	assert.Equal(t, common.Comparators{m2}, iter.exhaust())

	iter = sl.Iter(mockEntry(11))
	assert.Equal(t, common.Comparators{}, iter.exhaust())
}

func TestIterAtPosition(t *testing.T) {
	sl := New(uint8(0))
	m1 := newMockEntry(5)
	m2 := newMockEntry(10)

	sl.Insert(m1, m2)

	iter := sl.IterAtPosition(0)
	assert.Equal(t, common.Comparators{m1, m2}, iter.exhaust())

	iter = sl.IterAtPosition(1)
	assert.Equal(t, common.Comparators{m2}, iter.exhaust())

	iter = sl.IterAtPosition(2)
	assert.Equal(t, common.Comparators{}, iter.exhaust())
}

func BenchmarkInsert(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries[i%numItems])
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems])
	}
}

func BenchmarkDelete(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(entries[i])
	}
}

func BenchmarkPrepend(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))

	entries := make(common.Comparators, 0, numItems)
	for i := b.N; i < b.N+numItems; i++ {
		entries = append(entries, newMockEntry(uint64(i)))
	}

	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(newMockEntry(uint64(i)))
	}
}

func BenchmarkByPosition(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))
	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.ByPosition(uint64(i % numItems))
	}
}

func BenchmarkInsertAtPosition(b *testing.B) {
	numItems := b.N
	sl := New(uint64(0))
	entries := generateRandomMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.InsertAtPosition(0, entries[i%numItems])
	}
}
