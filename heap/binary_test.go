package heap

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {
	h := NewHeap(func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})

	h.Push(10)
	h.Push(15)
	h.Push(1)
	h.Push(5)
	h.Push(9)
	h.Push(7)
	h.Push(2)

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

	_, ok = h.Pop()
	assert.False(t, ok)
}

func BenchmarkBinaryVsDary2(b *testing.B) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	// Benchmark pushes then pops b.N elements for both heaps
	b.Run("binary", func(b *testing.B) {
		h := NewHeap[int](cmp)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Push(i)
		}
		for i := 0; i < b.N; i++ {
			h.Pop()
		}
	})
	b.Run("dary-2", func(b *testing.B) {
		h := NewDaryHeap[int](2, cmp)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Push(i)
		}
		for i := 0; i < b.N; i++ {
			h.Pop()
		}
	})
}

func TestHeapifyAndPeek(t *testing.T) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	values := []int{5, 3, 9, 1, 4, 8, 2}
	h := NewHeapFromSlice(values, cmp)
	peek, ok := h.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, peek)
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
func BenchmarkBinaryVsDary2_Mixed(b *testing.B) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	// 50/50 push/pop random workload.
	// Track heapSize instead of calling Len() in the hot loop to avoid lock overhead
	// and to prevent Pop on an empty heap.
	b.Run("binary-mixed-50-50", func(b *testing.B) {
		rng := rand.New(rand.NewSource(1))
		h := NewHeap[int](cmp)
		heapSize := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if rng.Intn(2) == 0 {
				h.Push(i)
				heapSize++
			} else if heapSize > 0 {
				_, ok := h.Pop()
				if ok {
					heapSize--
				}
			} else {
				h.Push(i)
				heapSize++
			}
		}
	})
	b.Run("dary2-mixed-50-50", func(b *testing.B) {
		rng := rand.New(rand.NewSource(1))
		h := NewDaryHeap[int](2, cmp)
		heapSize := 0
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if rng.Intn(2) == 0 {
				h.Push(i)
				heapSize++
			} else if heapSize > 0 {
				_, ok := h.Pop()
				if ok {
					heapSize--
				}
			} else {
				h.Push(i)
				heapSize++
			}
		}
	})
}

func BenchmarkBinaryVsDary2_Heapify(b *testing.B) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	// Build from slice then pop all
	const size = 10000
	values := make([]int, size)
	for i := 0; i < size; i++ {
		values[i] = size - i
	}
	b.Run("binary-heapify-then-pop", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := NewHeapFromSlice(values, cmp)
			for j := 0; j < size; j++ {
				h.Pop()
			}
		}
	})
	b.Run("dary2-heapify-then-pop", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := NewDaryHeapFromSlice(2, values, cmp)
			for j := 0; j < size; j++ {
				h.Pop()
			}
		}
	})
}

func BenchmarkBinaryVsDary2_Bursts(b *testing.B) {
	cmp := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}
	const burst = 64
	b.Run("binary-bursts", func(b *testing.B) {
		h := NewHeap[int](cmp)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i += burst {
			for j := 0; j < burst; j++ {
				h.Push(i + j)
			}
			for j := 0; j < burst; j++ {
				h.Pop()
			}
		}
	})
	b.Run("dary2-bursts", func(b *testing.B) {
		h := NewDaryHeap[int](2, cmp)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i += burst {
			for j := 0; j < burst; j++ {
				h.Push(i + j)
			}
			for j := 0; j < burst; j++ {
				h.Pop()
			}
		}
	})
}
