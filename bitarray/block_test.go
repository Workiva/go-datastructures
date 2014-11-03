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
