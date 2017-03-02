package fibheap

// Example usage of the Fibonacci heap

import (
	"fmt"
)

const SomeNumberAround0 float64 = -0.001
const SomeLargerNumberAround15 float64 = 15.77
const SomeNumberAroundMinus1000 float64 = -1002.2001
const SomeNumberAroundMinus1003 float64 = -1003.4

func Example() {
	heap1 := NewFloatFibHeap()
	fmt.Println("Created heap 1.")
	nodeh1_1 := heap1.Enqueue(SomeLargerNumberAround15)
	fmt.Printf("Heap 1 insert: %v\n", nodeh1_1.Priority)

	heap2 := NewFloatFibHeap()
	fmt.Println("Created heap 2.")
	fmt.Printf("Heap 2 is empty? %v\n", heap2.IsEmpty())
	nodeh2_1 := heap2.Enqueue(SomeNumberAroundMinus1000)
	fmt.Printf("Heap 2 insert: %v\n", nodeh2_1.Priority)
	nodeh2_2 := heap2.Enqueue(SomeNumberAround0)
	fmt.Printf("Heap 2 insert: %v\n", nodeh2_2.Priority)
	fmt.Printf("Heap 1 size: %v\n", heap1.Size())
	fmt.Printf("Heap 2 size: %v\n", heap2.Size())
	fmt.Printf("Heap 1 is empty? %v\n", heap1.IsEmpty())
	fmt.Printf("Heap 2 is empty? %v\n", heap2.IsEmpty())

	fmt.Printf("\nMerge Heap 1 and Heap 2.\n")
	mergedHeap, _ := heap1.Merge(&heap2)
	fmt.Printf("Merged heap size: %v\n", mergedHeap.Size())
	fmt.Printf("Set node with priority %v to new priority %v\n", SomeNumberAroundMinus1000, SomeNumberAroundMinus1003)

	mergedHeap.DecreaseKey(nodeh2_1, SomeNumberAroundMinus1003)
	min, _ := mergedHeap.DequeueMin()
	fmt.Printf("Dequeue minimum of merged heap: %v\n", min.Priority)
	fmt.Printf("Merged heap size: %v\n", mergedHeap.Size())

	fmt.Printf("Delete from merged heap: %v\n", SomeNumberAround0)
	mergedHeap.Delete(nodeh2_2)
	fmt.Printf("Merged heap size: %v\n", mergedHeap.Size())

	min, _ = mergedHeap.DequeueMin()
	fmt.Printf("Extracting minimum of merged heap: %v\n", min.Priority)
	fmt.Printf("Merged heap size: %v\n", mergedHeap.Size())
	fmt.Printf("Merged heap is empty? %v\n", mergedHeap.IsEmpty())

	// Output:
	// Created heap 1.
	// Heap 1 insert: 15.77
	// Created heap 2.
	// Heap 2 is empty? true
	// Heap 2 insert: -1002.2001
	// Heap 2 insert: -0.001
	// Heap 1 size: 1
	// Heap 2 size: 2
	// Heap 1 is empty? false
	// Heap 2 is empty? false
	//
	// Merge Heap 1 and Heap 2.
	// Merged heap size: 3
	// Set node with priority -1002.2001 to new priority -1003.4
	// Dequeue minimum of merged heap: -1003.4
	// Merged heap size: 2
	// Delete from merged heap: -0.001
	// Merged heap size: 1
	// Extracting minimum of merged heap: 15.77
	// Merged heap size: 0
	// Merged heap is empty? true
}
