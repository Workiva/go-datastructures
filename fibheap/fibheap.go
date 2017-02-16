/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Special thanks to  Keith Schwarz (htiek@cs.stanford.edu),
whose code and documentation have been used as a reference
for the algorithm implementation.
http://www.keithschwarz.com/interesting/code/?dir=fibonacci-heap
*/

/*
* An implementation of a priority queue backed by a Fibonacci heap,
* as described by Fredman and Tarjan.  Fibonacci heaps are interesting
* theoretically because they have asymptotically good runtime guarantees
* for many operations.  In particular, insert, peek, and decrease-key all
* run in amortized O(1) time.  dequeueMin and delete each run in amortized
* O(lg n) time.  This allows algorithms that rely heavily on decrease-key
* to gain significant performance boosts.  For example, Dijkstra's algorithm
* for single-source shortest paths can be shown to run in O(m + n lg n) using
* a Fibonacci heap, compared to O(m lg n) using a standard binary or binomial
* heap.
*
* Internally, a Fibonacci heap is represented as a circular, doubly-linked
* list of trees obeying the min-heap property.  Each node stores pointers
* to its parent (if any) and some arbitrary child.  Additionally, every
* node stores its degree (the number of children it has) and whether it
* is a "marked" node.  Finally, each Fibonacci heap stores a pointer to
* the tree with the minimum value.
*
* To insert a node into a Fibonacci heap, a singleton tree is created and
* merged into the rest of the trees.  The merge operation works by simply
* splicing together the doubly-linked lists of the two trees, then updating
* the min pointer to be the smaller of the minima of the two heaps.  Peeking
* at the smallest element can therefore be accomplished by just looking at
* the min element.  All of these operations complete in O(1) time.
*
* The tricky operations are dequeueMin and decreaseKey.  dequeueMin works
* by removing the root of the tree containing the smallest element, then
* merging its children with the topmost roots.  Then, the roots are scanned
* and merged so that there is only one tree of each degree in the root list.
* This works by maintaining a dynamic array of trees, each initially null,
* pointing to the roots of trees of each dimension.  The list is then scanned
* and this array is populated.  Whenever a conflict is discovered, the
* appropriate trees are merged together until no more conflicts exist.  The
* resulting trees are then put into the root list.  A clever analysis using
* the potential method can be used to show that the amortized cost of this
* operation is O(lg n), see "Introduction to Algorithms, Second Edition" by
* Cormen, Rivest, Leiserson, and Stein for more details.
*
* The other hard operation is decreaseKey, which works as follows.  First, we
* update the key of the node to be the new value.  If this leaves the node
* smaller than its parent, we're done.  Otherwise, we cut the node from its
* parent, add it as a root, and then mark its parent.  If the parent was
* already marked, we cut that node as well, recursively mark its parent,
* and continue this process.  This can be shown to run in O(1) amortized time
* using yet another clever potential function.  Finally, given this function,
* we can implement delete by decreasing a key to -\infty, then calling
* dequeueMin to extract it.
 */

package fibheap

import (
	"fmt"
	"math"
)

/******************************************
 ************** INTERFACE *****************
 ******************************************/

// The FloatingFibonacciHeap is an implementation of a fibonacci heap
// with only floating-point priorities and no user data attached.
type FloatingFibonacciHeap interface {
	// Adds and element to the heap
	Enqueue(priority float64) *Entry
	// Returns the minimum element in the heap
	Min() (*Entry, error)
	// Is the heap empty?
	IsEmpty() bool
	// The number of elements in the heap
	Size() uint
	// Removes and returns the minimal element
	// in the heap
	DequeueMin() (*Entry, error)
	// Decreases the key of the given element
	// and sets it to the new given priority
	// returns the node if succesfully set
	DecreaseKey(node *Entry, newPriority float64) (*Entry, error)
	// Deletes the given element in the heap
	Delete(node *Entry) error
	// Will merge two heaps
	Merge(otherHeap FloatingFibonacciHeap) (FloatingFibonacciHeap, error)
}

// Entry is the entry type that will be used
// for each node of the Fibonacci heap
type Entry struct {
	degree                    int
	marked                    bool
	next, prev, child, parent *Entry
	priority                  float64
}

/******************************************
 ************** END INTERFACE *************
 ******************************************/

type fibHeap struct {
	min  *Entry // The minimal element
	size uint   // Size of the heap
}

// ****************
// HELPER FUNCTIONS
// ****************

func newEntry(priority float64) *Entry {
	result := new(Entry)
	result.degree = 0
	result.marked = false
	result.child = nil
	result.parent = nil
	result.next = result
	result.prev = result
	result.priority = priority
	return result
}

// ***********
// ACTUAL CODE
// ***********

/*
NewFloatFibHeap creates a new, empty, Fibonacci heap object.

Remember that it's actually *fibHeap and not fibHeap
that fulfills the contract of the interface.
*/
func NewFloatFibHeap() FloatingFibonacciHeap { return &fibHeap{nil, 0} }

func (heap *fibHeap) Enqueue(priority float64) *Entry {
	singleton := newEntry(priority)

	// Merge singleton list with heap
	heap.min = mergeLists(heap.min, singleton)
	heap.size++
	return singleton
}

func (heap *fibHeap) Min() (*Entry, error) {
	if heap.IsEmpty() {
		return nil, fmt.Errorf("Trying to get minimum element of empty heap")
	}
	return heap.min, nil
}

func (heap *fibHeap) IsEmpty() bool {
	return heap.size == 0
}

func (heap *fibHeap) Size() uint {
	return heap.size
}

func (heap *fibHeap) DequeueMin() (*Entry, error) {
	if heap.IsEmpty() {
		return nil, fmt.Errorf("Heap is empty")
	}

	heap.size--

	// Copy pointer. Will need it later.
	min := heap.min

	if min.next == min { // This is the only root node
		heap.min = nil
	} else { // There are more root nodes
		heap.min.prev.next = heap.min.next
		heap.min.next.prev = heap.min.prev
		heap.min = heap.min.next // Arbitrary element of the root list
	}

	if min.child != nil {
		// Keep track of the first visited node
		curr := min.child
		for ok := true; ok; ok = (curr != min.child) {
			curr.parent = nil
			curr = curr.next
		}
	}

	heap.min = mergeLists(heap.min, min.child)

	if heap.min == nil {
		// If there are no entries left, we're done.
		return min, nil
	}

	treeSlice := make([]*Entry, 0, heap.size)
	toVisit := make([]*Entry, 0, heap.size)

	for curr := heap.min; len(toVisit) == 0 || toVisit[0] != curr; curr = curr.next {
		toVisit = append(toVisit, curr)
	}

	for _, curr := range toVisit {
		for {
			for curr.degree >= len(treeSlice) {
				treeSlice = append(treeSlice, nil)
			}

			if treeSlice[curr.degree] == nil {
				treeSlice[curr.degree] = curr
				break
			}

			other := treeSlice[curr.degree]
			treeSlice[curr.degree] = nil

			// Determine which of two trees has the smaller root
			var minT, maxT *Entry
			if other.priority < curr.priority {
				minT = other
				maxT = curr
			} else {
				minT = curr
				maxT = other
			}

			// Break max out of the root list,
			// then merge it into min's child list
			maxT.next.prev = maxT.prev
			maxT.prev.next = maxT.next

			// Make it a singleton so that we can merge it
			maxT.prev = maxT
			maxT.next = maxT
			minT.child = mergeLists(minT.child, maxT)

			// Reparent max appropriately
			maxT.parent = minT

			// Clear max's mark, since it can now lose another child
			maxT.marked = false

			// Increase min's degree. It has another child.
			minT.degree++

			// Continue merging this tree
			curr = minT
		}

		/* Update the global min based on this node.  Note that we compare
		 * for <= instead of < here.  That's because if we just did a
		 * reparent operation that merged two different trees of equal
		 * priority, we need to make sure that the min pointer points to
		 * the root-level one.
		 */
		if curr.priority <= heap.min.priority {
			heap.min = curr
		}
	}

	// All done. Return minimum element and no error
	return min, nil
}

func (heap *fibHeap) DecreaseKey(node *Entry, newPriority float64) (*Entry, error) {

	if newPriority > node.priority {
		return nil, fmt.Errorf("The given new priority is larger than the old")
	}

	decreaseKeyUnchecked(heap, node, newPriority)
	return node, nil
}

func (heap *fibHeap) Delete(node *Entry) error {

	decreaseKeyUnchecked(heap, node, -math.MaxFloat64)
	heap.DequeueMin()
	return nil
}

/*
 * Given two Fibonacci heaps, returns a new Fibonacci heap that contains
 * all of the elements of the two heaps.  Each of the input heaps is
 * destructively modified by having all its elements removed.  You can
 * continue to use those heaps, but be aware that they will be empty
 * after this call completes.
 */
func (heap *fibHeap) Merge(other FloatingFibonacciHeap) (FloatingFibonacciHeap, error) {

	otherHeap, ok := other.(*fibHeap)
	if !ok {
		// throw an error
		return nil, fmt.Errorf("The passed object is of type %T, not of internal type *fibHeap. Please provide your own implementation of merge", other)
	}

	resultSize := heap.size + otherHeap.size

	resultMin := mergeLists(heap.min, otherHeap.min)

	heap.min = nil
	otherHeap.min = nil
	heap.size = 0
	otherHeap.size = 0

	return &fibHeap{resultMin, resultSize}, nil
}

func mergeLists(one, two *Entry) *Entry {
	if one == nil && two == nil {
		return nil
	} else if one != nil && two == nil {
		return one
	} else if one == nil && two != nil {
		return two
	}
	// Both trees non-null; actually do the merge.
	oneNext := one.next
	one.next = two.next
	one.next.prev = one
	two.next = oneNext
	two.next.prev = two

	if one.priority < two.priority {
		return one
	}
	return two

}

func decreaseKeyUnchecked(heap *fibHeap, node *Entry, priority float64) {
	node.priority = priority

	if node.parent != nil && node.priority <= node.parent.priority {
		cutNode(heap, node)
	}

	if node.priority <= heap.min.priority {
		heap.min = node
	}
}

func cutNode(heap *fibHeap, node *Entry) {
	node.marked = false

	if node.parent == nil {
		return
	}

	// Rewire siblings if it has any
	if node.next != node {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	// Rewrite pointer if this is the representative child node
	if node.parent.child == node {
		if node.next != node {
			node.parent.child = node.next
		} else {
			node.parent.child = nil
		}
	}

	node.parent.degree--

	node.prev = node
	node.next = node
	heap.min = mergeLists(heap.min, node)

	// cut parent recursively if marked
	if node.parent.marked {
		cutNode(heap, node.parent)
	} else {
		node.parent.marked = true
	}

	node.parent = nil
}
