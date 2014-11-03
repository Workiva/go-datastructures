package augmentedtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkMultiDimensionalRedBlack(tb testing.TB, it *multiDimensionalTree) {
	for _, tree := range it.dimensions {
		checkRedBlack(tb, tree.root, 1)
	}
}

func constructMultiDimensionQueryTestTree() (
	*multiDimensionalTree, Interval, Interval, Interval) {

	it := newMultiDimensionalTree(2)

	iv1 := constructMultiDimensionInterval(
		0, &dimension{low: 5, high: 10}, &dimension{low: 5, high: 10},
	)
	it.Insert(iv1)

	iv2 := constructMultiDimensionInterval(
		1, &dimension{low: 4, high: 5}, &dimension{low: 4, high: 5},
	)
	it.Insert(iv2)

	iv3 := constructMultiDimensionInterval(
		2, &dimension{low: 7, high: 12}, &dimension{low: 7, high: 12},
	)
	it.Insert(iv3)

	return it, iv1, iv2, iv3
}

func TestRootInsertMultipleDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)
	iv := constructMultiDimensionInterval(
		1, &dimension{low: 0, high: 5}, &dimension{low: 1, high: 6},
	)

	it.Insert(iv)

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Equal(t, Intervals{iv}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{100, 200}, &dimension{100, 200},
		),
	)
	assert.Len(t, result, 0)
}

func TestMultipleInsertMultipleDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	checkMultiDimensionalRedBlack(t, it)

	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 100}, &dimension{0, 100},
		),
	)
	assert.Equal(t, Intervals{iv2, iv1, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{3, 5}, &dimension{3, 5},
		),
	)
	assert.Equal(t, Intervals{iv2}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{5, 8}, &dimension{5, 8},
		),
	)
	assert.Equal(t, Intervals{iv1, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{11, 15}, &dimension{11, 15},
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{15, 20}, &dimension{15, 20},
		),
	)
	assert.Len(t, result, 0)
}

func TestInsertRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	for i := int64(0); i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		it.Insert(iv)
	}

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 10)
	assert.Equal(t, 10, it.Len())
}

func TestInsertRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	for i := int64(9); i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		it.Insert(iv)
	}

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 10)
	assert.Equal(t, 10, it.Len())
}

func TestInsertRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	starts := []int64{0, 4, 2, 1, 3}

	for i, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{start, start + 1}, &dimension{start, start + 1},
		)
		it.Insert(iv)
	}

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 5)
	assert.Equal(t, 5, it.Len())
}

func TestInsertLargeNumbersMultiDimensions(t *testing.T) {
	numItems := int64(1000)
	it := newMultiDimensionalTree(2)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		it.Insert(iv)
	}

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, numItems}, &dimension{0, numItems},
		),
	)
	assert.Len(t, result, int(numItems))
	assert.Equal(t, numItems, it.Len())
}

func BenchmarkInsertItemsMultiDimensions(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		intervals = append(intervals, iv)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := newMultiDimensionalTree(2)
		it.Insert(intervals...)
	}
}

func BenchmarkQueryItemsMultiDimensions(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		intervals = append(intervals, iv)
	}

	it := newMultiDimensionalTree(2)
	it.Insert(intervals...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it.Query(
			constructMultiDimensionInterval(
				0, &dimension{0, numItems}, &dimension{0, numItems},
			),
		)
	}
}

func TestRootDeleteMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)
	iv := constructMultiDimensionInterval(
		0, &dimension{low: 5, high: 10}, &dimension{low: 5, high: 10},
	)
	it.Insert(iv)

	it.Delete(iv)

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 100}, &dimension{0, 100},
		),
	)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, it.Len())
}

func TestDeleteMultiDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	checkMultiDimensionalRedBlack(t, it)

	it.Delete(iv1)

	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 100}, &dimension{0, 100},
		),
	)
	assert.Equal(t, Intervals{iv2, iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{3, 5}, &dimension{3, 5},
		),
	)
	assert.Equal(t, Intervals{iv2}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{5, 8}, &dimension{5, 8},
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{11, 15}, &dimension{11, 15},
		),
	)
	assert.Equal(t, Intervals{iv3}, result)

	result = it.Query(
		constructMultiDimensionInterval(
			0, &dimension{15, 20}, &dimension{15, 20},
		),
	)
	assert.Len(t, result, 0)
}

func TestDeleteRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	var toDelete *mockInterval

	for i := int64(0); i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		it.Insert(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 9)
	assert.Equal(t, 9, it.Len())
}

func TestDeleteRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	var toDelete *mockInterval

	for i := int64(9); i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		it.Insert(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 9)
	assert.Equal(t, 9, it.Len())
}

func TestDeleteRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	starts := []int64{0, 4, 2, 1, 3}

	var toDelete *mockInterval

	for i, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{start, start + 1}, &dimension{start, start + 1},
		)
		it.Insert(iv)
		if start == 1 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkMultiDimensionalRedBlack(t, it)
	result := it.Query(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Len(t, result, 4)
	assert.Equal(t, 4, it.Len())
}

func TestDeleteEmptyTreeMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	it.Delete(
		constructMultiDimensionInterval(
			0, &dimension{0, 10}, &dimension{0, 10},
		),
	)
	assert.Equal(t, 0, it.Len())
}

func BenchmarkDeleteItemsMultiDimensions(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{i, i + 1}, &dimension{i, i + 1},
		)
		intervals = append(intervals, iv)
	}

	trees := make([]*multiDimensionalTree, 0, b.N)
	for i := 0; i < b.N; i++ {
		it := newMultiDimensionalTree(2)
		it.Insert(intervals...)
		trees = append(trees, it)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trees[i].Delete(intervals...)
	}
}

func TestMaxMultiDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	assert.Equal(t, 12, it.Max(1))
	assert.Equal(t, 12, it.Max(2))
	assert.Equal(t, 0, it.Max(3))

	it.Delete(iv1, iv2, iv3)

	assert.Equal(t, 0, it.Max(1))
}

func TestMinMultiDimensions(t *testing.T) {
	it, iv1, iv2, iv3 := constructMultiDimensionQueryTestTree()

	assert.Equal(t, 4, it.Min(1))
	assert.Equal(t, 4, it.Min(2))
	assert.Equal(t, 0, it.Min(3))

	it.Delete(iv1, iv2, iv3)

	assert.Equal(t, 0, it.Min(1))
}

func TestInsertDeleteDuplicatesRebalanceInOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	intervals := make(Intervals, 0, 10)

	for i := 0; i < 10; i++ {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{0, 10}, &dimension{0, 10},
		)
		intervals = append(intervals, iv)
	}

	it.Insert(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.dimensions[0].Len())
	assert.Equal(t, 0, it.dimensions[1].Len())
}

func TestInsertDeleteDuplicatesRebalanceReverseOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	intervals := make(Intervals, 0, 10)

	for i := 9; i >= 0; i-- {
		iv := constructMultiDimensionInterval(
			uint64(i), &dimension{0, 10}, &dimension{0, 10},
		)
		intervals = append(intervals, iv)
	}

	it.Insert(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.dimensions[0].Len())
	assert.Equal(t, 0, it.dimensions[1].Len())
}

func TestInsertDeleteDuplicatesRebalanceRandomOrderMultiDimensions(t *testing.T) {
	it := newMultiDimensionalTree(2)

	intervals := make(Intervals, 0, 5)
	starts := []int{0, 4, 2, 1, 3}

	for _, start := range starts {
		iv := constructMultiDimensionInterval(
			uint64(start), &dimension{0, 10}, &dimension{0, 10},
		)
		intervals = append(intervals, iv)
	}

	it.Insert(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.dimensions[0].Len())
	assert.Equal(t, 0, it.dimensions[1].Len())
}
