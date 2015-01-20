package skip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntriesInsert(t *testing.T) {
	e1 := newMockEntry(10)
	e2 := newMockEntry(5)

	entries := Entries{e1}
	result := entries.insert(e2)
	assert.Equal(t, Entries{e2, e1}, entries)
	assert.Nil(t, result)

	e3 := newMockEntry(15)
	result = entries.insert(e3)
	assert.Equal(t, Entries{e2, e1, e3}, entries)
	assert.Nil(t, result)
}

func TestInsertOverwrite(t *testing.T) {
	e1 := newMockEntry(10)
	e2 := newMockEntry(10)

	entries := Entries{e1}
	result := entries.insert(e2)
	assert.Equal(t, e1, result)
}

func TestEntriesDelete(t *testing.T) {
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)

	entries := Entries{e1, e2}

	result := entries.delete(10)
	assert.Equal(t, Entries{e1}, entries)
	assert.Equal(t, e2, result)

	result = entries.delete(5)
	assert.Equal(t, Entries{}, entries)
	assert.Equal(t, e1, result)
}
