package augmentedtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func min(one, two int64) int64 {
	if one == -1 {
		return two
	}

	if two == -1 {
		return one
	}

	if one > two {
		return two
	}

	return one
}

func max(one, two int64) int64 {
	if one == -1 {
		return two
	}

	if two == -1 {
		return one
	}

	if one > two {
		return one
	}

	return two
}

func checkRedBlack(tb testing.TB, node *node, dimension int) (int64, int64, int64) {
	lh, rh := 0, 0
	if node == nil {
		return 1, -1, -1
	}

	if isRed(node) {
		if isRed(node.children[0]) || isRed(node.children[1]) {
			tb.Errorf(`Node is red and has red children: %+v`, node)
		}
	}

	fn := func(min, max int64) {
		if min != -1 && min < node.min {
			tb.Errorf(`Min not set correctly: %+v, node: %+v`, min, node)
		}

		if max != -1 && max > node.max {
			tb.Errorf(`Max not set correctly: %+v, node: %+v`, max, node)
		}
	}

	left, minL, maxL := checkRedBlack(tb, node.children[0], dimension)
	fn(minL, maxL)
	right, minR, maxR := checkRedBlack(tb, node.children[1], dimension)
	fn(minR, maxR)

	min := min(minL, minR)
	if min == -1 && node.min != node.low {
		tb.Errorf(`Min not set correctly, node: %+v`, node)
	} else if min != -1 && node.children[0] != nil && node.children[0].min != node.min {
		tb.Errorf(`Min not set correctly: node: %+v, child: %+v`, node, node.children[0])
	} else if min != -1 && node.children[0] == nil && node.min != node.low {
		tb.Errorf(`Min not set correctly: %+v`, node)
	}

	max := max(maxL, maxR)
	if max == -1 && node.max != node.high {
		tb.Errorf(`Max not set correctly, node: %+v`, node)
	} else if max > node.high && max != node.max {
		tb.Errorf(`Max not set correctly, max: %+v, node: %+v`, max, node)
	}

	if left != 0 && right != 0 && lh != rh {
		tb.Errorf(`Black violation: left: %d, right: %d`, left, right)
	}

	if left != 0 && right != 0 {
		if isRed(node) {
			return left, node.min, node.max
		}

		return left + 1, node.min, node.max
	}

	return 0, node.min, node.max
}

func constructSingleDimensionTestTree(number int) (*tree, Intervals) {
	tree := newTree(1)

	ivs := make(Intervals, 0, number)
	for i := 0; i < number; i++ {
		iv := constructSingleDimensionInterval(int64(i), int64(i)+10, uint64(i))
		ivs = append(ivs, iv)
	}

	tree.Add(ivs...)
	return tree, ivs
}

func TestSimpleAddNilRoot(t *testing.T) {
	it := newTree(1)

	iv := constructSingleDimensionInterval(5, 10, 0)

	it.Add(iv)

	expected := newNode(iv, 5, 10, 1)
	expected.red = false

	assert.Equal(t, expected, it.root)
	assert.Equal(t, 1, it.Len())
	checkRedBlack(t, it.root, 1)
}

func TestSimpleAddRootLeft(t *testing.T) {
	it := newTree(1)

	iv := constructSingleDimensionInterval(5, 10, 0)
	it.Add(iv)

	expectedRoot := newNode(iv, 4, 11, 1)
	expectedRoot.red = false

	iv = constructSingleDimensionInterval(4, 11, 1)
	it.Add(iv)

	expectedChild := newNode(iv, 4, 11, 1)
	expectedRoot.children[0] = expectedChild

	assert.Equal(t, expectedRoot, it.root)
	assert.Equal(t, 2, it.Len())
	checkRedBlack(t, it.root, 1)
}

func TestSimpleAddRootRight(t *testing.T) {
	it := newTree(1)

	iv := constructSingleDimensionInterval(5, 10, 0)
	it.Add(iv)

	expectedRoot := newNode(iv, 5, 11, 1)
	expectedRoot.red = false

	iv = constructSingleDimensionInterval(7, 11, 1)
	it.Add(iv)

	expectedChild := newNode(iv, 7, 11, 1)
	expectedRoot.children[1] = expectedChild

	assert.Equal(t, expectedRoot, it.root)
	assert.Equal(t, 2, it.Len())
	checkRedBlack(t, it.root, 1)
}

func TestAddRootLeftAndRight(t *testing.T) {
	it := newTree(1)

	iv := constructSingleDimensionInterval(5, 10, 0)
	it.Add(iv)

	expectedRoot := newNode(iv, 4, 12, 1)
	expectedRoot.red = false

	iv = constructSingleDimensionInterval(4, 11, 1)
	it.Add(iv)

	expectedLeft := newNode(iv, 4, 11, 1)
	expectedRoot.children[0] = expectedLeft

	iv = constructSingleDimensionInterval(7, 12, 1)
	it.Add(iv)

	expectedRight := newNode(iv, 7, 12, 1)
	expectedRoot.children[1] = expectedRight

	assert.Equal(t, expectedRoot, it.root)
	assert.Equal(t, 3, it.Len())
	checkRedBlack(t, it.root, 1)
}

func TestAddRebalanceInOrder(t *testing.T) {
	it := newTree(1)

	for i := int64(0); i < 10; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		it.add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 10)
	assert.Equal(t, 10, it.Len())
}

func TestAddRebalanceOutOfOrder(t *testing.T) {
	it := newTree(1)

	for i := int64(9); i >= 0; i-- {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		it.add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 10)
	assert.Equal(t, 10, it.Len())
}

func TestAddRebalanceRandomOrder(t *testing.T) {
	it := newTree(1)

	starts := []int64{0, 4, 2, 1, 3}

	for _, start := range starts {
		iv := constructSingleDimensionInterval(start, start+1, uint64(start))
		it.add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 5)
	assert.Equal(t, 5, it.Len())
}

func TestAddLargeNumberOfItems(t *testing.T) {
	numItems := int64(1000)
	it := newTree(1)

	for i := int64(0); i < numItems; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		it.add(iv)
	}

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, numItems, 0))
	assert.Len(t, result, int(numItems))
	assert.Equal(t, numItems, it.Len())
}

func BenchmarkAddItems(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		intervals = append(intervals, iv)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		it := newTree(1)
		it.Add(intervals...)
	}
}

func BenchmarkQueryItems(b *testing.B) {
	numItems := int64(1000)
	intervals := make(Intervals, 0, numItems)

	for i := int64(0); i < numItems; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		intervals = append(intervals, iv)
	}

	it := newTree(1)
	it.Add(intervals...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it.Query(constructSingleDimensionInterval(0, numItems, 0))
	}
}

func constructSingleDimensionQueryTestTree() (
	*tree, Interval, Interval, Interval) {

	it := newTree(1)

	iv1 := constructSingleDimensionInterval(5, 10, 0)
	it.Add(iv1)

	iv2 := constructSingleDimensionInterval(4, 5, 1)
	it.Add(iv2)

	iv3 := constructSingleDimensionInterval(7, 12, 2)
	it.Add(iv3)

	return it, iv1, iv2, iv3
}

func TestSimpleQuery(t *testing.T) {
	it, iv1, iv2, _ := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(3, 6, 0))

	expected := Intervals{iv2, iv1}
	assert.Equal(t, expected, result)
}

func TestRightQuery(t *testing.T) {
	it, iv1, _, iv3 := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(5, 8, 0))

	expected := Intervals{iv1, iv3}
	assert.Equal(t, expected, result)
}

func TestLeftQuery(t *testing.T) {
	it, _, iv2, _ := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(3, 5, 0))

	expected := Intervals{iv2}
	assert.Equal(t, expected, result)
}

func TestMatchingQuery(t *testing.T) {
	it, _, iv2, _ := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(4, 5, 0))

	expected := Intervals{iv2}
	assert.Equal(t, expected, result)
}

func TestNoMatchLeft(t *testing.T) {
	it, _, _, _ := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(1, 4, 0))

	expected := Intervals{}
	assert.Equal(t, expected, result)
}

func TestNoMatchRight(t *testing.T) {
	it, _, _, _ := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(12, 13, 0))

	expected := Intervals{}
	assert.Equal(t, expected, result)
}

func TestAllQuery(t *testing.T) {
	it, iv1, iv2, iv3 := constructSingleDimensionQueryTestTree()

	result := it.Query(constructSingleDimensionInterval(1, 14, 0))

	expected := Intervals{iv2, iv1, iv3}
	assert.Equal(t, expected, result)
}

func TestQueryDuplicate(t *testing.T) {
	it, _, iv2, _ := constructSingleDimensionQueryTestTree()
	iv4 := constructSingleDimensionInterval(4, 5, 3)
	it.Add(iv4)

	result := it.Query(constructSingleDimensionInterval(4, 5, 0))

	expected := Intervals{iv2, iv4}
	assert.Equal(t, expected, result)
}

func TestRootDelete(t *testing.T) {
	it := newTree(1)
	iv := constructSingleDimensionInterval(1, 5, 1)
	it.add(iv)

	it.Delete(iv)

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(1, 10, 0))
	assert.Len(t, result, 0)
	assert.Equal(t, 0, it.Len())
}

func TestDeleteLeft(t *testing.T) {
	it, iv1, iv2, iv3 := constructSingleDimensionQueryTestTree()

	it.Delete(iv2)

	expected := Intervals{iv1, iv3}

	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	checkRedBlack(t, it.root, 1)
	assert.Equal(t, expected, result)
	assert.Equal(t, 2, it.Len())
}

func TestDeleteRight(t *testing.T) {
	it, iv1, iv2, iv3 := constructSingleDimensionQueryTestTree()

	it.Delete(iv3)

	expected := Intervals{iv2, iv1}

	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	checkRedBlack(t, it.root, 1)
	assert.Equal(t, expected, result)
	assert.Equal(t, 2, it.Len())
}

func TestDeleteCenter(t *testing.T) {
	it, iv1, iv2, iv3 := constructSingleDimensionQueryTestTree()

	it.Delete(iv1)

	expected := Intervals{iv2, iv3}

	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	checkRedBlack(t, it.root, 1)
	assert.Equal(t, expected, result)
	assert.Equal(t, 2, it.Len())
}

func TestDeleteRebalanceInOrder(t *testing.T) {
	it := newTree(1)

	var toDelete *mockInterval

	for i := int64(0); i < 10; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		it.add(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 9)
	assert.Equal(t, 9, it.Len())
}

func TestDeleteRebalanceOutOfOrder(t *testing.T) {
	it := newTree(1)

	var toDelete *mockInterval
	for i := int64(9); i >= 0; i-- {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		it.add(iv)
		if i == 5 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 9)
	assert.Equal(t, 9, it.Len())
}

func TestDeleteRebalanceRandomOrder(t *testing.T) {
	it := newTree(1)

	starts := []int64{0, 4, 2, 1, 3}

	var toDelete *mockInterval
	for _, start := range starts {
		iv := constructSingleDimensionInterval(start, start+1, uint64(start))
		it.add(iv)
		if start == 1 {
			toDelete = iv
		}
	}

	it.Delete(toDelete)

	checkRedBlack(t, it.root, 1)
	result := it.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Len(t, result, 4)
	assert.Equal(t, 4, it.Len())
}

func TestDeleteEmptyTree(t *testing.T) {
	it := newTree(1)

	it.Delete(constructSingleDimensionInterval(0, 1, 1))

	assert.Equal(t, 0, it.Len())
}

func BenchmarkDeleteItems(b *testing.B) {
	numItems := int64(1000)

	intervals := make(Intervals, 0, numItems)
	for i := int64(0); i < numItems; i++ {
		iv := constructSingleDimensionInterval(i, i+1, uint64(i))
		intervals = append(intervals, iv)
	}

	trees := make([]*tree, 0, b.N)
	for i := 0; i < b.N; i++ {
		it := newTree(1)
		it.Add(intervals...)
		trees = append(trees, it)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trees[i].Delete(intervals...)
	}
}

func TestMax(t *testing.T) {
	it, _, _, _ := constructSingleDimensionQueryTestTree()

	assert.Equal(t, 12, it.Max(1))
	assert.Equal(t, 0, it.Max(2))

	it.root = nil

	assert.Equal(t, 0, it.Max(1))
}

func TestMin(t *testing.T) {
	it, _, _, _ := constructSingleDimensionQueryTestTree()

	assert.Equal(t, 4, it.Min(1))
	assert.Equal(t, 0, it.Min(2))

	it.root = nil

	assert.Equal(t, 0, it.Min(1))
}

func TestAddDuplicateRanges(t *testing.T) {
	it := newTree(1)
	iv1 := constructSingleDimensionInterval(0, 10, 1)
	iv2 := constructSingleDimensionInterval(0, 10, 2)
	iv3 := constructSingleDimensionInterval(0, 10, 3)

	it.Add(iv1, iv2, iv3)
	it.Delete(iv1, iv2, iv3)

	assert.Equal(t, 0, it.Len())
}

func TestAddDeleteDuplicatesRebalanceInOrder(t *testing.T) {
	it := newTree(1)

	intervals := make(Intervals, 0, 10)

	for i := 0; i < 10; i++ {
		iv := constructSingleDimensionInterval(0, 10, uint64(i))
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.Len())
}

func TestAddDeleteDuplicatesRebalanceReverseOrder(t *testing.T) {
	it := newTree(1)

	intervals := make(Intervals, 0, 10)

	for i := 9; i >= 0; i-- {
		iv := constructSingleDimensionInterval(0, 10, uint64(i))
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.Len())
}

func TestAddDeleteDuplicatesRebalanceRandomOrder(t *testing.T) {
	it := newTree(1)

	starts := []int{0, 4, 2, 1, 3}
	intervals := make(Intervals, 0, 5)

	for _, start := range starts {
		iv := constructSingleDimensionInterval(0, 10, uint64(start))
		intervals = append(intervals, iv)
	}

	it.Add(intervals...)
	it.Delete(intervals...)

	assert.Equal(t, 0, it.Len())
}

func TestInsertSingleAtDimension(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, 10, 1)
	assert.Len(t, deleted, 0)
	assert.Equal(t, ivs[1:], modified)

	result := tree.Query(constructSingleDimensionInterval(10, 20, 0))
	assert.Equal(t, ivs[1:], result)
	checkRedBlack(t, tree.root, 1)

	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 13, tree.root.max)
}

func TestInsertMultipleAtDimension(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, 10, 2)
	assert.Len(t, deleted, 0)
	assert.Equal(t, ivs[1:], modified)

	result := tree.Query(constructSingleDimensionInterval(11, 20, 0))
	assert.Equal(t, ivs[1:], result)
	checkRedBlack(t, tree.root, 1)

	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 14, tree.root.max)
}

func TestInsertAtLowestIndex(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, -1, 1)
	assert.Equal(t, ivs[0:], modified)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructSingleDimensionInterval(0, 1, 0))
	assert.Len(t, result, 0)

	result = tree.Query(constructSingleDimensionInterval(1, 4, 0))
	assert.Equal(t, ivs, result)

	checkRedBlack(t, tree.root, 1)

	assert.Equal(t, 1, tree.root.min)
	assert.Equal(t, 13, tree.root.max)
}

func TestDeleteSingleAtDimension(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, 10, -1)
	assert.Equal(t, ivs[1:], modified)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructSingleDimensionInterval(10, 20, 0))
	assert.Equal(t, ivs[2:], result)

	result = tree.Query(constructSingleDimensionInterval(9, 20, 0))
	assert.Equal(t, ivs, result)

	checkRedBlack(t, tree.root, 1)

	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 11, tree.root.max)
}

func TestDeleteBelowLowestIndex(t *testing.T) {
	tree := newTree(1)

	ivs := make(Intervals, 0, 3)
	for i := int64(0); i < 3; i++ {
		iv := constructSingleDimensionInterval(i+1, i+11, uint64(i))
		ivs = append(ivs, iv)
	}

	tree.Add(ivs...)

	modified, deleted := tree.Insert(1, 0, -1)
	assert.Equal(t, ivs, modified)
	assert.Len(t, deleted, 0)

	result := tree.Query(constructSingleDimensionInterval(0, 1, 0))
	assert.Equal(t, ivs[:1], result)

	result = tree.Query(constructSingleDimensionInterval(0, 10, 0))
	assert.Equal(t, ivs, result)

	checkRedBlack(t, tree.root, 1)
	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 12, tree.root.max)
}

func TestInsertDeletesInterval(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, 0, -10)
	assert.Equal(t, ivs[1:], modified)
	assert.Equal(t, ivs[:1], deleted)

	result := tree.Query(constructSingleDimensionInterval(2, 10, 0))
	assert.Len(t, result, 0)

	result = tree.Query(constructSingleDimensionInterval(0, 2, 0))
	assert.Equal(t, ivs[1:], result)

	checkRedBlack(t, tree.root, 1)
	assert.Equal(t, 2, tree.Len())
	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 2, tree.root.max)
}

func TestInsertInvalidDimension(t *testing.T) {
	tree, _ := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(2, 0, 1)
	assert.Len(t, deleted, 0)
	assert.Len(t, modified, 0)
}

func TestDeleteMiddleOfRange(t *testing.T) {
	tree, ivs := constructSingleDimensionTestTree(3)

	modified, deleted := tree.Insert(1, 5, -10)
	assert.Equal(t, ivs, modified)
	assert.Len(t, deleted, 0)

	checkRedBlack(t, tree.root, 1)
	assert.Equal(t, 0, tree.root.min)
	assert.Equal(t, 5, tree.root.max)
}

func BenchmarkInsertPositive(b *testing.B) {
	numItems := 1000

	tree, _ := constructSingleDimensionTestTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(1, 0, 1)
	}
}

func BenchmarkInsertNegative(b *testing.B) {
	numItems := 1000

	tree, _ := constructSingleDimensionTestTree(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(1, 0, int64(numItems))
	}
}
