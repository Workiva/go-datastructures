package skiplist

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/rangetree"
)

func generateMultiDimensionalEntries(num int) rangetree.Entries {
	entries := make(rangetree.Entries, 0, num)
	for i := 0; i < num; i++ {
		entries = append(entries, newMockEntry(int64(i), int64(i)))
	}

	return entries
}

func generateRandomMultiDimensionalEntries(num int) rangetree.Entries {
	entries := make(rangetree.Entries, 0, num)
	for i := 0; i < num; i++ {
		value := rand.Int63()
		entries = append(entries, newMockEntry(value, value))
	}

	return entries
}

func TestRTSingleDimensionAdd(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(3)
	m2 := newMockEntry(5)

	overwritten := rt.Add(m1, m2)
	assert.Equal(t, rangetree.Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(2), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, m2}, rt.Get(m1, m2))
}

func TestRTMultiDimensionAdd(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 5)
	m2 := newMockEntry(4, 6)

	overwritten := rt.Add(m1, m2)
	assert.Equal(t, rangetree.Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(2), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, m2}, rt.Get(m1, m2))
}

func TestRTSingleDimensionOverwrite(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(5)
	m2 := newMockEntry(5)

	overwritten := rt.Add(m1)
	assert.Equal(t, rangetree.Entries{nil}, overwritten)
	assert.Equal(t, uint64(1), rt.Len())

	overwritten = rt.Add(m2)
	assert.Equal(t, rangetree.Entries{m1}, overwritten)
	assert.Equal(t, uint64(1), rt.Len())
	assert.Equal(t, rangetree.Entries{m2}, rt.Get(m2))
}

func TestRTMultiDimensionOverwrite(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(5, 6)
	m2 := newMockEntry(5, 6)

	overwritten := rt.Add(m1)
	assert.Equal(t, rangetree.Entries{nil}, overwritten)
	assert.Equal(t, uint64(1), rt.Len())

	overwritten = rt.Add(m2)
	assert.Equal(t, rangetree.Entries{m1}, overwritten)
	assert.Equal(t, uint64(1), rt.Len())
	assert.Equal(t, rangetree.Entries{m2}, rt.Get(m2))
}

func TestRTSingleDimensionDelete(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(5)
	m2 := newMockEntry(2)
	rt.Add(m1, m2)

	rt.Delete(m1, m2)
	assert.Equal(t, uint64(0), rt.Len())
	assert.Equal(t, rangetree.Entries{nil, nil}, rt.Get(m1, m2))
}

func TestRTMultiDimensionDelete(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 5)
	m2 := newMockEntry(4, 6)
	rt.Add(m1, m2)

	rt.Delete(m1, m2)
	assert.Equal(t, uint64(0), rt.Len())
	assert.Equal(t, rangetree.Entries{nil, nil}, rt.Get(m1, m2))
}

func TestRTSingleDimensionQuery(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(3)
	m2 := newMockEntry(6)
	m3 := newMockEntry(9)
	rt.Add(m1, m2, m3)

	result := rt.Query(newMockInterval([]int64{1}, []int64{7}))
	assert.Equal(t, rangetree.Entries{m1, m2}, result)

	result = rt.Query(newMockInterval([]int64{6}, []int64{10}))
	assert.Equal(t, rangetree.Entries{m2, m3}, result)

	result = rt.Query(newMockInterval([]int64{9}, []int64{11}))
	assert.Equal(t, rangetree.Entries{m3}, result)

	result = rt.Query(newMockInterval([]int64{0}, []int64{3}))
	assert.Len(t, result, 0)

	result = rt.Query(newMockInterval([]int64{10}, []int64{13}))
	assert.Len(t, result, 0)
}

func TestRTMultiDimensionQuery(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	m2 := newMockEntry(6, 6)
	m3 := newMockEntry(9, 9)
	rt.Add(m1, m2, m3)

	result := rt.Query(newMockInterval([]int64{1, 1}, []int64{7, 7}))
	assert.Equal(t, rangetree.Entries{m1, m2}, result)

	result = rt.Query(newMockInterval([]int64{6, 6}, []int64{10, 10}))
	assert.Equal(t, rangetree.Entries{m2, m3}, result)

	result = rt.Query(newMockInterval([]int64{9, 9}, []int64{11, 11}))
	assert.Equal(t, rangetree.Entries{m3}, result)

	result = rt.Query(newMockInterval([]int64{0, 0}, []int64{3, 3}))
	assert.Len(t, result, 0)

	result = rt.Query(newMockInterval([]int64{10, 10}, []int64{13, 13}))
	assert.Len(t, result, 0)

	result = rt.Query(newMockInterval([]int64{0, 0}, []int64{3, 3}))
	assert.Len(t, result, 0)

	result = rt.Query(newMockInterval([]int64{6, 1}, []int64{7, 6}))
	assert.Len(t, result, 0)

	result = rt.Query(newMockInterval([]int64{0, 0}, []int64{7, 4}))
	assert.Equal(t, rangetree.Entries{m1}, result)
}

func TestRTSingleDimensionInsert(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(3)
	m2 := newMockEntry(6)
	m3 := newMockEntry(9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(0, 0, 1)
	assert.Equal(t, rangetree.Entries{m1, m2, m3}, affected)
	assert.Len(t, deleted, 0)
	assert.Equal(t, uint64(3), rt.Len())
	assert.Equal(t, rangetree.Entries{nil, nil, nil}, rt.Get(m1, m2, m3))
	e1 := newMockEntry(4)
	e2 := newMockEntry(7)
	e3 := newMockEntry(10)
	assert.Equal(t, rangetree.Entries{m1, m2, m3}, rt.Get(e1, e2, e3))
}

func TestRTSingleDimensionInsertNegative(t *testing.T) {
	rt := new(1)
	m1 := newMockEntry(3)
	m2 := newMockEntry(6)
	m3 := newMockEntry(9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(0, 6, -2)
	assert.Equal(t, rangetree.Entries{m3}, affected)
	assert.Equal(t, rangetree.Entries{m2}, deleted)
	assert.Equal(t, uint64(2), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, nil}, rt.Get(m1, m2))

	e2 := newMockEntry(4)
	e3 := newMockEntry(7)
	assert.Equal(t, rangetree.Entries{nil, m3}, rt.Get(e2, e3))
}

func TestRTMultiDimensionInsert(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	m2 := newMockEntry(6, 6)
	m3 := newMockEntry(9, 9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(1, 4, 2)
	assert.Equal(t, rangetree.Entries{m2, m3}, affected)
	assert.Len(t, deleted, 0)
	assert.Equal(t, uint64(3), rt.Len())

	e2 := newMockEntry(6, 8)
	e3 := newMockEntry(9, 11)
	assert.Equal(t, rangetree.Entries{m1, nil, nil}, rt.Get(m1, m2, m3))
	assert.Equal(t, rangetree.Entries{m2, m3}, rt.Get(e2, e3))
}

func TestRTMultiDimensionInsertNegative(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	m2 := newMockEntry(6, 6)
	m3 := newMockEntry(9, 9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(1, 6, -2)
	assert.Equal(t, rangetree.Entries{m3}, affected)
	assert.Equal(t, rangetree.Entries{m2}, deleted)
	assert.Equal(t, uint64(2), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, nil, nil}, rt.Get(m1, m2, m3))

	e2 := newMockEntry(6, 4)
	e3 := newMockEntry(9, 7)
	assert.Equal(t, rangetree.Entries{nil, m3}, rt.Get(e2, e3))
}

func TestRTInsertInZeroDimensionMultiDimensionList(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	m2 := newMockEntry(6, 6)
	m3 := newMockEntry(9, 9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(0, 4, 2)
	assert.Equal(t, rangetree.Entries{m2, m3}, affected)
	assert.Len(t, deleted, 0)
	assert.Equal(t, uint64(3), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, nil, nil}, rt.Get(m1, m2, m3))

	e2 := newMockEntry(8, 6)
	e3 := newMockEntry(11, 9)
	assert.Equal(t, rangetree.Entries{m2, m3}, rt.Get(e2, e3))
}

func TestRTInsertNegativeInZeroDimensionMultiDimensionList(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	m2 := newMockEntry(6, 6)
	m3 := newMockEntry(9, 9)
	rt.Add(m1, m2, m3)

	affected, deleted := rt.InsertAtDimension(0, 6, -2)
	assert.Equal(t, rangetree.Entries{m3}, affected)
	assert.Equal(t, rangetree.Entries{m2}, deleted)
	assert.Equal(t, uint64(2), rt.Len())
	assert.Equal(t, rangetree.Entries{m1, nil, nil}, rt.Get(m1, m2, m3))

	e2 := newMockEntry(4, 6)
	e3 := newMockEntry(7, 9)
	assert.Equal(t, rangetree.Entries{nil, m3}, rt.Get(e2, e3))
}

func TestRTInsertBeyondDimension(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	rt.Add(m1)

	affected, deleted := rt.InsertAtDimension(4, 0, 1)
	assert.Len(t, affected, 0)
	assert.Len(t, deleted, 0)
	assert.Equal(t, rangetree.Entries{m1}, rt.Get(m1))
}

func TestRTInsertZero(t *testing.T) {
	rt := new(2)
	m1 := newMockEntry(3, 3)
	rt.Add(m1)

	affected, deleted := rt.InsertAtDimension(1, 0, 0)
	assert.Len(t, affected, 0)
	assert.Len(t, deleted, 0)
	assert.Equal(t, rangetree.Entries{m1}, rt.Get(m1))
}

func BenchmarkMultiDimensionInsert(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateMultiDimensionalEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.Add(entries[i%numItems])
	}
}

func BenchmarkMultiDimensionInsertReverse(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateMultiDimensionalEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		index := numItems - (i % numItems) - 1
		rt.Add(entries[index])
	}
}

func BenchmarkMultiDimensionRandomInsert(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.Add(entries[i%numItems])
	}
}

func BenchmarkMultiDimensionalGet(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)
	rt.Add(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.Get(entries[i%numItems])
	}
}

func BenchmarkMultiDimensionDelete(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)
	rt.Add(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.Delete(entries[i%numItems])
	}
}

func BenchmarkMultiDimensionQuery(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)
	rt.Add(entries...)
	iv := newMockInterval([]int64{0, 0}, []int64{math.MaxInt64, math.MaxInt64})
	var result rangetree.Entries

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result = rt.Query(iv)
	}

	assert.Len(b, result, numItems)
}

func BenchmarkMultiDimensionInsertAtZeroDimension(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)
	rt.Add(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.InsertAtDimension(0, 0, 1)
	}
}

func BenchmarkMultiDimensionInsertNegativeAtZeroDimension(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateRandomMultiDimensionalEntries(numItems)
	rt.Add(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.InsertAtDimension(0, 0, -1)
	}
}
