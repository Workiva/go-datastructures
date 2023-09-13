package heaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDaryHeap(t *testing.T) {
	h := NewDaryHeap(6, func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})

	h.Push(10)
	h.Push(1)
	h.Push(5)
	h.Push(9)
	h.Push(7)
	h.Push(2)
	h.Push(15)

	h.Push(20)

	ok, value := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 5, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 7, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 9, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 10, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 15, value)

	ok, value = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 20, value)

	ok, _ = h.Pop()
	assert.False(t, ok)
}

func BenchmarkDaryHeap(b *testing.B) {
	h := NewDaryHeap(3, func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	for i := 0; i < b.N; i++ {
		h.Push(b.N)
	}
	for i := 0; i < b.N; i++ {
		h.Pop()
	}
}
