package heap

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

	value, ok := h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 5, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 7, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 9, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 10, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 15, value)

	value, ok = h.Pop()
	assert.True(t, ok)
	assert.Equal(t, 20, value)

	_, ok = h.Pop()
	assert.False(t, ok)
}

func TestDaryHeapVariantsAndHeapify(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	for _, d := range []int{2, 3, 6, 8} {
		values := []int{10, 1, 5, 9, 7, 2, 15, 20, 3, 3, -1}
		h := NewDaryHeapFromSlice(d, values, cmp)
		peek, ok := h.Peek()
		assert.True(t, ok)
		assert.Equal(t, -1, peek)
		prev := -1 << 31
		for {
			v, ok := h.Pop()
			if !ok {
				break
			}
			assert.GreaterOrEqual(t, v, prev)
			prev = v
		}
	}
}
