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

func TestImmutableMultiDimensionAdd(t *testing.T) {
	tree := newImmutableRangeTree(2)
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

func TestImmutableMultiDimensionMultipleAdds(t *testing.T) {
	tree := newImmutableRangeTree(2)
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

func TestImmutableMultiDimensionBulkAdd(t *testing.T) {
	tree := newImmutableRangeTree(2)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	entries := Entries{e1, e2, e3}

	tree1 := tree.Add(entries...)

	result := tree1.Query(constructMockInterval(dimension{0, 10}, dimension{0, 10}))
	assert.Equal(t, entries, result)
	assert.Equal(t, 3, tree1.Len())
}

func BenchmarkImmutableMultiDimensionInserts(b *testing.B) {
	numItems := int64(1000)

	entries := make(Entries, 0, numItems)
	for i := int64(0); i < numItems; i++ {
		e := constructMockEntry(uint64(i), i, i)
		entries = append(entries, e)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newImmutableRangeTree(2)
		for _, e := range entries {
			tree = tree.Add(e)
		}
	}
}

func BenchmarkImmutableMultiDimensionBulkInsert(b *testing.B) {
	numItems := int64(100000)

	entries := make(Entries, 0, numItems)
	for i := int64(0); i < numItems; i++ {
		e := constructMockEntry(uint64(i), i, i)
		entries = append(entries, e)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newImmutableRangeTree(2)
		tree.Add(entries...)
	}
}

func BenchmarkMultiDimensionBulkInsert(b *testing.B) {
	numItems := int64(100000)

	entries := make(Entries, 0, numItems)
	for i := int64(0); i < numItems; i++ {
		e := constructMockEntry(uint64(i), i, i)
		entries = append(entries, e)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newOrderedTree(2)
		tree.Add(entries...)
	}
}

func TestImmutableSingleDimensionDelete(t *testing.T) {
	tree := newImmutableRangeTree(1)
	entry := constructMockEntry(0, int64(0), int64(0))
	tree2 := tree.Add(entry)
	tree3 := tree2.Delete(entry)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	result := tree3.Query(iv)
	assert.Len(t, result, 0)
}

func TestImmutableSingleDimensionMultipleDeletes(t *testing.T) {
	tree := newImmutableRangeTree(1)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	tree1 := tree.Add(e1)
	tree2 := tree1.Add(e2)
	tree3 := tree2.Add(e3)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	tree4 := tree3.Delete(e3)
	result := tree4.Query(iv)
	assert.Equal(t, Entries{e1, e2}, result)
	assert.Equal(t, 2, tree4.Len())

	tree5 := tree4.Delete(e2)
	result = tree5.Query(iv)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, 1, tree5.Len())

	tree6 := tree5.Delete(e1)
	result = tree6.Query(iv)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, tree6.Len())

	result = tree3.Query(iv)
	assert.Equal(t, Entries{e1, e2, e3}, result)
	assert.Equal(t, 3, tree3.Len())

	tree7 := tree3.Delete(constructMockEntry(0, int64(3), int64(3)))
	assert.Equal(t, tree3, tree7)
}

func TestImmutableSingleDimensionBulkDeletes(t *testing.T) {
	tree := newImmutableRangeTree(1)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	tree1 := tree.Add(e1, e2, e3)
	tree2 := tree1.Delete(e2, e3)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	result := tree2.Query(iv)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, 1, tree2.Len())

	tree3 := tree2.Delete(e1)

	result = tree3.Query(iv)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, tree3.Len())
}

func TestImmutableMultiDimensionDelete(t *testing.T) {
	tree := newImmutableRangeTree(2)
	entry := constructMockEntry(0, int64(0), int64(0))
	tree2 := tree.Add(entry)
	tree3 := tree2.Delete(entry)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	result := tree3.Query(iv)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, tree3.Len())
}

func TestImmutableMultiDimensionMultipleDeletes(t *testing.T) {
	tree := newImmutableRangeTree(2)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	tree1 := tree.Add(e1)
	tree2 := tree1.Add(e2)
	tree3 := tree2.Add(e3)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})
	tree4 := tree3.Delete(e3)

	result := tree4.Query(iv)
	assert.Equal(t, Entries{e1, e2}, result)
	assert.Equal(t, 2, tree4.Len())

	tree5 := tree4.Delete(e2)
	result = tree5.Query(iv)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, 1, tree5.Len())

	tree6 := tree5.Delete(e1)
	result = tree6.Query(iv)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, tree6.Len())

	result = tree3.Query(iv)
	assert.Equal(t, Entries{e1, e2, e3}, result)
	assert.Equal(t, 3, tree3.Len())

	tree7 := tree3.Delete(constructMockEntry(0, int64(3), int64(3)))
	assert.Equal(t, tree3, tree7)
}

func TestImmutableMultiDimensionBulkDeletes(t *testing.T) {
	tree := newImmutableRangeTree(2)
	e1 := constructMockEntry(0, int64(0), int64(0))
	e2 := constructMockEntry(0, int64(1), int64(1))
	e3 := constructMockEntry(0, int64(2), int64(2))

	tree1 := tree.Add(e1, e2, e3)
	tree2 := tree1.Delete(e2, e3)

	iv := constructMockInterval(dimension{0, 10}, dimension{0, 10})

	result := tree2.Query(iv)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, 1, tree2.Len())

	tree3 := tree2.Delete(e1)

	result = tree3.Query(iv)
	assert.Len(t, result, 0)
	assert.Equal(t, 0, tree3.Len())
}

func constructMultiDimensionalImmutableTree(number int64) (*immutableRangeTree, Entries) {
	tree := newImmutableRangeTree(2)
	entries := make(Entries, 0, number)
	for i := int64(0); i < number; i++ {
		entries = append(entries, constructMockEntry(uint64(i), i, i))
	}

	return tree.Add(entries...), entries
}

func TestImmutableInsertPositiveIndexFirstDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalImmutableTree(2)

	tree1, modified, deleted := tree.InsertAtDimension(1, 1, 1)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree1.Query(constructMockInterval(dimension{2, 10}, dimension{1, 10}))
	assert.Equal(t, entries[1:], result)

	result = tree.Query(constructMockInterval(dimension{2, 10}, dimension{0, 10}))
	assert.Len(t, result, 0)
}

func TestImmutableInsertPositiveIndexSecondDimension(t *testing.T) {
	tree, entries := constructMultiDimensionalImmutableTree(3)

	tree1, modified, deleted := tree.InsertAtDimension(2, 1, 1)
	assert.Len(t, deleted, 0)
	assert.Equal(t, entries[1:], modified)

	result := tree1.Query(constructMockInterval(dimension{1, 10}, dimension{2, 10}))
	assert.Equal(t, entries[1:], result)

	result = tree.Query(constructMockInterval(dimension{1, 10}, dimension{2, 10}))
	assert.Equal(t, entries[2:], result)
}
