package yfast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRTAddSingleDimension(t *testing.T) {
	rt := new(1, uint8(0))
	e1 := newMockEntry(2)
	e2 := newMockEntry(5)

	overwritten := rt.Add(e1, e2)
	assert.Len(t, overwritten, 2)
	assert.Equal(t, Entries{nil, nil}, overwritten)

	assert.Equal(t, Entries{e1}, rt.Get(newMockEntry(2)))
	assert.Equal(t, Entries{e2}, rt.Get(newMockEntry(5)))
	assert.Equal(t, Entries{e1, e2}, rt.Get(newMockEntry(2), newMockEntry(5)))

	assert.Equal(t, Entries{nil, nil}, rt.Get(newMockEntry(18), newMockEntry(19)))
	assert.Equal(t, Entries{e1, nil}, rt.Get(newMockEntry(2), newMockEntry(3)))
}

func TestRTAddSingleDimensionOverwrite(t *testing.T) {
	rt := new(1, uint8(0))
	e1 := newMockEntry(2)
	e2 := newMockEntry(2)

	rt.Add(e1)
	overwritten := rt.Add(e2)

	assert.Equal(t, Entries{e1}, overwritten)
	assert.Equal(t, Entries{e2}, rt.Get(newMockEntry(2)))
}

func TestRTAddMultiDimension(t *testing.T) {
	rt := new(2, uint8(0))

	e1 := newMockEntry(2, 3)
	e2 := newMockEntry(3, 4)

	overwritten := rt.Add(e1, e2)
	assert.Len(t, overwritten, 2)
	assert.Equal(t, Entries{nil, nil}, overwritten)
}
