package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	s := Int64Slice{3, 6, 1, 0, -1}
	s.Sort()

	assert.Equal(t, Int64Slice{-1, 0, 1, 3, 6}, s)
}

func TestSearch(t *testing.T) {
	s := Int64Slice{1, 3, 6}

	assert.Equal(t, 1, s.Search(3))
	assert.Equal(t, 1, s.Search(2))
	assert.Equal(t, 3, s.Search(7))
}

func TestExists(t *testing.T) {
	s := Int64Slice{1, 3, 6}

	assert.True(t, s.Exists(3))
	assert.False(t, s.Exists(4))
}

func TestInsert(t *testing.T) {
	s := Int64Slice{1, 3, 6}
	s = s.Insert(2)
	assert.Equal(t, Int64Slice{1, 2, 3, 6}, s)

	s = s.Insert(7)
	assert.Equal(t, Int64Slice{1, 2, 3, 6, 7}, s)
}
