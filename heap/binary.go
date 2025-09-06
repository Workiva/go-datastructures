/*
MIT License

Copyright (c) 2021 Florimond Husquinet

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

/* A generic implementation of a binary heap */
package heap

import "sync"

type Heap[T any] struct {
	mu      sync.RWMutex
	data    []T
	compare func(T, T) int
}

// NewHeap constructs a heap using the provided comparator.
// The comparator should return -1 if a < b, 0 if a == b, and 1 if a > b.
// If compare orders values in ascending order, the heap behaves as a min-heap.
// To build a max-heap, invert the comparator (e.g., return -compare(a, b)).
func NewHeap[T any](compare func(T, T) int) *Heap[T] {
	return &Heap[T]{
		data:    make([]T, 0),
		compare: compare,
	}
}

// NewHeapFromSlice builds a heap in O(n) from an initial slice.
func NewHeapFromSlice[T any](values []T, compare func(T, T) int) *Heap[T] {
	h := &Heap[T]{
		data:    append([]T(nil), values...),
		compare: compare,
	}
	// heapify bottom-up
	for i := (len(h.data) / 2) - 1; i >= 0; i-- {
		h.sinkDown(i)
	}
	return h
}

// Peek returns the top element without removing it.
func (h *Heap[T]) Peek() (value T, ok bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.data) == 0 {
		return value, false
	}
	return h.data[0], true
}

func (h *Heap[T]) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}

func (h *Heap[T]) Push(value T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data = append(h.data, value)
	idx := len(h.data) - 1
	h.bubbleUp(idx)
}

func (h *Heap[T]) Pop() (value T, ok bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	n := len(h.data)
	if n == 0 {
		return value, false
	}
	top := h.data[0]
	h.data[0] = h.data[n-1]
	h.data = h.data[:n-1]
	h.sinkDown(0)
	return top, true
}

// Min heap: if a node is less than its parent, swap them.
func (h *Heap[T]) bubbleUp(index int) {
	if index == 0 {
		return
	}
	var parent = (index - 1) / 2
	if h.compare(h.data[index], h.data[parent]) < 0 {
		h.swap(index, parent)
		h.bubbleUp(parent)
	}
}

// Min heap: if a node is greater than its children, swap the node with the smallest child.
func (h *Heap[T]) sinkDown(index int) {
	n := len(h.data)
	left := index*2 + 1
	right := index*2 + 2
	smallest := index
	if left < n && h.compare(h.data[left], h.data[smallest]) < 0 {
		smallest = left
	}
	if right < n && h.compare(h.data[right], h.data[smallest]) < 0 {
		smallest = right
	}
	if smallest != index {
		h.swap(index, smallest)
		h.sinkDown(smallest)
	}
}

func (h *Heap[T]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}
