/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bitarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBitOperations(t *testing.T) {
	ba := newBitArray(10)

	err := ba.SetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ba.GetBit(5)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Errorf(`Expected true at position: %d`, 5)
	}

	result, err = ba.GetBit(3)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Errorf(`Expected false at position %d`, 3)
	}

	err = ba.ClearBit(5)
	if err != nil {
		t.Fatal(err)
	}

	result, err = ba.GetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Errorf(`Expected false at position: %d`, 5)
	}

	ba = newBitArray(24)
	err = ba.SetBit(16)
	if err != nil {
		t.Fatal(err)
	}

	result, err = ba.GetBit(16)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Errorf(`Expected true at position: %d`, 16)
	}
}

func TestDuplicateOperation(t *testing.T) {
	ba := newBitArray(10)

	err := ba.SetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	err = ba.SetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ba.GetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	if !result {
		t.Errorf(`Expected true at position: %d`, 5)
	}

	err = ba.ClearBit(5)
	if err != nil {
		t.Fatal(err)
	}

	err = ba.ClearBit(5)
	if err != nil {
		t.Fatal(err)
	}

	result, err = ba.GetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Errorf(`Expected false at position: %d`, 5)
	}
}

func TestOutOfBounds(t *testing.T) {
	ba := newBitArray(4)

	err := ba.SetBit(s + 1)

	if _, ok := err.(OutOfRangeError); !ok {
		t.Errorf(`Expected out of range error.`)
	}

	_, err = ba.GetBit(s + 1)
	if _, ok := err.(OutOfRangeError); !ok {
		t.Errorf(`Expected out of range error.`)
	}
}

func TestIsEmpty(t *testing.T) {
	ba := newBitArray(10)
	assert.True(t, ba.IsEmpty())

	ba.SetBit(5)
	assert.False(t, ba.IsEmpty())
}

func TestCount(t *testing.T) {
	ba := newBitArray(500)
	assert.Equal(t, 0, ba.Count())

	require.NoError(t, ba.SetBit(0))
	assert.Equal(t, 1, ba.Count())

	require.NoError(t, ba.SetBit(40))
	require.NoError(t, ba.SetBit(64))
	require.NoError(t, ba.SetBit(100))
	require.NoError(t, ba.SetBit(200))
	require.NoError(t, ba.SetBit(469))
	require.NoError(t, ba.SetBit(500))
	assert.Equal(t, 7, ba.Count())

	require.NoError(t, ba.ClearBit(200))
	assert.Equal(t, 6, ba.Count())

	ba.Reset()
	assert.Equal(t, 0, ba.Count())
}

func TestClear(t *testing.T) {
	ba := newBitArray(10)

	err := ba.SetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	err = ba.SetBit(9)
	if err != nil {
		t.Fatal(err)
	}

	ba.Reset()

	assert.False(t, ba.anyset)
	result, err := ba.GetBit(5)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Errorf(`BA not reset.`)
	}

	result, err = ba.GetBit(9)
	if err != nil {
		t.Fatal(err)
	}

	if result {
		t.Errorf(`BA not reset.`)
	}
}

func BenchmarkGetBit(b *testing.B) {
	numItems := uint64(168000)

	ba := newBitArray(numItems)

	for i := uint64(0); i < numItems; i++ {
		ba.SetBit(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := uint64(0); j < numItems; j++ {
			ba.GetBit(j)
		}
	}
}

func TestGetSetBits(t *testing.T) {
	ba := newBitArray(1000)
	buf := make([]uint64, 0, 5)

	require.NoError(t, ba.SetBit(1))
	require.NoError(t, ba.SetBit(4))
	require.NoError(t, ba.SetBit(8))
	require.NoError(t, ba.SetBit(63))
	require.NoError(t, ba.SetBit(64))
	require.NoError(t, ba.SetBit(200))
	require.NoError(t, ba.SetBit(1000))

	assert.Equal(t, []uint64{1, 4, 8, 63, 64}, ba.GetSetBits(0, buf))
	assert.Equal(t, []uint64{63, 64, 200, 1000}, ba.GetSetBits(10, buf))
	assert.Equal(t, []uint64{63, 64, 200, 1000}, ba.GetSetBits(63, buf))
	assert.Equal(t, []uint64{200, 1000}, ba.GetSetBits(128, buf))

	require.NoError(t, ba.ClearBit(4))
	require.NoError(t, ba.ClearBit(64))
	assert.Equal(t, []uint64{1, 8, 63, 200, 1000}, ba.GetSetBits(0, buf))
	assert.Empty(t, ba.GetSetBits(1001, buf))

	ba.Reset()
	assert.Empty(t, ba.GetSetBits(0, buf))
}

func BenchmarkGetSetBits(b *testing.B) {
	numItems := uint64(168000)

	ba := newBitArray(numItems)
	for i := uint64(0); i < numItems; i++ {
		if i%13 == 0 || i%5 == 0 {
			require.NoError(b, ba.SetBit(i))
		}
	}

	buf := make([]uint64, 0, ba.Capacity())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.GetSetBits(0, buf)
	}
}

func TestEquality(t *testing.T) {
	ba := newBitArray(s + 1)
	other := newBitArray(s + 1)

	if !ba.Equals(other) {
		t.Errorf(`Expected equality.`)
	}

	ba.SetBit(s + 1)
	other.SetBit(s + 1)

	if !ba.Equals(other) {
		t.Errorf(`Expected equality.`)
	}

	other.SetBit(0)

	if ba.Equals(other) {
		t.Errorf(`Expected inequality.`)
	}
}

func BenchmarkEquality(b *testing.B) {
	ba := newBitArray(160000)
	other := newBitArray(ba.Capacity())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.Equals(other)
	}
}

func TestIntersects(t *testing.T) {
	ba := newBitArray(10)
	other := newBitArray(ba.Capacity())

	ba.SetBit(1)
	ba.SetBit(2)

	other.SetBit(1)

	if !ba.Intersects(other) {
		t.Errorf(`Is intersecting.`)
	}

	other.SetBit(5)

	if ba.Intersects(other) {
		t.Errorf(`Is not intersecting.`)
	}

	other = newBitArray(ba.Capacity() + 1)
	other.SetBit(1)

	if ba.Intersects(other) {
		t.Errorf(`Is not intersecting.`)
	}
}

func BenchmarkIntersects(b *testing.B) {
	ba := newBitArray(162432)
	other := newBitArray(ba.Capacity())

	ba.SetBit(159999)
	other.SetBit(159999)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.Intersects(other)
	}
}

func TestComplement(t *testing.T) {
	ba := newBitArray(10)

	ba.SetBit(5)

	ba.complement()

	if ok, _ := ba.GetBit(5); ok {
		t.Errorf(`Expected clear.`)
	}

	if ok, _ := ba.GetBit(4); !ok {
		t.Errorf(`Expected set.`)
	}
}

func BenchmarkComplement(b *testing.B) {
	ba := newBitArray(160000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.complement()
	}
}

func TestSetHighestLowest(t *testing.T) {
	ba := newBitArray(10)

	assert.False(t, ba.anyset)
	assert.Equal(t, uint64(0), ba.lowest)
	assert.Equal(t, uint64(0), ba.highest)

	ba.SetBit(5)

	assert.True(t, ba.anyset)
	assert.Equal(t, uint64(5), ba.lowest)
	assert.Equal(t, uint64(5), ba.highest)

	ba.SetBit(8)
	assert.Equal(t, uint64(5), ba.lowest)
	assert.Equal(t, uint64(8), ba.highest)
}

func TestGetBitAtCapacity(t *testing.T) {
	ba := newBitArray(s * 2)
	_, err := ba.GetBit(s * 2)
	assert.Error(t, err)
}

func TestSetBitAtCapacity(t *testing.T) {
	ba := newBitArray(s * 2)
	err := ba.SetBit(s * 2)
	assert.Error(t, err)
}

func TestClearBitAtCapacity(t *testing.T) {
	ba := newBitArray(s * 2)
	err := ba.ClearBit(s * 2)
	assert.Error(t, err)
}

func TestClearHighestLowest(t *testing.T) {
	ba := newBitArray(10)

	ba.SetBit(5)
	ba.ClearBit(5)

	assert.False(t, ba.anyset)
	assert.Equal(t, uint64(0), ba.lowest)
	assert.Equal(t, uint64(0), ba.highest)

	ba.SetBit(3)
	ba.SetBit(5)
	ba.SetBit(7)

	ba.ClearBit(7)
	assert.True(t, ba.anyset)
	assert.Equal(t, uint64(5), ba.highest)
	assert.Equal(t, uint64(3), ba.lowest)

	ba.SetBit(7)
	ba.ClearBit(3)
	assert.True(t, ba.anyset)
	assert.Equal(t, uint64(5), ba.lowest)
	assert.Equal(t, uint64(7), ba.highest)

	ba.ClearBit(7)
	assert.True(t, ba.anyset)
	assert.Equal(t, uint64(5), ba.lowest)
	assert.Equal(t, uint64(5), ba.highest)

	ba.ClearBit(5)
	assert.False(t, ba.anyset)
	assert.Equal(t, uint64(0), ba.lowest)
	assert.Equal(t, uint64(0), ba.highest)
}

func TestComplementResetsBounds(t *testing.T) {
	ba := newBitArray(5)

	ba.complement()
	assert.True(t, ba.anyset)
	assert.Equal(t, uint64(0), ba.lowest)
	assert.Equal(t, uint64(s-1), ba.highest)
}

func TestBitArrayIntersectsSparse(t *testing.T) {
	ba := newBitArray(s * 2)
	cba := newSparseBitArray()

	assert.True(t, ba.Intersects(cba))

	cba.SetBit(5)
	assert.False(t, ba.Intersects(cba))

	ba.SetBit(5)
	assert.True(t, ba.Intersects(cba))

	cba.SetBit(s + 1)
	assert.False(t, ba.Intersects(cba))

	ba.SetBit(s + 1)
	assert.True(t, ba.Intersects(cba))
}

func TestBitArrayEqualsSparse(t *testing.T) {
	ba := newBitArray(s * 2)
	cba := newSparseBitArray()

	assert.True(t, ba.Equals(cba))

	ba.SetBit(5)
	assert.False(t, ba.Equals(cba))

	cba.SetBit(5)
	assert.True(t, ba.Equals(cba))

	ba.SetBit(s + 1)
	assert.False(t, ba.Equals(cba))

	cba.SetBit(s + 1)
	assert.True(t, ba.Equals(cba))
}

func TestConstructorSetBitArray(t *testing.T) {
	ba := newBitArray(8, true)

	result, err := ba.GetBit(7)
	assert.Nil(t, err)
	assert.True(t, result)
	assert.Equal(t, s-1, ba.highest)
	assert.Equal(t, uint64(0), ba.lowest)
	assert.True(t, ba.anyset)
}

func TestCopyBitArray(t *testing.T) {
	ba := newBitArray(10)
	ba.SetBit(5)
	ba.SetBit(1)

	result := ba.copy().(*bitArray)
	assert.Equal(t, ba.anyset, result.anyset)
	assert.Equal(t, ba.lowest, result.lowest)
	assert.Equal(t, ba.highest, result.highest)
	assert.Equal(t, ba.blocks, result.blocks)
}

func BenchmarkDenseIntersectsCompressed(b *testing.B) {
	numBits := uint64(162432)
	ba := newBitArray(numBits)
	other := newSparseBitArray()

	for i := uint64(0); i < numBits; i++ {
		ba.SetBit(i)
		other.SetBit(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ba.intersectsSparseBitArray(other)
	}
}

func TestBitArrayToNums(t *testing.T) {
	ba := newBitArray(s * 2)

	ba.SetBit(s - 1)
	ba.SetBit(s + 1)

	expected := []uint64{s - 1, s + 1}

	result := ba.ToNums()

	assert.Equal(t, expected, result)
}

func BenchmarkBitArrayToNums(b *testing.B) {
	numItems := uint64(1000)
	ba := newBitArray(numItems)

	for i := uint64(0); i < numItems; i++ {
		ba.SetBit(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ba.ToNums()
	}
}
