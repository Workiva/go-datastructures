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
package hilbert

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHilbert(t *testing.T) {
	h := Encode(0, 0)
	x, y := Decode(h)
	assert.Equal(t, int64(0), h)
	assert.Equal(t, int32(0), x)
	assert.Equal(t, int32(0), y)

	h = Encode(1, 0)
	x, y = Decode(h)
	assert.Equal(t, int64(3), h)
	assert.Equal(t, int32(1), x)
	assert.Equal(t, int32(0), y)

	h = Encode(1, 1)
	x, y = Decode(h)
	assert.Equal(t, int64(2), h)
	assert.Equal(t, int32(1), x)
	assert.Equal(t, int32(1), y)

	h = Encode(0, 1)
	x, y = Decode(h)
	assert.Equal(t, int64(1), h)
	assert.Equal(t, int32(0), x)
	assert.Equal(t, int32(1), y)
}

func TestHilbertAtMaxRange(t *testing.T) {
	x, y := int32(math.MaxInt32), int32(math.MaxInt32)
	h := Encode(x, y)
	resultx, resulty := Decode(h)
	assert.Equal(t, x, resultx)
	assert.Equal(t, y, resulty)
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encode(int32(i), int32(i))
	}
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode(int64(i))
	}
}
