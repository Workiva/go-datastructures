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

/*
A generic implementation of a d-ary heap.

The d-ary heap or d-heap is a priority queue data structure, a generalization
of the binary heap in which the nodes have d children instead of 2.
*/
package heaps

import "math"

type DaryHeap[T any] struct {
	d       int
	data    []T
	compare func(T, T) int
}

func NewDaryHeap[T any](d int, compare func(T, T) int) *DaryHeap[T] {
	return &DaryHeap[T]{
		d:       d,
		data:    make([]T, 0),
		compare: compare,
	}
}

func (h *DaryHeap[T]) Len() int {
	return len(h.data)
}

func (h *DaryHeap[T]) Push(value T) {
	h.data = append(h.data, value)
	h.bubbleUp(h.Len() - 1)
}

func (h *DaryHeap[T]) Pop() (ok bool, value T) {
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
func (h *DaryHeap[T]) bubbleUp(index int) {
	if index == 0 {
		return
	}
	var parent = (index - 1) / h.d // Todo: make test fail if d is not 2 but you divide by 2
	if h.compare(h.data[index], h.data[parent]) < 0 {
		h.swap(index, parent)
		h.bubbleUp(parent)
	}
}

// Min heap: if a node is greater than its children, swap the node with the smallest child.
func (h *DaryHeap[T]) sinkDown(index int) {
	var childrenIndex = int(math.Pow(2, float64(index)))
	var smallest = index

	for i := 0; i < h.d; i++ {
		var child = childrenIndex + i
		if child >= h.Len() {
			break
		}
		if h.compare(h.data[child], h.data[smallest]) < 0 {
			smallest = child
		}
	}
	if smallest != index {
		h.swap(index, smallest)
		h.sinkDown(smallest)
	}
}

func (h *DaryHeap[T]) swap(i, j int) {
	var tmp = h.data[i]
	h.data[i] = h.data[j]
	h.data[j] = tmp
}
