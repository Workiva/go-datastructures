package yfast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntriesInsert(t *testing.T) {
	es := Entries{}

	e1 := newMockEntry(5)
	e2 := newMockEntry(1)

	es.insert(e1)
	es.insert(e2)

	assert.Equal(t, Entries{e2, e1}, es)

	e3 := newMockEntry(3)
	es.insert(e3)

	assert.Equal(t, Entries{e2, e3, e1}, es)
}

func TestEntriesDelete(t *testing.T) {
	es := Entries{}

	e1 := newMockEntry(5)
	e2 := newMockEntry(1)
	es.insert(e1)
	es.insert(e2)

	es.delete(5)
	assert.Equal(t, Entries{e2}, es)

	es.delete(1)
	assert.Equal(t, Entries{}, es)
}

func TestEntriesMax(t *testing.T) {
	es := Entries{}
	max, ok := es.max()
	assert.Equal(t, uint64(0), max)
	assert.False(t, ok)

	e2 := newMockEntry(1)
	es.insert(e2)
	max, ok = es.max()
	assert.Equal(t, uint64(1), max)
	assert.True(t, ok)

	e1 := newMockEntry(5)
	es.insert(e1)
	max, ok = es.max()
	assert.Equal(t, uint64(5), max)
	assert.True(t, ok)
}

func TestEntriesGet(t *testing.T) {
	es := Entries{}

	e1 := newMockEntry(5)
	e2 := newMockEntry(1)
	es.insert(e1)
	es.insert(e2)

	result := es.get(5)
	assert.Equal(t, e1, result)

	result = es.get(1)
	assert.Equal(t, e2, result)

	result = es.get(10)
	assert.Nil(t, result)
}

func TestEntriesSuccessor(t *testing.T) {
	es := Entries{}

	successor, i := es.successor(5)
	assert.Equal(t, -1, i)
	assert.Nil(t, successor)

	e1 := newMockEntry(5)
	e2 := newMockEntry(1)
	es.insert(e1)
	es.insert(e2)

	successor, i = es.successor(0)
	assert.Equal(t, 0, i)
	assert.Equal(t, e2, successor)

	successor, i = es.successor(2)
	assert.Equal(t, 1, i)
	assert.Equal(t, e1, successor)

	successor, i = es.successor(5)
	assert.Equal(t, 1, i)
	assert.Equal(t, e1, successor)

	successor, i = es.successor(10)
	assert.Equal(t, -1, i)
	assert.Nil(t, successor)
}

func TestEntriesPredecessor(t *testing.T) {
	es := Entries{}

	predecessor, i := es.predecessor(5)
	assert.Equal(t, -1, i)
	assert.Nil(t, predecessor)

	e1 := newMockEntry(5)
	e2 := newMockEntry(1)
	es.insert(e1)
	es.insert(e2)

	predecessor, i = es.predecessor(0)
	assert.Equal(t, -1, i)
	assert.Nil(t, predecessor)

	predecessor, i = es.predecessor(2)
	assert.Equal(t, 0, i)
	assert.Equal(t, e2, predecessor)

	predecessor, i = es.predecessor(5)
	assert.Equal(t, 1, i)
	assert.Equal(t, e1, predecessor)

	predecessor, i = es.predecessor(10)
	assert.Equal(t, 1, i)
	assert.Equal(t, e1, predecessor)
}
