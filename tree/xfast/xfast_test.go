package xfast

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	xft := New()
	e1 := newMockEntry(5)
	xft.Insert(e1)

	assert.True(t, xft.Exists(5))

	key := uint64(math.MaxUint64 - 235325)
	e2 := newMockEntry(key)
	xft.Insert(e2)

	assert.True(t, xft.Exists(key))
	assert.Equal(t, uint64(2), xft.Len())
}
