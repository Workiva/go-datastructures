package xfast

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/slice"
)

func TestInsert(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	assert.True(t, xft.Exists(5))

	key := uint64(math.MaxUint64 - 235325)
	e2 := newMockEntry(key)
	xft.Insert(e2)

	assert.True(t, xft.Exists(key))
	assert.Equal(t, uint64(2), xft.Len())
}

func TestSuccessorDoesNotExist(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Successor(6)
	assert.Nil(t, result)
}

func TestSuccessorIsExactValue(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorGreaterThanKey(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(math.MaxUint64)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorCloseToKey(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorBetweenTwoKeys(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	for i := uint64(11); i < 20; i++ {
		result := xft.Successor(i)
		assert.Equal(t, e2, result)
	}

	for i := uint64(21); i < 100; i++ {
		result := xft.Successor(i)
		assert.Nil(t, result)
	}
}

func TestPredecessorDoesNotExist(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(4)
	assert.Nil(t, result)
}

func TestPredecessorIsExactValue(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(5)
	assert.Equal(t, e1, result)
}

func TestPredecessorLessThanKey(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(0)
	xft.Insert(e1)

	result := xft.Predecessor(math.MaxUint64)
	assert.Equal(t, e1, result)
}

func TestPredecessorCloseToKey(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(10)
	assert.Equal(t, e1, result)
}

func TestPredecessorBetweenTwoKeys(t *testing.T) {
	xft := New(uint64(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	for i := uint64(16); i < 17; i++ {
		result := xft.Predecessor(i)
		assert.Equal(t, e1, result)
	}

	for i := uint64(0); i < 10; i++ {
		result := xft.Predecessor(i)
		assert.Nil(t, result)
	}
}

func BenchmarkSuccessor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xft := New(uint64(0))
		e := newMockEntry(uint64(i))
		xft.Insert(e)
		xft.Successor(0)
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xft := New(uint64(0))
		e := newMockEntry(uint64(i))
		xft.Insert(e)
	}
}

// benchmarked against a flat list
func BenchmarkListInsert(b *testing.B) {
	numItems := 100000

	s := make(slice.Int64Slice, 0, numItems)
	for j := int64(0); j < int64(numItems); j++ {
		s = append(s, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Insert(int64(i))
	}
}

func BenchmarkListSearch(b *testing.B) {
	numItems := 1000000

	s := make(slice.Int64Slice, 0, numItems)
	for j := int64(0); j < int64(numItems); j++ {
		s = append(s, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Search(int64(i))
	}
}
