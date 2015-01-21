package skip

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

func dumpNodes(n *node) {
	for n != nil {
		log.Printf(`N: %+v, N POINTER: %p`, n, n)
		n = n.forward[0]
	}
}

func TestSimpleImmutableInsert(t *testing.T) {
	sl1 := NewImmutable(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)

	sl2, overwritten := sl1.Insert(e1)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, uint64(0), sl1.Len())
	assert.Equal(t, uint64(1), sl2.Len())
	assert.Equal(t, Entries{e1}, sl2.Get(5))
	assert.Equal(t, Entries{nil}, sl1.Get(5))

	sl3, overwritten := sl2.Insert(e2)

	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, Entries{e1, e2}, sl3.Get(5, 10))
	assert.Equal(t, Entries{e1, nil}, sl2.Get(5, 10))
	assert.Equal(t, uint64(2), sl3.Len())
	assert.Equal(t, uint64(1), sl2.Len())

	e3 := newMockEntry(3)

	sl4, overwritten := sl3.Insert(e3)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, Entries{e3}, sl4.Get(3))
	assert.Equal(t, Entries{nil}, sl3.Get(3))
}

func TestImmutableOverwrite(t *testing.T) {
	sl1 := NewImmutable(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(5)

	sl2, overwritten := sl1.Insert(e1)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, Entries{nil}, sl1.Get(5))
	assert.Equal(t, Entries{e1}, sl2.Get(5))

	sl3, overwritten := sl2.Insert(e2)
	assert.Equal(t, Entries{e1}, overwritten)
	assert.Equal(t, Entries{e1}, sl2.Get(5))
	assert.Equal(t, Entries{e2}, sl3.Get(5))

	e3 := newMockEntry(3)
	sl4, overwritten := sl3.Insert(e3)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, Entries{e3}, sl4.Get(3))
	assert.Equal(t, Entries{nil}, sl3.Get(3))
}

func TestImmutableDelete(t *testing.T) {
	sl1 := NewImmutable(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)

	sl2, _ := sl1.Insert(e1, e2)

	sl3, _ := sl2.Delete(e1.Key())
	assert.Equal(t, uint64(2), sl2.Len())
	assert.Equal(t, uint64(1), sl3.Len())
	assert.Equal(t, Entries{e1, e2}, sl2.Get(e1.Key(), e2.Key()))
	assert.Equal(t, Entries{nil, e2}, sl3.Get(e1.Key(), e2.Key()))

	sl4, _ := sl3.Delete(e2.Key())
	assert.Equal(t, uint64(1), sl3.Len())
	assert.Equal(t, uint64(0), sl4.Len())
	assert.Equal(t, Entries{nil, e2}, sl3.Get(e1.Key(), e2.Key()))
	assert.Equal(t, Entries{nil, nil}, sl4.Get(e1.Key(), e2.Key()))
}

func TestImmutableBulkDelete(t *testing.T) {
	sl1 := NewImmutable(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)
	e3 := newMockEntry(15)
	sl2, _ := sl1.Insert(e1)
	sl2, _ = sl2.Insert(e2)
	sl2, _ = sl2.Insert(e3)
	dumpNodes(sl2.head)

	/*
		sl3, _ := sl2.Delete(e1.Key())
		assert.Equal(t, Entries{e1, e2, e3}, sl2.Get(e1.Key(), e2.Key(), e3.Key()))
		sl3, _ = sl3.Delete(e2.Key())
		assert.Equal(t, Entries{e1, e2, e3}, sl2.Get(e1.Key(), e2.Key(), e3.Key()))
		dumpNodes(sl2.head)
		sl3, _ = sl3.Delete(e3.Key())
		assert.Equal(t, Entries{e1, e2, e3}, sl2.Get(e1.Key(), e2.Key(), e3.Key()))
		//assert.Equal(t, Entries{e1, e2, e3}, deleted)
		//assert.Equal(t, uint64(3), sl2.Len())
		assert.Equal(t, uint64(0), sl3.Len())
		assert.Equal(t, Entries{nil, nil, nil}, sl3.Get(e1.Key(), e2.Key(), e3.Key()))*/
}

func BenchmarkImmutableInsert(b *testing.B) {
	numItems := b.N
	sl := NewImmutable(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl, _ = sl.Insert(entries[i%numItems])
	}
}

func BenchmarkBulkImmutableInsert(b *testing.B) {
	numItems := b.N
	sl := NewImmutable(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries...) // every sl should be empty when this is called
	}
}

func BenchmarkBulkImmutableDelete(b *testing.B) {
	numItems := b.N
	sl := NewImmutable(uint64(0))

	entries := generateMockEntries(numItems)
	keys := make([]uint64, 0, len(entries))
	for _, e := range entries {
		keys = append(keys, e.Key())
	}
	sl, _ = sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(keys...) // every sl should be empty when this is called
	}
}

func BenchmarkImmutableGet(b *testing.B) {
	numItems := b.N
	sl := NewImmutable(uint64(0))

	entries := generateMockEntries(numItems)
	sl, _ = sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems].Key())
	}
}

func BenchmarkImmutableDelete(b *testing.B) {
	numItems := b.N
	sl := NewImmutable(uint64(0))

	entries := generateMockEntries(numItems)
	sl, _ = sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl, _ = sl.Delete(entries[i].Key())
	}
}
