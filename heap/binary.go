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
package heaps

type Heap[T any] struct {
	data    []T
	compare func(T, T) int
}

func NewHeap[T any](compare func(T, T) int) *Heap[T] {
	return &Heap[T]{
		data:    make([]T, 0),
		compare: compare,
	}
}

func (h *Heap[T]) Len() int {
	return len(h.data)
}

func (h *Heap[T]) Push(value T) {
	h.data = append(h.data, value)
	h.bubbleUp(h.Len() - 1)
}

func (h *Heap[T]) Pop() (ok bool, value T) {
	if h.Len() == 0 {
		return false, value
	}
	var top = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.sinkDown(0)
	return true, top
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
	var left = index*2 + 1
	var right = index*2 + 2
	var smallest = index
	if left < h.Len() && h.compare(h.data[left], h.data[smallest]) < 0 {
		smallest = left
	}
	if right < h.Len() && h.compare(h.data[right], h.data[smallest]) < 0 {
		smallest = right
	}
	if smallest != index {
		h.swap(index, smallest)
		h.sinkDown(smallest)
	}
}

func (h *Heap[T]) swap(i, j int) {
	var tmp = h.data[i]
	h.data[i] = h.data[j]
	h.data[j] = tmp
}
