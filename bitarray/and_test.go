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

// checkBit is a helper method for these unit tests
func checkBit(t *testing.T, ba BitArray, position uint64, expected bool) {
	ok, err := ba.GetBit(position)
	if assert.NoError(t, err) {
		if expected {
			assert.True(t, ok, "Bitarray at position %d should be set", position)
		} else {
			assert.False(t, ok, "Bitarray at position %d should be unset", position)
		}
	}
}

func TestAndSparseWithSparseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	// bits for which only one of the arrays is set
	sba.SetBit(3)
	sba.SetBit(280)
	other.SetBit(9)
	other.SetBit(100)

	// bits for which both arrays are set
	sba.SetBit(1)
	other.SetBit(1)
	sba.SetBit(2680)
	other.SetBit(2680)
	sba.SetBit(30)
	other.SetBit(30)

	ba := andSparseWithSparseBitArray(sba, other)

	checkBit(t, ba, 1, true)
	checkBit(t, ba, 30, true)
	checkBit(t, ba, 2680, true)

	checkBit(t, ba, 3, false)
	checkBit(t, ba, 9, false)
	checkBit(t, ba, 100, false)
	checkBit(t, ba, 2, false)
	checkBit(t, ba, 280, false)
}

func TestAndSparseWithDenseBitArray(t *testing.T) {
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

	ba := andSparseWithDenseBitArray(sba, other)

	checkBit(t, ba, 1, true)
	checkBit(t, ba, 150, true)
	checkBit(t, ba, 155, false)
	checkBit(t, ba, 156, false)
	checkBit(t, ba, 300, true)
}

// Maks sure that the sparse array is trimmed correctly if compared against a
// smaller dense bit array.
func TestAndSparseWithSmallerDenseBitArray(t *testing.T) {
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

	ba := andSparseWithDenseBitArray(sba, other)

	checkBit(t, ba, 1, true)
	checkBit(t, ba, 150, true)
	checkBit(t, ba, 155, false)
	checkBit(t, ba, 128, false)
	checkBit(t, ba, 500, false)
	checkBit(t, ba, 1200, false)
	checkBit(t, ba, 1500, false)
}

func TestAndDenseWithDenseBitArray(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(2000)

	dba.SetBit(1)
	other.SetBit(18)
	dba.SetBit(222)
	other.SetBit(222)
	other.SetBit(1501)

	ba := andDenseWithDenseBitArray(dba, other)

	checkBit(t, ba, 0, false)
	checkBit(t, ba, 1, false)
	checkBit(t, ba, 3, false)
	checkBit(t, ba, 18, false)
	checkBit(t, ba, 222, true)

	// check that the ba is the minimum of the size of `dba` and `other`
	// (dense bitarrays return an error on an out-of-bounds access)
	_, err := ba.GetBit(1500)
	assert.Equal(t, OutOfRangeError(1500), err)
	_, err = ba.GetBit(1501)
	assert.Equal(t, OutOfRangeError(1501), err)
}

func TestAndSparseWithEmptySparse(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	sba.SetBit(5)

	ba := andSparseWithSparseBitArray(sba, other)
	checkBit(t, ba, 0, false)
	checkBit(t, ba, 5, false)
	checkBit(t, ba, 100, false)
}

func TestAndSparseWithEmptyDense(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(1000)

	sba.SetBit(5)
	ba := andSparseWithDenseBitArray(sba, other)
	checkBit(t, ba, 5, false)

	sba.Reset()
	other.SetBit(5)

	ba = andSparseWithDenseBitArray(sba, other)
	checkBit(t, ba, 5, false)
}

func TestAndDenseWithEmptyDense(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(1000)

	dba.SetBit(5)
	ba := andDenseWithDenseBitArray(dba, other)
	checkBit(t, ba, 5, false)

	dba.Reset()
	other.SetBit(5)
	ba = andDenseWithDenseBitArray(dba, other)
	checkBit(t, ba, 5, false)
}
