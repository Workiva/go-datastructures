package avl

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

func generateMockEntries(num int) Entries {
	entries := make(Entries, 0, num)
	for i := 0; i < num; i++ {
		entries = append(entries, mockEntry(i))
	}

	return entries
}

func TestAVLSimpleInsert(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(5)
	m2 := mockEntry(10)

	i2, overwritten := i1.Insert(m1, m2)
	assert.Equal(t, Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(2), i2.Len())
	assert.Equal(t, uint64(0), i1.Len())
	assert.Equal(t, Entries{nil, nil}, i1.Get(m1, m2))
	assert.Equal(t, Entries{m1, m2}, i2.Get(m1, m2))

	m3 := mockEntry(1)

	i3, overwritten := i2.Insert(m3)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, uint64(3), i3.Len())
	assert.Equal(t, uint64(2), i2.Len())
	assert.Equal(t, uint64(0), i1.Len())
	assert.Equal(t, Entries{m1, m2, m3}, i3.Get(m1, m2, m3))
}

func TestAVLInsertRightLeaning(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(1)
	m2 := mockEntry(5)
	m3 := mockEntry(10)

	i2, overwritten := i1.Insert(m1, m2, m3)
	assert.Equal(t, Entries{nil, nil, nil}, overwritten)
	assert.Equal(t, uint64(0), i1.Len())
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{m1, m2, m3}, i2.Get(m1, m2, m3))
	assert.Equal(t, Entries{nil, nil, nil}, i1.Get(m1, m2, m3))

	m4 := mockEntry(15)
	m5 := mockEntry(20)

	i3, overwritten := i2.Insert(m4, m5)
	assert.Equal(t, Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(5), i3.Len())
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{nil, nil}, i2.Get(m4, m5))
	assert.Equal(t, Entries{m4, m5}, i3.Get(m4, m5))
}

func TestAVLInsertRightLeaningDoubleRotation(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(1)
	m2 := mockEntry(10)
	m3 := mockEntry(5)

	i2, overwritten := i1.Insert(m1, m2, m3)
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{nil, nil, nil}, overwritten)
	assert.Equal(t, Entries{nil, nil, nil}, i1.Get(m1, m2, m3))
	assert.Equal(t, Entries{m1, m2, m3}, i2.Get(m1, m2, m3))

	m4 := mockEntry(20)
	m5 := mockEntry(15)

	i3, overwritten := i2.Insert(m4, m5)
	assert.Equal(t, Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(5), i3.Len())
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{nil, nil}, i2.Get(m4, m5))
	assert.Equal(t, Entries{m4, m5}, i3.Get(m4, m5))
}

func TestAVLInsertLeftLeaning(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(20)
	m2 := mockEntry(15)
	m3 := mockEntry(10)

	i2, overwritten := i1.Insert(m1, m2, m3)
	assert.Equal(t, Entries{nil, nil, nil}, overwritten)
	assert.Equal(t, uint64(0), i1.Len())
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{nil, nil, nil}, i1.Get(m1, m2, m3))
	assert.Equal(t, Entries{m1, m2, m3}, i2.Get(m1, m2, m3))

	m4 := mockEntry(5)
	m5 := mockEntry(1)

	i3, overwritten := i2.Insert(m4, m5)
	assert.Equal(t, Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, uint64(5), i3.Len())
	assert.Equal(t, Entries{nil, nil}, i2.Get(m4, m5))
	assert.Equal(t, Entries{m4, m5}, i3.Get(m4, m5))
}

func TestAVLInsertLeftLeaningDoubleRotation(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(20)
	m2 := mockEntry(10)
	m3 := mockEntry(15)

	i2, overwritten := i1.Insert(m1, m2, m3)
	assert.Equal(t, Entries{nil, nil, nil}, overwritten)
	assert.Equal(t, uint64(0), i1.Len())
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, Entries{nil, nil, nil}, i1.Get(m1, m2, m3))
	assert.Equal(t, Entries{m1, m2, m3}, i2.Get(m1, m2, m3))

	m4 := mockEntry(1)
	m5 := mockEntry(5)

	i3, overwritten := i2.Insert(m4, m5)
	assert.Equal(t, Entries{nil, nil}, overwritten)
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, uint64(5), i3.Len())
	assert.Equal(t, Entries{nil, nil}, i2.Get(m4, m5))
	assert.Equal(t, Entries{m4, m5}, i3.Get(m4, m5))
	assert.Equal(t, Entries{m1, m2, m3}, i3.Get(m1, m2, m3))
}

func TestAVLInsertOverwrite(t *testing.T) {
	i1 := NewImmutable()
	m1 := mockEntry(20)
	m2 := mockEntry(10)
	m3 := mockEntry(15)

	i2, _ := i1.Insert(m1, m2, m3)
	m4 := mockEntry(15)

	i3, overwritten := i2.Insert(m4)
	assert.Equal(t, Entries{m3}, overwritten)
	assert.Equal(t, uint64(3), i2.Len())
	assert.Equal(t, uint64(3), i3.Len())
	assert.Equal(t, Entries{m4}, i3.Get(m4))
	assert.Equal(t, Entries{m3}, i2.Get(m3))
}

func BenchmarkImmutableInsert(b *testing.B) {
	numItems := b.N
	sl := NewImmutable()

	entries := generateMockEntries(numItems)
	sl, _ = sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl, _ = sl.Insert(entries[i%numItems])
	}
}

func BenchmarkImmutableGet(b *testing.B) {
	numItems := b.N
	sl := NewImmutable()

	entries := generateMockEntries(numItems)
	sl, _ = sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems])
	}
}

func BenchmarkImmutableBulkInsert(b *testing.B) {
	numItems := b.N
	sl := NewImmutable()

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries...)
	}
}
