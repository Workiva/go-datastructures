package bitarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCompressedBit(t *testing.T) {
	ba := newSparseBitArray()

	result, err := ba.GetBit(5)
	assert.Nil(t, err)
	assert.False(t, result)
}

func BenchmarkGetCompressedBit(b *testing.B) {
	numItems := 1000
	ba := newSparseBitArray()

	for i := 0; i < numItems; i++ {
		ba.SetBit(uint64(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.GetBit(s)
	}
}

func TestGetSetCompressedBit(t *testing.T) {
	ba := newSparseBitArray()

	ba.SetBit(5)

	result, err := ba.GetBit(5)
	assert.Nil(t, err)
	assert.True(t, result)
	result, err = ba.GetBit(7)
	assert.Nil(t, err)
	assert.False(t, result)

	ba.SetBit(s * 2)
	result, _ = ba.GetBit(s * 2)
	assert.True(t, result)
	result, _ = ba.GetBit(s*2 + 1)
	assert.False(t, result)
}

func BenchmarkSetCompressedBit(b *testing.B) {
	numItems := 1000
	ba := newSparseBitArray()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < numItems; j++ {
			ba.SetBit(uint64(j))
		}
	}
}

func TestClearCompressedBit(t *testing.T) {
	ba := newSparseBitArray()
	ba.SetBit(5)
	ba.ClearBit(5)

	result, err := ba.GetBit(5)
	assert.Nil(t, err)
	assert.False(t, result)
	assert.Len(t, ba.blocks, 0)
	assert.Len(t, ba.indices, 0)

	ba.SetBit(s * 2)
	ba.ClearBit(s * 2)

	result, _ = ba.GetBit(s * 2)
	assert.False(t, result)
	assert.Len(t, ba.indices, 0)
	assert.Len(t, ba.blocks, 0)
}

func BenchmarkClearCompressedBit(b *testing.B) {
	numItems := 1000
	ba := newSparseBitArray()
	for i := 0; i < numItems; i++ {
		ba.SetBit(uint64(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.ClearBit(uint64(i))
	}
}

func TestClearCompressedBitArray(t *testing.T) {
	ba := newSparseBitArray()
	ba.SetBit(5)
	ba.SetBit(s * 2)

	result, err := ba.GetBit(5)
	assert.Nil(t, err)
	assert.True(t, result)
	result, _ = ba.GetBit(s * 2)
	assert.True(t, result)

	ba.Reset()

	result, err = ba.GetBit(5)
	assert.Nil(t, err)
	assert.False(t, result)
	result, _ = ba.GetBit(s * 2)
	assert.False(t, result)
}

func TestCompressedEquals(t *testing.T) {
	ba := newSparseBitArray()
	other := newSparseBitArray()

	assert.True(t, ba.Equals(other))

	ba.SetBit(5)
	assert.False(t, ba.Equals(other))

	other.SetBit(5)
	assert.True(t, ba.Equals(other))

	ba.ClearBit(5)
	assert.False(t, ba.Equals(other))
}

func TestCompressedIntersects(t *testing.T) {
	ba := newSparseBitArray()
	other := newSparseBitArray()

	assert.True(t, ba.Intersects(other))

	other.SetBit(5)

	assert.False(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))

	ba.SetBit(5)

	assert.True(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))

	other.SetBit(10)

	assert.False(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))
}

func TestLongCompressedIntersects(t *testing.T) {
	ba := newSparseBitArray()
	other := newSparseBitArray()

	ba.SetBit(5)
	other.SetBit(5)

	assert.True(t, ba.Intersects(other))

	other.SetBit(s * 2)

	assert.False(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))

	ba.SetBit(s * 2)

	assert.True(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))

	other.SetBit(s*2 + 1)

	assert.False(t, ba.Intersects(other))
	assert.True(t, other.Intersects(ba))
}

func BenchmarkCompressedIntersects(b *testing.B) {
	numItems := uint64(1000)

	ba := newSparseBitArray()
	other := newSparseBitArray()

	for i := uint64(0); i < numItems; i++ {
		ba.SetBit(i)
		other.SetBit(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.Intersects(other)
	}
}

func TestSparseIntersectsBitArray(t *testing.T) {
	cba := newSparseBitArray()
	ba := newBitArray(s * 2)

	assert.True(t, cba.Intersects(ba))
	ba.SetBit(5)

	assert.False(t, cba.Intersects(ba))
	cba.SetBit(5)

	assert.True(t, cba.Intersects(ba))
	cba.SetBit(10)

	assert.True(t, cba.Intersects(ba))
	ba.SetBit(s + 1)

	assert.False(t, cba.Intersects(ba))
	cba.SetBit(s + 1)

	assert.True(t, cba.Intersects(ba))
	cba.SetBit(s * 3)

	assert.True(t, cba.Intersects(ba))
}

func TestSparseEqualsBitArray(t *testing.T) {
	cba := newSparseBitArray()
	ba := newBitArray(s * 2)

	assert.True(t, cba.Equals(ba))

	ba.SetBit(5)
	assert.False(t, cba.Equals(ba))

	cba.SetBit(5)
	assert.True(t, cba.Equals(ba))

	ba.SetBit(s + 1)
	assert.False(t, cba.Equals(ba))

	cba.SetBit(s + 1)
	assert.True(t, cba.Equals(ba))
}

func BenchmarkCompressedEquals(b *testing.B) {
	numItems := uint64(1000)
	cba := newSparseBitArray()
	other := newSparseBitArray()

	for i := uint64(0); i < numItems; i++ {
		cba.SetBit(i)
		other.SetBit(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cba.Equals(other)
	}
}

func TestInsertPreviousBlockInSparse(t *testing.T) {
	sba := newSparseBitArray()

	sba.SetBit(s * 2)
	sba.SetBit(s - 1)

	result, err := sba.GetBit(s - 1)
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestSparseBitArrayToNums(t *testing.T) {
	sba := newSparseBitArray()

	sba.SetBit(s - 1)
	sba.SetBit(s + 1)

	expected := []uint64{s - 1, s + 1}

	results := sba.ToNums()
	assert.Equal(t, expected, results)
}

func BenchmarkSparseBitArrayToNums(b *testing.B) {
	numItems := uint64(1000)
	sba := newSparseBitArray()

	for i := uint64(0); i < numItems; i++ {
		sba.SetBit(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sba.ToNums()
	}
}
