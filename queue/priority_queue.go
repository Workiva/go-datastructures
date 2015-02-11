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
	// is in ascending order, this should return > logic.
	// Return 1 to indicate this object is greater than the
	// the other logic, 0 to indicate equality, and -1 to indicate
	// less than other.
	Compare(other Item) int
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
		(*items)[i] = nil
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

	equalFound := false
	i := sort.Search(len(*items), func(i int) bool {
		result := (*items)[i].Compare(item)
		if result == 0 {
			equalFound = true
		}
		return result >= 0
	})

	if equalFound {
		return
	}

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
		return disposedError
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
		return nil, disposedError
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
			return nil, disposedError
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
