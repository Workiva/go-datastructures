package skip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStarInsert(t *testing.T) {
	ssl := NewStar(uint8(0))

	e1 := newMockEntry(7)
	e2 := newMockEntry(10)

	result := ssl.Insert(e1, e2)
	assert.Equal(t, Entries{nil, nil}, result)
	assert.Equal(t, Entries{e1}, ssl.Get(7))
	assert.Equal(t, Entries{e2}, ssl.Get(10))
	assert.Equal(t, Entries{e1, e2}, ssl.Get(7, 10))
	assert.Equal(t, Entries{e1, nil}, ssl.Get(7, 13))
	assert.Equal(t, Entries{e2, nil}, ssl.Get(10, 13))
	assert.Equal(t, uint64(2), ssl.Len())
}

func TestStarOverwrite(t *testing.T) {
	ssl := NewStar(uint8(0))
	e1 := newMockEntry(7)
	e2 := newMockEntry(7)

	result := ssl.Insert(e1)
	assert.Equal(t, Entries{nil}, result)
	assert.Equal(t, uint64(1), ssl.Len())

	result = ssl.Insert(e2)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, uint64(1), ssl.Len())
}

func TestStarDelete(t *testing.T) {
	ssl := NewStar(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)
	ssl.Insert(e1, e2)

	result := ssl.Delete(e1.Key(), e2.Key())
	assert.Equal(t, Entries{e1, e2}, result)
	assert.Equal(t, uint64(0), ssl.Len())
}

func TestStarIter(t *testing.T) {
	ssl := NewStar(uint8(0))

	iter := ssl.Iter(0)
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e1 := newMockEntry(5)
	e2 := newMockEntry(10)
	ssl.Insert(e1, e2)

	iter = ssl.Iter(0)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())

	iter = ssl.Iter(5)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())

	iter = ssl.Iter(6)
	assert.Equal(t, Entries{e2}, iter.exhaust())

	iter = ssl.Iter(10)
	assert.Equal(t, Entries{e2}, iter.exhaust())

	iter = ssl.Iter(11)
	assert.Equal(t, Entries{}, iter.exhaust())
}

func BenchmarkStarInsert(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries[i%numItems])
	}
}

func BenchmarkStarGet(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems].Key())
	}
}
