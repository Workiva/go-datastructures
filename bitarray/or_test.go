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

func TestOrSparseWithSparseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	ctx := false
	for i := uint64(0); i < 1000; i += s {
		if ctx {
			sba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	sba.SetBit(s - 1)
	other.SetBit(s - 1)

	result := orSparseWithSparseBitArray(sba, other)

	for i := uint64(0); i < 1000; i += s {
		ok, err := result.GetBit(i)
		assert.Nil(t, err)
		assert.True(t, ok)
	}

	ok, err := result.GetBit(s - 1)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = result.GetBit(s - 2)
	assert.Nil(t, err)
	assert.False(t, ok)

	other.SetBit(2000)
	result = orSparseWithSparseBitArray(sba, other)

	ok, err = result.GetBit(2000)
	assert.Nil(t, err)
	assert.True(t, ok)

	sba.SetBit(2000)
	result = orSparseWithSparseBitArray(sba, other)

	ok, err = result.GetBit(2000)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func BenchmarkOrSparseWithSparse(b *testing.B) {
	numItems := uint64(160000)
	sba := newSparseBitArray()
	other := newSparseBitArray()

	ctx := false
	for i := uint64(0); i < numItems; i += s {
		if ctx {
			sba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		orSparseWithSparseBitArray(sba, other)
	}
}

func TestOrSparseWithDenseBitArray(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(2000)

	ctx := false
	for i := uint64(0); i < 1000; i += s {
		if ctx {
			sba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	other.SetBit(1500)
	other.SetBit(s - 1)
	sba.SetBit(s - 1)

	result := orSparseWithDenseBitArray(sba, other)

	for i := uint64(0); i < 1000; i += s {
		ok, err := result.GetBit(i)
		assert.Nil(t, err)
		assert.True(t, ok)
	}

	ok, err := result.GetBit(1500)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = result.GetBit(s - 1)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = result.GetBit(s - 2)
	assert.Nil(t, err)
	assert.False(t, ok)

	sba.SetBit(2500)
	result = orSparseWithDenseBitArray(sba, other)

	ok, err = result.GetBit(2500)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func BenchmarkOrSparseWithDense(b *testing.B) {
	numItems := uint64(160000)
	sba := newSparseBitArray()
	other := newBitArray(numItems)

	ctx := false
	for i := uint64(0); i < numItems; i += s {
		if ctx {
			sba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		orSparseWithDenseBitArray(sba, other)
	}
}

func TestOrDenseWithDenseBitArray(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(2000)

	ctx := false
	for i := uint64(0); i < 1000; i += s {
		if ctx {
			dba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	other.SetBit(1500)
	other.SetBit(s - 1)
	dba.SetBit(s - 1)

	result := orDenseWithDenseBitArray(dba, other)

	for i := uint64(0); i < 1000; i += s {
		ok, err := result.GetBit(i)
		assert.Nil(t, err)
		assert.True(t, ok)
	}

	ok, err := result.GetBit(s - 1)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = result.GetBit(1500)
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = result.GetBit(1700)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func BenchmarkOrDenseWithDense(b *testing.B) {
	numItems := uint64(160000)
	dba := newBitArray(numItems)
	other := newBitArray(numItems)

	ctx := false
	for i := uint64(0); i < numItems; i += s {
		if ctx {
			dba.SetBit(i)
		} else {
			other.SetBit(i)
		}

		ctx = !ctx
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		orDenseWithDenseBitArray(dba, other)
	}
}

func TestOrSparseWithEmptySparse(t *testing.T) {
	sba := newSparseBitArray()
	other := newSparseBitArray()

	sba.SetBit(5)

	result := orSparseWithSparseBitArray(sba, other)
	assert.Equal(t, sba, result)

	sba.Reset()
	other.SetBit(5)

	result = orSparseWithSparseBitArray(sba, other)
	assert.Equal(t, other, result)
}

func TestOrSparseWithEmptyDense(t *testing.T) {
	sba := newSparseBitArray()
	other := newBitArray(1000)

	sba.SetBit(5)
	result := orSparseWithDenseBitArray(sba, other)
	assert.Equal(t, sba, result)

	sba.Reset()
	other.SetBit(5)

	result = orSparseWithDenseBitArray(sba, other)
	assert.Equal(t, other, result)
}

func TestOrDenseWithEmptyDense(t *testing.T) {
	dba := newBitArray(1000)
	other := newBitArray(1000)

	dba.SetBit(5)
	result := orDenseWithDenseBitArray(dba, other)
	assert.Equal(t, dba, result)

	dba.Reset()
	other.SetBit(5)
	result = orDenseWithDenseBitArray(dba, other)
	assert.Equal(t, other, result)
}
