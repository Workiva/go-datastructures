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
package queue

import (
	"errors"
	"log"
	"runtime"
	"sync/atomic"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

// roundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type node struct {
	position uint64
	data     interface{}
}

type nodes []*node

type RingBuffer struct {
	nodes                          nodes
	queue, dequeue, mask, disposed uint64
}

func (rb *RingBuffer) init(size uint64) {
	size = roundUp(size)
	rb.nodes = make(nodes, size)
	for i := uint64(0); i < size; i++ {
		rb.nodes[i] = &node{position: i}
	}
	rb.mask = size - 1 // so we don't have to do this with every put/get operation
}

func (rb *RingBuffer) Put(item interface{}) error {
	var n *node
	pos := atomic.LoadUint64(&rb.queue)
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return disposedError
		}

		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - pos; {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.queue, pos, pos+1) {
				break L
			}
		case dif < 0:
			return errors.New(`Full.`)
		default:
			pos = atomic.LoadUint64(&rb.queue)
		}
		runtime.Gosched() // free up the cpu before the next iteration
	}

	n.data = item
	atomic.StoreUint64(&n.position, pos+1)
	return nil
}

func (rb *RingBuffer) Get() (interface{}, error) {
	var n *node
	pos := atomic.LoadUint64(&rb.dequeue)
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return nil, disposedError
		}

		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - (pos + 1); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.dequeue, pos, pos+1) {
				break L
			}
		case dif < 0:
			return nil, errors.New(`Queue empty.`)
		default:
			pos = atomic.LoadUint64(&rb.dequeue)
		}
		runtime.Gosched() // free up cpu before next iteration
	}
	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.position, pos+rb.mask+1)
	return data, nil
}

func (rb *RingBuffer) Len() uint64 {
	return atomic.LoadUint64(&rb.queue) - atomic.LoadUint64(&rb.dequeue)
}

func (rb *RingBuffer) Cap() uint64 {
	return uint64(len(rb.nodes))
}

func (rb *RingBuffer) Dispose() {
	atomic.CompareAndSwapUint64(&rb.disposed, 0, 1)
}

func (rb *RingBuffer) IsDisposed() bool {
	return atomic.LoadUint64(&rb.disposed) == 1
}

func NewRingBuffer(size uint64) *RingBuffer {
	rb := &RingBuffer{}
	rb.init(size)
	return rb
}
