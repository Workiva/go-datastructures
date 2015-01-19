package skip

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Printf(`BOTCHED.`)
}

func generateMockEntries(num int) Entries {
	entries := make(Entries, 0, num)
	for i := uint64(0); i < uint64(num); i++ {
		entries = append(entries, newMockEntry(i))
	}

	return entries
}

func TestSimpleInsert(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1)
	assert.Equal(t, Entries{m1}, sl.Get(5))
	assert.Equal(t, uint64(1), sl.Len())
	assert.Equal(t, Entries{nil}, overwritten)

	overwritten = sl.Insert(m2)
	assert.Equal(t, Entries{m2}, sl.Get(6))
	assert.Equal(t, Entries{nil}, sl.Get(7))
	assert.Equal(t, uint64(2), sl.Len())
	assert.Equal(t, Entries{nil}, overwritten)
}

func TestSimpleOverwrite(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(5)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1)
	assert.Equal(t, Entries{nil}, overwritten)
	assert.Equal(t, uint64(1), sl.Len())

	overwritten = sl.Insert(m2)
	assert.Equal(t, Entries{m1}, overwritten)
	assert.Equal(t, uint64(1), sl.Len())
}

func TestInsertOutOfOrder(t *testing.T) {
	m1 := newMockEntry(6)
	m2 := newMockEntry(5)

	sl := New(uint8(0))

	overwritten := sl.Insert(m1, m2)
	assert.Equal(t, Entries{nil, nil}, overwritten)

	assert.Equal(t, Entries{m1, m2}, sl.Get(6, 5))
}

func TestDelete(t *testing.T) {
	m1 := newMockEntry(5)
	sl := New(uint8(0))
	sl.Insert(m1)

	deleted := sl.Delete(m1.Key())
	assert.Equal(t, Entries{m1}, deleted)
	assert.Equal(t, uint64(0), sl.Len())
	assert.Equal(t, Entries{nil}, sl.Get(5))
}

func TestDeleteAll(t *testing.T) {
	m1 := newMockEntry(5)
	m2 := newMockEntry(6)
	sl := New(uint8(0))
	sl.Insert(m1, m2)

	deleted := sl.Delete(m1.Key(), m2.Key())
	assert.Equal(t, Entries{m1, m2}, deleted)
	assert.Equal(t, uint64(0), sl.Len())
	assert.Equal(t, Entries{nil, nil}, sl.Get(m1.Key(), m2.Key()))
}

func BenchmarkInsert(b *testing.B) {
	numItems := 10000
	sl := New(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries[i%numItems])
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := 10000
	sl := New(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems].Key())
	}
}
