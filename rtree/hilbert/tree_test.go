package hilbert

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/rtree"
)

func TestSimpleInsert(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	tree := newTree(3)

	tree.Insert(r1)
	assert.Equal(t, uint64(1), tree.Len())

	q := newMockRectangle(5, 5, 15, 15)
	result := tree.Search(q)
	assert.Equal(t, []rtree.Rectangle{r1}, result)
}

func TestTwoInsert(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	tree := newTree(3)

	tree.Insert(r1, r2)
	assert.Equal(t, uint64(2), tree.Len())

	q := newMockRectangle(0, 0, 20, 20)
	result := tree.Search(q)
	assert.Equal(t, []rtree.Rectangle{r1, r2}, result)

	q = newMockRectangle(0, 0, 4, 4)
	result = tree.Search(q)
	assert.Equal(t, []rtree.Rectangle{r1}, result)

	q = newMockRectangle(11, 11, 20, 20)
	result = tree.Search(q)
	assert.Equal(t, []rtree.Rectangle{r2}, result)
}

func TestInsertCausesRootSplitOddAry(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	r3 := newMockRectangle(10, 10, 20, 20)
	tree := newTree(3)

	tree.Insert(r1, r2, r3)
	assert.Equal(t, uint64(3), tree.Len())

	q := newMockRectangle(0, 0, 20, 20)
	result := tree.Search(q)
	assert.Contains(t, result, r1)
	assert.Contains(t, result, r2)
	assert.Contains(t, result, r3)
}

func TestInsertCausesRootSplitEvenAry(t *testing.T) {
	r1 := newMockRectangle(0, 0, 10, 10)
	r2 := newMockRectangle(5, 5, 15, 15)
	r3 := newMockRectangle(10, 10, 20, 20)
	r4 := newMockRectangle(15, 15, 25, 25)
	tree := newTree(4)

	tree.Insert(r1, r2, r3, r4)
	assert.Equal(t, uint64(4), tree.Len())

	q := newMockRectangle(0, 0, 25, 25)
	result := tree.Search(q)
	assert.Contains(t, result, r1)
	assert.Contains(t, result, r2)
	assert.Contains(t, result, r3)
	assert.Contains(t, result, r4)
}
