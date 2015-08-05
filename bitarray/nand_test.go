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
)

func TestNandSparseWithSparseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	// bits for which only one of the arrays is set
	sba.SetBit(3)
	sba.SetBit(280)
	other.SetBit(9)
	other.SetBit(100)
	sba.SetBit(1000)
	other.SetBit(1001)

	// bits for which both arrays are set
	sba.SetBit(1)
	other.SetBit(1)
	sba.SetBit(2680)
	other.SetBit(2680)
	sba.SetBit(30)
	other.SetBit(30)

	ba := nandSparseWithSparseBitArray(sba, other)

	// Bits in both
	checkBit(t, ba, 1, false)
	checkBit(t, ba, 30, false)
	checkBit(t, ba, 2680, false)

	// Bits in sba but not other
	checkBit(t, ba, 3, true)
	checkBit(t, ba, 280, true)
	checkBit(t, ba, 1000, true)

	// Bits in other but not sba
	checkBit(t, ba, 9, false)
	checkBit(t, ba, 100, false)
	checkBit(t, ba, 2, false)

	nums := ba.ToNums()
	assert.Equal(t, []uint64{3, 280, 1000}, nums)
}

func TestNandSparseWithDenseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(300)

	other.SetBit(1)
	sba.SetBit(1)
	other.SetBit(150)
	sba.SetBit(150)
	sba.SetBit(155)
	other.SetBit(156)
	sba.SetBit(300)
	other.SetBit(300)

	ba := nandSparseWithDenseBitArray(sba, other)

	// Bits in both
	checkBit(t, ba, 1, false)
	checkBit(t, ba, 150, false)
	checkBit(t, ba, 300, false)

	// Bits in sba but not other
	checkBit(t, ba, 155, true)

	// Bits in other but not sba
	checkBit(t, ba, 156, false)

	nums := ba.ToNums()
	assert.Equal(t, []uint64{155}, nums)
}

func TestNandDenseWithSparseBitArray(t *testing.T) {
	sba := newBitArray(300)
	other := newSparseBitArray()

	other.SetBit(1)
	sba.SetBit(1)
	other.SetBit(150)
	sba.SetBit(150)
	sba.SetBit(155)
	other.SetBit(156)
	sba.SetBit(300)
	other.SetBit(300)

	ba := nandDenseWithSparseBitArray(sba, other)

	// Bits in both
	checkBit(t, ba, 1, false)
	checkBit(t, ba, 150, false)
	checkBit(t, ba, 300, false)

	// Bits in sba but not other
	checkBit(t, ba, 155, true)

	// Bits in other but not sba
	checkBit(t, ba, 156, false)

	nums := ba.ToNums()
	assert.Equal(t, []uint64{155}, nums)
}

func TestNandSparseWithSmallerDenseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(512)

	other.SetBit(1)
	sba.SetBit(1)
	other.SetBit(150)
	sba.SetBit(150)
	sba.SetBit(155)
	sba.SetBit(500)

	other.SetBit(128)
	sba.SetBit(1500)
	sba.SetBit(1200)

	ba := nandSparseWithDenseBitArray(sba, other)

	// Bits in both
	checkBit(t, ba, 1, false)
	checkBit(t, ba, 150, false)

	// Bits in sba but not other
	checkBit(t, ba, 155, true)
	checkBit(t, ba, 500, true)
	checkBit(t, ba, 1200, true)
	checkBit(t, ba, 1500, true)

	// Bits in other but not sba
	checkBit(t, ba, 128, false)

	nums := ba.ToNums()
	assert.Equal(t, []uint64{155, 500, 1200, 1500}, nums)
}

func TestNandDenseWithDenseBitArray(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(2000)

	dba.SetBit(1)
	other.SetBit(18)
	dba.SetBit(222)
	other.SetBit(222)
	other.SetBit(1501)

	ba := nandDenseWithDenseBitArray(dba, other)

	// Bits in both
	checkBit(t, ba, 222, false)

	// Bits in dba and not other
	checkBit(t, ba, 1, true)

	// Bits in other
	checkBit(t, ba, 18, false)

	// Bits in neither
	checkBit(t, ba, 0, false)
	checkBit(t, ba, 3, false)

	// check that the ba is the minimum of the size of `dba` and `other`
	// (dense bitarrays return an error on an out-of-bounds access)
	_, err := ba.GetBit(1500)
	assert.Equal(t, OutOfRangeError(1500), err)
	_, err = ba.GetBit(1501)
	assert.Equal(t, OutOfRangeError(1501), err)

	nums := ba.ToNums()
	assert.Equal(t, []uint64{1}, nums)
}

func TestNandSparseWithEmptySparse(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	sba.SetBit(5)

	ba := nandSparseWithSparseBitArray(sba, other)

	checkBit(t, ba, 0, false)
	checkBit(t, ba, 5, true)
	checkBit(t, ba, 100, false)
}

func TestNandSparseWithEmptyDense(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(1000)

	sba.SetBit(5)
	ba := nandSparseWithDenseBitArray(sba, other)
	checkBit(t, ba, 5, true)

	sba.Reset()
	other.SetBit(5)

	ba = nandSparseWithDenseBitArray(sba, other)
	checkBit(t, ba, 5, false)
}

func TestNandDenseWithEmptyDense(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(1000)

	dba.SetBit(5)
	ba := nandDenseWithDenseBitArray(dba, other)
	checkBit(t, ba, 5, true)

	dba.Reset()
	other.SetBit(5)
	ba = nandDenseWithDenseBitArray(dba, other)
	checkBit(t, ba, 5, false)
}

func BenchmarkNandSparseWithSparse(b *testing.B) {
	numItems := uint64(160000)
	sba := newSparseBitArray()
	other := newSparseBitArray()

	for i := uint64(0); i < numItems; i += s {
		if i%200 == 0 {
			sba.SetBit(i)
		} else if i%300 == 0 {
			other.SetBit(i)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nandSparseWithSparseBitArray(sba, other)
	}
}

func BenchmarkNandSparseWithDense(b *testing.B) {
	numItems := uint64(160000)
	sba := newSparseBitArray()
	other := newBitArray(numItems)

	for i := uint64(0); i < numItems; i += s {
		if i%2 == 0 {
			sba.SetBit(i)
		} else if i%3 == 0 {
			other.SetBit(i)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nandSparseWithDenseBitArray(sba, other)
	}
}

func BenchmarkNandDenseWithSparse(b *testing.B) {
	numItems := uint64(160000)
	ba := newBitArray(numItems)
	other := newSparseBitArray()

	for i := uint64(0); i < numItems; i += s {
		if i%2 == 0 {
			ba.SetBit(i)
		} else if i%3 == 0 {
			other.SetBit(i)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nandDenseWithSparseBitArray(ba, other)
	}
}

func BenchmarkNandDenseWithDense(b *testing.B) {
	numItems := uint64(160000)
	dba := newBitArray(numItems)
	other := newBitArray(numItems)

	for i := uint64(0); i < numItems; i += s {
		if i%2 == 0 {
			dba.SetBit(i)
		} else if i%3 == 0 {
			other.SetBit(i)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		nandDenseWithDenseBitArray(dba, other)
	}
}
