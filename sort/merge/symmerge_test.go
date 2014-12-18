package merge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockComparator int

func (mc mockComparator) Compare(other Comparator) int {
	if mc == other.(mockComparator) {
		return 0
	}

	if mc > other.(mockComparator) {
		return 1
	}

	return -1
}

func constructMockComparators(values ...int) Comparators {
	comparators := make(Comparators, 0, len(values))
	for _, v := range values {
		comparators = append(comparators, mockComparator(v))
	}

	return comparators
}

func constructOrderedMockComparators(upTo int) Comparators {
	comparators := make(Comparators, 0, upTo)
	for i := 0; i < upTo; i++ {
		comparators = append(comparators, mockComparator(i))
	}

	return comparators
}

func reverseComparators(comparators Comparators) Comparators {
	for i := 0; i < len(comparators); i++ {
		li := len(comparators) - i - 1
		comparators[i], comparators[li] = comparators[li], comparators[i]
	}
	return comparators
}

func TestMultiThreadedSortEvenNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(10)
	comparators = reverseComparators(comparators)

	result := multithreadedSortComparators(comparators)
	comparators = reverseComparators(comparators)

	assert.Equal(t, comparators, result)
}

func TestMultiThreadedSortOddNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(9)
	comparators = reverseComparators(comparators)

	result := multithreadedSortComparators(comparators)
	comparators = reverseComparators(comparators)

	assert.Equal(t, comparators, result)
}

func BenchmarkMultiThreadedSort(b *testing.B) {
	numCells := 100000

	comparators := constructOrderedMockComparators(numCells)
	comparators = reverseComparators(comparators)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		multithreadedSortComparators(comparators)
	}
}
func TestDecomposeForSymMergeOddNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(7)

	v1, w, v2 := decomposeForSymMerge(3, comparators)
	assert.Equal(t, comparators[:2], v1)
	assert.Equal(t, comparators[2:5], w)
	assert.Equal(t, comparators[5:], v2)
}

func TestDecomposeForSymMergeEvenNumber(t *testing.T) {
	comparators := constructOrderedMockComparators(8)
	v1, w, v2 := decomposeForSymMerge(5, comparators)

	assert.Equal(t, comparators[:1], v1)
	assert.Equal(t, comparators[1:6], w)
	assert.Equal(t, comparators[6:], v2)
}

func TestNearCompleteDecomposeForSymMerge(t *testing.T) {
	comparators := constructOrderedMockComparators(8)
	v1, w, v2 := decomposeForSymMerge(7, comparators)

	assert.Len(t, v1, 0)
	assert.Equal(t, comparators[:7], w)
	assert.Equal(t, comparators[7:], v2)
}

func TestDecomposePanicsWithWrongLength(t *testing.T) {
	comparators := constructOrderedMockComparators(8)
	assert.Panics(t, func() {
		decomposeForSymMerge(8, comparators)
	})
}

func TestSymSearch(t *testing.T) {
	u := constructMockComparators(1, 5, 7, 9)
	w := constructMockComparators(1, 3, 9, 10)

	result := symSearch(u, w)
	assert.Equal(t, 2, result)

	u = constructMockComparators(1, 5, 7)
	w = constructMockComparators(1, 3, 9)

	result = symSearch(u, w)
	assert.Equal(t, 1, result)
}

func TestSwap(t *testing.T) {
	u := constructMockComparators(1, 5, 7, 9)
	w := constructMockComparators(2, 8, 11, 13)
	u1 := constructMockComparators(1, 5, 2, 8)
	w1 := constructMockComparators(7, 9, 11, 13)

	swap(u, w, 2)

	assert.Equal(t, u1, u)
	assert.Equal(t, w1, w)
}

func TestSymMergeSmallLists(t *testing.T) {
	u := constructMockComparators(1, 5)
	w := constructMockComparators(2, 8)
	expected := constructMockComparators(1, 2, 5, 8)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestSymMergeAlreadySorted(t *testing.T) {
	u := constructMockComparators(1, 5)
	w := constructMockComparators(6, 7)
	expected := constructMockComparators(1, 5, 6, 7)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestSymMergeAlreadySortedReverseOrder(t *testing.T) {
	u := constructMockComparators(6, 7)
	w := constructMockComparators(1, 5)
	expected := constructMockComparators(1, 5, 6, 7)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestSymMergeUnevenLists(t *testing.T) {
	u := constructMockComparators(1, 3, 7)
	w := constructMockComparators(2, 4)
	expected := constructMockComparators(1, 2, 3, 4, 7)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestSymMergeUnevenListsWrongOrder(t *testing.T) {
	u := constructMockComparators(2, 4)
	w := constructMockComparators(1, 3, 7)
	expected := constructMockComparators(1, 2, 3, 4, 7)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenLists(t *testing.T) {
	u := constructMockComparators(1, 3, 7, 12, 15)
	w := constructMockComparators(2, 4)
	expected := constructMockComparators(1, 2, 3, 4, 7, 12, 15)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenListsWrongOrder(t *testing.T) {
	u := constructMockComparators(2, 4)
	w := constructMockComparators(1, 3, 7, 12, 15)
	expected := constructMockComparators(1, 2, 3, 4, 7, 12, 15)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenListsAlreadySorted(t *testing.T) {
	u := constructMockComparators(2, 4)
	w := constructMockComparators(5, 7, 9, 10, 11, 12)
	expected := constructMockComparators(2, 4, 5, 7, 9, 10, 11, 12)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenListsAlreadySortedWrongOrder(t *testing.T) {
	w := constructMockComparators(2, 4)
	u := constructMockComparators(5, 7, 9, 10, 11, 12)
	expected := constructMockComparators(2, 4, 5, 7, 9, 10, 11, 12)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenListIsSubset(t *testing.T) {
	u := constructMockComparators(2, 4)
	w := constructMockComparators(1, 3, 5, 7, 9)
	expected := constructMockComparators(1, 2, 3, 4, 5, 7, 9)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeVeryUnevenListIsSubsetReverseOrder(t *testing.T) {
	w := constructMockComparators(2, 4)
	u := constructMockComparators(1, 3, 5, 7, 9)
	expected := constructMockComparators(1, 2, 3, 4, 5, 7, 9)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeUnevenOneListIsOne(t *testing.T) {
	u := constructMockComparators(1)
	w := constructMockComparators(0, 3, 5, 7, 9)
	expected := constructMockComparators(0, 1, 3, 5, 7, 9)

	u = SymMerge(u, w)
	assert.Equal(t, expected, u)
}

func TestMergeEmptyList(t *testing.T) {
	u := constructMockComparators(1, 3, 5)
	expected := constructMockComparators(1, 3, 5)

	u = SymMerge(u, nil)
	assert.Equal(t, expected, u)
}
