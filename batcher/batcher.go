/*
Copyright 2015 Workiva, LLC

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

package batcher

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Batcher provides an API for accumulating items into a batch for processing.
type Batcher interface {
	// Put adds items to the batcher.
	Put(interface{}) error

	// Get retrieves a batch from the batcher. This call will block until
	// one of the conditions for a "complete" batch is reached.
	Get() ([]interface{}, error)

	// Flush forcibly completes the batch currently being built
	Flush() error

	// Dispose will dispose of the batcher. Any calls to Put or Flush
	// will return ErrDisposed, calls to Get will return an error iff
	// there are no more ready batches.
	Dispose()

	// IsDisposed will determine if the batcher is disposed
	IsDisposed() bool
}

// ErrDisposed is the error returned for a disposed Batcher
var ErrDisposed = errors.New("batcher: disposed")

// CalculateBytes evaluates the number of bytes in an item added to a Batcher.
type CalculateBytes func(interface{}) uint

type basicBatcher struct {
	maxTime        time.Duration
	maxItems       uint
	maxBytes       uint
	calculateBytes CalculateBytes
	disposed       bool
	items          []interface{}
	lock           sync.RWMutex
	batchChan      chan []interface{}
	availableBytes uint
	waiting        int32
}

// New creates a new Batcher using the provided arguments.
// Batch readiness can be determined in three ways:
//   - Maximum number of bytes per batch
//   - Maximum number of items per batch
//   - Maximum amount of time waiting for a batch
// Values of zero for one of these fields indicate they should not be
// taken into account when evaluating the readiness of a batch.
func New(maxTime time.Duration, maxItems, maxBytes, queueLen uint, calculate CalculateBytes) (Batcher, error) {
	if maxBytes > 0 && calculate == nil {
		return nil, errors.New("batcher: must provide CalculateBytes function")
	}

	return &basicBatcher{
		maxTime:        maxTime,
		maxItems:       maxItems,
		maxBytes:       maxBytes,
		calculateBytes: calculate,
		items:          make([]interface{}, 0, maxItems),
		batchChan:      make(chan []interface{}, queueLen),
	}, nil
}

// Put adds items to the batcher. If Put is continually called without calls to
// Get, an unbounded number of go-routines will be generated.
// Note: there is no order guarantee for items entering/leaving the batcher.
func (b *basicBatcher) Put(item interface{}) error {
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed
	}

	b.items = append(b.items, item)
	if b.calculateBytes != nil {
		b.availableBytes += b.calculateBytes(item)
	}
	if b.ready() {
		b.flush()
	}

	b.lock.Unlock()
	return nil
}

// Get retrieves a batch from the batcher. This call will block until
// one of the conditions for a "complete" batch is reached. If Put is
// continually called without calls to Get, an unbounded number of
// go-routines will be generated.
// Note: there is no order guarantee for items entering/leaving the batcher.
func (b *basicBatcher) Get() ([]interface{}, error) {
	// Don't check disposed yet so any items remaining in the queue
	// will be returned properly.

	var timeout <-chan time.Time
	if b.maxTime > 0 {
		timeout = time.After(b.maxTime)
	}

	select {
	case items, ok := <-b.batchChan:
		if !ok {
			return nil, ErrDisposed
		}
		return items, nil
	case <-timeout:
		b.lock.Lock()
		if b.disposed {
			b.lock.Unlock()
			return nil, ErrDisposed
		}
		items := b.items
		b.items = make([]interface{}, 0, b.maxItems)
		b.availableBytes = 0
		b.lock.Unlock()
		return items, nil
	}
}

// Flush forcibly completes the batch currently being built
func (b *basicBatcher) Flush() error {
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed
	}
	b.flush()
	b.lock.Unlock()
	return nil
}

// Dispose will dispose of the batcher. Any calls to Put or Flush
// will return ErrDisposed, calls to Get will return an error iff
// there are no more ready batches.
func (b *basicBatcher) Dispose() {
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return
	}
	b.flush()
	b.disposed = true
	b.items = nil

	// Drain the batch channel and all routines waiting to put on the channel
	for len(b.batchChan) > 0 || atomic.LoadInt32(&b.waiting) > 0 {
		<-b.batchChan
	}
	close(b.batchChan)
	b.lock.Unlock()
}

// IsDisposed will determine if the batcher is disposed
func (b *basicBatcher) IsDisposed() bool {
	b.lock.RLock()
	disposed := b.disposed
	b.lock.RUnlock()
	return disposed
}

// flush adds the batch currently being built to the queue of completed batches.
// flush is not threadsafe, so should be synchronized externally.
func (b *basicBatcher) flush() {
	// Note: This needs to be in a go-routine to avoid locking out gets when
	// the batch channel is full.
	cpItems := make([]interface{}, len(b.items))
	for i, val := range b.items {
		cpItems[i] = val
	}
	// Signal one more waiter for the batch channel
	atomic.AddInt32(&b.waiting, 1)
	// Don't block on the channel put
	go func() {
		b.batchChan <- cpItems
		atomic.AddInt32(&b.waiting, -1)
	}()
	b.items = make([]interface{}, 0, b.maxItems)
	b.availableBytes = 0
}

func (b *basicBatcher) ready() bool {
	if b.maxItems != 0 && uint(len(b.items)) >= b.maxItems {
		return true
	}
	if b.maxBytes != 0 && b.availableBytes >= b.maxBytes {
		return true
	}
	return false
}
