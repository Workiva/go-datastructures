package skiplist

import (
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

func BenchmarkMultiDimensionInsert(b *testing.B) {
	numItems := b.N
	rt := new(2)
	entries := generateMultiDimensionalEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rt.Add(entries[i%numItems])
	}
}
