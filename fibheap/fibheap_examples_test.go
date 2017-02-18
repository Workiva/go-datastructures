package fibheap

// Tests for the Fibonacci heap with floating point number priorities

import (
	"fmt"
)

const SomeNumber float64 = 15.5
const SomeSmallerNumber float64 = -10.1
const SomeLargerNumber float64 = 112.211

func ExampleFloatingFibonacciHeap_Enqueue() {
	heap := NewFloatFibHeap()
	// The function returns a pointer
	// to the node that contains the new value
	node := heap.Enqueue(SomeNumber)
	fmt.Println(node.Priority)
	// Output: 15.5
}

func ExampleFloatingFibonacciHeap_Min() {
	heap := NewFloatFibHeap()
	heap.Enqueue(SomeNumber)
	heap.Enqueue(SomeLargerNumber)
	min, _ := heap.Min()
	fmt.Println(min.Priority)
	// Output: 15.5
}

func ExampleFloatingFibonacciHeap_IsEmpty() {
	heap := NewFloatFibHeap()
	fmt.Printf("Empty before insert? %v\n", heap.IsEmpty())
	heap.Enqueue(SomeNumber)
	fmt.Printf("Empty after insert? %v\n", heap.IsEmpty())
	// Output:
	// Empty before insert? true
	// Empty after insert? false
}

func ExampleFloatingFibonacciHeap_Size() {
	heap := NewFloatFibHeap()
	fmt.Printf("Size before insert: %v\n", heap.Size())
	heap.Enqueue(SomeNumber)
	fmt.Printf("Size after insert: %v\n", heap.Size())
	// Output:
	// Size before insert: 0
	// Size after insert: 1
}

func ExampleFloatingFibonacciHeap_DequeueMin() {
	heap := NewFloatFibHeap()
	heap.Enqueue(SomeNumber)
	node, _ := heap.DequeueMin()
	fmt.Printf("Dequeueing minimal element: %v\n", node.Priority)
	// Output:
	// Dequeueing minimal element: 15.5
}

func ExampleFloatingFibonacciHeap_DecreaseKey() {
	heap := NewFloatFibHeap()
	node := heap.Enqueue(SomeNumber)
	min, _ := heap.Min()
	fmt.Printf("Minimal element before decreasing key: %v\n", min.Priority)
	heap.DecreaseKey(node, SomeSmallerNumber)
	min, _ = heap.Min()
	fmt.Printf("Minimal element after decreasing key: %v\n", min.Priority)
	// Output:
	// Minimal element before decreasing key: 15.5
	// Minimal element after decreasing key: -10.1
}

func ExampleFloatingFibonacciHeap_Delete() {
	heap := NewFloatFibHeap()
	node := heap.Enqueue(SomeNumber)
	heap.Enqueue(SomeLargerNumber)
	min, _ := heap.Min()
	fmt.Printf("Minimal element before deletion: %v\n", min.Priority)
	heap.Delete(node)
	min, _ = heap.Min()
	fmt.Printf("Minimal element after deletion: %v\n", min.Priority)
	// Output:
	// Minimal element before deletion: 15.5
	// Minimal element after deletion: 112.211
}

func ExampleFloatingFibonacciHeap_Merge() {
	heap1 := NewFloatFibHeap()
	heap2 := NewFloatFibHeap()
	heap1.Enqueue(SomeNumber)
	heap1.Enqueue(SomeLargerNumber)
	heap2.Enqueue(SomeSmallerNumber)
	min, _ := heap1.Min()
	fmt.Printf("Minimal element of heap 1: %v\n", min.Priority)
	min, _ = heap2.Min()
	fmt.Printf("Minimal element of heap 2: %v\n", min.Priority)
	heap, _ := heap1.Merge(heap2)
	min, _ = heap.Min()
	fmt.Printf("Minimal element of merged heap: %v\n", min.Priority)
	// Output:
	// Minimal element of heap 1: 15.5
	// Minimal element of heap 2: -10.1
	// Minimal element of merged heap: -10.1
}
