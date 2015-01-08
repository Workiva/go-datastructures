package fastinteger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)

	assert.True(t, hm.Exists(5))
	assert.Equal(t, 5, hm.Get(5))
}

func TestInsertOverwrite(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)
	hm.Set(5, 10)

	assert.True(t, hm.Exists(5))
	assert.Equal(t, 10, hm.Get(5))
}

func TestMultipleInserts(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)
	hm.Set(6, 6)

	assert.True(t, hm.Exists(6))
	assert.Equal(t, 6, hm.Get(6))
}

func TestRebuild(t *testing.T) {
	numItems := uint64(100)

	hm := New(10)

	for i := uint64(0); i < numItems; i++ {
		hm.Set(i, i)
	}

	for i := uint64(0); i < numItems; i++ {
		assert.Equal(t, i, hm.Get(i))
	}
}
