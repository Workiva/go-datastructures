package rangetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmutableSingleDimensionAdd(t *testing.T) {
	tree := newImmutableRangeTree(1)
	entry := constructMockEntry(0, int64(0), int64(0))
	tree2 := tree.Add(entry)

	result := tree.Query(
		constructMockInterval(dimension{0, 10}, dimension{0, 10}),
	)
	assert.Len(t, result, 0)

	result = tree2.Query(
		constructMockInterval(dimension{0, 10}, dimension{0, 10}),
	)
	assert.Equal(t, Entries{entry}, result)
}

func TestImmutableSingleDimensionMultipleAdds(t *testing.T) {
	tree := newImmutableRangeTree(1)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	tree1 := tree.Add(e1)
	tree2 := tree1.Add(e2)
	tree3 := tree2.Add(e3)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	result := tree1.Query(iv)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, 1, tree1.Len())

	result = tree2.Query(iv)
	assert.Equal(t, Entries{e1, e2}, result)
	assert.Equal(t, 2, tree2.Len())

	result = tree3.Query(iv)
	assert.Equal(t, Entries{e1, e2, e3}, result)
	assert.Equal(t, 3, tree3.Len())
}

func TestImmutableSingleDimensionBulkAdd(t *testing.T) {
	tree := newImmutableRangeTree(1)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	entries := Entries{e1, e2, e3}

	tree1 := tree.Add(entries...)

	result := tree1.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Equal(t, entries, result)
	assert.Equal(t, 3, tree1.Len())
}
