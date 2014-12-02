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

func TestBlockToNums(t *testing.T) {
	b := block(0)

	b = b.insert(s - 2)
	b = b.insert(s - 6)

	expected := []uint64{s - 6, s - 2}

	result := make([]uint64, 0, 0)
	b.toNums(0, &result)
	assert.Equal(t, expected, result)
}

func BenchmarkBlockToNums(b *testing.B) {
	block := block(0)
	for i := uint64(0); i < s; i++ {
		block = block.insert(i)
	}

	nums := make([]uint64, 0, 0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		block.toNums(0, &nums)
	}
}
