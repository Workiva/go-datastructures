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
