package batcher

import (
	"errors"
	"sync"
	"time"
)

// Batcher provides an API for accumulating items into a batch for processing.
type Batcher interface {
	// Put adds items to the batcher.
	Put(interface{}) error

	// Get retrieves a batch from the batcher. This call will block until
	// one of the conditions for a "complete" batch is reached.
	Get() ([]interface{}, error)

	// Dispose will dispose of the batcher. Any calls to Put or Get
	// will return errors.
	Dispose()
}

// ErrDisposed is the error returned for a disposed Batcher
var ErrDisposed = errors.New("batcher: disposed")

// CalculateBytes evaluates the number of bytes in an item added to a Batcher.
type CalculateBytes func(interface{}) uint

type basicBatcher struct {
	maxTime        time.Duration
	maxItems       uint
	maxBytes       uint
	queueLen       uint
	calculateBytes CalculateBytes
	disposed       bool
	items          []interface{}
	lock           sync.RWMutex
	batchChan      chan []interface{}
	availableBytes uint
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
		queueLen:       queueLen,
		calculateBytes: calculate,
		items:          make([]interface{}, 0, maxItems),
		batchChan:      make(chan []interface{}, queueLen),
	}, nil
}

// Put adds items to the batcher.
func (b *basicBatcher) Put(item interface{}) error {
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed
	}

	b.items = append(b.items, item)
	b.availableBytes += b.calculateBytes(item)
	if b.ready() {
		b.batchChan <- b.items
		b.items = make([]interface{}, 0, b.maxItems)
		b.availableBytes = 0
	}

	b.lock.Unlock()
	return nil
}

// Get retrieves a batch from the batcher. This call will block until
// one of the conditions for a "complete" batch is reached.
func (b *basicBatcher) Get() ([]interface{}, error) {
	b.lock.RLock()
	if b.disposed {
		b.lock.RUnlock()
		return nil, ErrDisposed
	}
	b.lock.RUnlock()

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

// Dispose will dispose of the batcher. Any calls to Put or Get
// will return errors.
func (b *basicBatcher) Dispose() {
	b.lock.Lock()
	b.disposed = true
	b.items = nil
	close(b.batchChan)
	b.lock.Unlock()
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
