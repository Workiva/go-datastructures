package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiThreadedSortEvenNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(10)
	comparators = reverseComparators(comparators)

	result := MultithreadedSortComparators(comparators)
	comparators = reverseComparators(comparators)

	assert.Equal(t, comparators, result)
}

func TestMultiThreadedSortOddNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(9)
	comparators = reverseComparators(comparators)

	result := MultithreadedSortComparators(comparators)
	comparators = reverseComparators(comparators)

	assert.Equal(t, comparators, result)
}

func BenchmarkMultiThreadedSort(b *testing.B) {
	numCells := 100000

	comparators := constructOrderedMockComparators(numCells)
	comparators = reverseComparators(comparators)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultithreadedSortComparators(comparators)
	}
}
