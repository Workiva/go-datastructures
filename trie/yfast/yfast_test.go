package yfast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateEntries(num int) Entries {
	entries := make(Entries, 0, num)
	for i := uint64(0); i < uint64(num); i++ {
		entries = append(entries, newMockEntry(i))
	}

	return entries
}

func TestTrieSimpleInsert(t *testing.T) {
	yfast := New(uint8(0))

	e1 := newMockEntry(3)
	e2 := newMockEntry(7)
	e3 := newMockEntry(8)

	yfast.Insert(e1, e2, e3)

	result := yfast.get(3)
	assert.Equal(t, e1, result)

	result = yfast.get(7)
	assert.Equal(t, e2, result)

	result = yfast.get(8)
	assert.Equal(t, e3, result)

	result = yfast.get(250)
	assert.Nil(t, result)

	assert.Equal(t, uint64(3), yfast.Len())
}

func TestTrieDelete(t *testing.T) {
	yfast := New(uint8(0))

	e1 := newMockEntry(3)
	e2 := newMockEntry(7)
	e3 := newMockEntry(8)

	yfast.Insert(e1, e2, e3)

	result := yfast.Delete(3)
	assert.Equal(t, Entries{e1}, result)
	assert.Nil(t, yfast.Get(3))
	assert.Equal(t, uint64(2), yfast.Len())

	result = yfast.Delete(7)
	assert.Equal(t, Entries{e2}, result)
	assert.Nil(t, yfast.Get(7))
	assert.Equal(t, uint64(1), yfast.Len())

	result = yfast.Delete(8)
	assert.Equal(t, Entries{e3}, result)
	assert.Nil(t, yfast.Get(8))
	assert.Equal(t, uint64(0), yfast.Len())

	result = yfast.Delete(5)
	assert.Equal(t, Entries{nil}, result)
	assert.Equal(t, uint64(0), yfast.Len())
}

func TestTrieSuccessor(t *testing.T) {
	yfast := New(uint8(0))

	e3 := newMockEntry(13)
	yfast.Insert(e3)

	successor := yfast.Successor(0)
	assert.Equal(t, e3, successor)

	e1 := newMockEntry(3)
	e2 := newMockEntry(7)

	yfast.Insert(e1, e2)

	successor = yfast.Successor(0)
	assert.Equal(t, e1, successor)

	successor = yfast.Successor(3)
	assert.Equal(t, e1, successor)

	successor = yfast.Successor(4)
	assert.Equal(t, e2, successor)

	successor = yfast.Successor(8)
	assert.Equal(t, e3, successor)

	successor = yfast.Successor(14)
	assert.Nil(t, successor)

	successor = yfast.Successor(100)
	assert.Nil(t, successor)
}

func TestTriePredecessor(t *testing.T) {
	yfast := New(uint8(0))

	predecessor := yfast.Predecessor(5)
	assert.Nil(t, predecessor)

	e1 := newMockEntry(5)
	yfast.Insert(e1)

	predecessor = yfast.Predecessor(13)
	assert.Equal(t, e1, predecessor)

	e2 := newMockEntry(12)
	yfast.Insert(e2)

	predecessor = yfast.Predecessor(11)
	assert.Equal(t, e1, predecessor)

	predecessor = yfast.Predecessor(5)
	assert.Equal(t, e1, predecessor)

	predecessor = yfast.Predecessor(4)
	assert.Nil(t, predecessor)

	predecessor = yfast.Predecessor(100)
	assert.Equal(t, e2, predecessor)
}

func BenchmarkInsert(b *testing.B) {
	yfast := New(uint64(0))
	entries := generateEntries(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		yfast.Insert(entries[i])
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := 1000

	entries := generateEntries(numItems)

	yfast := New(uint32(0))
	yfast.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		yfast.Get(uint64(numItems / 2))
	}
}

func BenchmarkDelete(b *testing.B) {
	entries := generateEntries(b.N)
	yfast := New(uint64(0))
	yfast.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		yfast.Delete(uint64(i))
	}
}

func BenchmarkSuccessor(b *testing.B) {
	numItems := 100000

	entries := make(Entries, 0, numItems)
	for i := uint64(0); i < uint64(numItems); i++ {
		entries = append(entries, newMockEntry(i+uint64(b.N/2)))
	}

	yfast := New(uint64(0))
	yfast.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		yfast.Successor(uint64(i))
	}
}

func BenchmarkPredecessor(b *testing.B) {
	numItems := 100000

	entries := make(Entries, 0, numItems)
	for i := uint64(0); i < uint64(numItems); i++ {
		entries = append(entries, newMockEntry(i+uint64(b.N/2)))
	}

	yfast := New(uint64(0))
	yfast.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		yfast.Predecessor(uint64(i))
	}
}
