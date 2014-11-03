/*
The priority queue is almost a spitting image of the logic
used for a regular queue.  In order to keep the logic fast,
this code is repeated instead of using casts to cast to interface{}
back and forth.  If Go had inheritance and generics, this problem
would be easier to solve.
*/
package queue

import (
	"sort"
	"sync"
)

// Item is an item that can be added to the priority queue.
type Item interface {
	// Compare returns a bool that can be used to determine
	// ordering in the priority queue.  Assuming the queue
	// is in ascending order, this should return >= logic.
	Compare(other Item) bool
}

type priorityItems []Item

func (items *priorityItems) get(number int) []Item {
	returnItems := make([]Item, 0, number)
	index := 0
	for i := 0; i < number; i++ {
		if i >= len(*items) {
			break
		}

		returnItems = append(returnItems, (*items)[i])
		index++
	}

	*items = (*items)[index:]
	return returnItems
}

func (items *priorityItems) insert(item Item) {
	if len(*items) == 0 {
		*items = append(*items, item)
		return
	}

	i := sort.Search(len(*items), func(i int) bool {
		return (*items)[i].Compare(item)
	})

	if i == len(*items) {
		*items = append(*items, item)
		return
	}

	*items = append(*items, nil)
	copy((*items)[i+1:], (*items)[i:])
	(*items)[i] = item
}

// PriorityQueue is similar to queue except that it takes
// items that implement the Item interface and adds them
// to the queue in priority order.
type PriorityQueue struct {
	waiters     waiters
	items       priorityItems
	lock        sync.Mutex
	disposeLock sync.Mutex
	disposed    bool
}

// Put adds items to the queue.
func (pq *PriorityQueue) Put(items ...Item) error {
	if len(items) == 0 {
		return nil
	}

	pq.lock.Lock()
	if pq.disposed {
		pq.lock.Unlock()
		return DisposedError{}
	}

	for _, item := range items {
		pq.items.insert(item)
	}

	for {
		sema := pq.waiters.get()
		if sema == nil {
			break
		}

		sema.response.Add(1)
		sema.wg.Done()
		sema.response.Wait()
		if len(pq.items) == 0 {
			break
		}
	}

	pq.lock.Unlock()
	return nil
}

// Get retrieves items from the queue.  If the queue is empty,
// this call blocks until the next item is added to the queue.  This
// will attempt to retrieve number of items.
func (pq *PriorityQueue) Get(number int) ([]Item, error) {
	if number < 1 {
		return nil, nil
	}

	pq.lock.Lock()

	if pq.disposed {
		pq.lock.Unlock()
		return nil, DisposedError{}
	}

	var items []Item

	if len(pq.items) == 0 {
		sema := newSema()
		pq.waiters.put(sema)
		sema.wg.Add(1)
		pq.lock.Unlock()

		sema.wg.Wait()
		pq.disposeLock.Lock()
		if pq.disposed {
			pq.disposeLock.Unlock()
			return nil, DisposedError{}
		}
		pq.disposeLock.Unlock()

		items = pq.items.get(number)
		sema.response.Done()
		return items, nil
	}

	items = pq.items.get(number)
	pq.lock.Unlock()
	return items, nil
}

// Peek will look at the next item without removing it from the queue.
func (pq *PriorityQueue) Peek() Item {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	if len(pq.items) > 0 {
		return pq.items[0]
	}
	return nil
}

// Empty returns a bool indicating if there are any items left
// in the queue.
func (pq *PriorityQueue) Empty() bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return len(pq.items) == 0
}

// Len returns a number indicating how many items are in the queue.
func (pq *PriorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return len(pq.items)
}

// Disposed returns a bool indicating if this queue has been disposed.
func (pq *PriorityQueue) Disposed() bool {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	return pq.disposed
}

// Dispose will prevent any further reads/writes to this queue
// and frees available resources.
func (pq *PriorityQueue) Dispose() {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	pq.disposeLock.Lock()
	defer pq.disposeLock.Unlock()

	pq.disposed = true
	for _, waiter := range pq.waiters {
		waiter.response.Add(1)
		waiter.wg.Done()
	}

	pq.items = nil
	pq.waiters = nil
}

// NewPriorityQueue is the constructor for a priority queue.
func NewPriorityQueue(hint int) *PriorityQueue {
	return &PriorityQueue{
		items: make(priorityItems, 0, hint),
	}
}
