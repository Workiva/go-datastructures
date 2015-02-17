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
	t.Logf(`TREE.ROOT: %+v`, tree.root)
}
