package plus

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Print(`hate this`)
}

func TestSearchKeys(t *testing.T) {
	keys := keys{newMockKey(1, 1), newMockKey(2, 2), newMockKey(4, 4)}

	testKey := newMockKey(5, 5)
	assert.Equal(t, 3, keySearch(keys, testKey))

	testKey = newMockKey(2, 2)
	assert.Equal(t, 1, keySearch(keys, testKey))

	testKey = newMockKey(0, 0)
	assert.Equal(t, 0, keySearch(keys, testKey))

	testKey = newMockKey(3, 3)
	assert.Equal(t, 2, keySearch(keys, testKey))

	assert.Equal(t, 0, keySearch(nil, testKey))
}

func TestTreeInsert2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)

	tree.Insert(keys...)

	assert.Len(t, tree.root.(*inode).keys, 2)
	assert.Len(t, tree.root.(*inode).nodes, 3)
	assert.IsType(t, &inode{}, tree.root)
}

func TestTreeInsert3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	assert.Len(t, tree.root.(*inode).keys, 1)
	assert.Len(t, tree.root.(*inode).nodes, 2)
	assert.IsType(t, &inode{}, tree.root)
}

func TestTreeInsertQuery2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsertQuery3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsertReverseOrder2_3_4(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()
	keys.reverse() // we want to fetch things in the correct
	// ascending order

	assert.Equal(t, keys, result)
}

func TestTreeInsertReverseOrder3_4_5(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()
	keys.reverse() // we want to fetch things in the correct
	// ascending order

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5_WithEndDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	keys = append(keys, newMockKey(4, 5))
	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5_WithMiddleDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	key := newMockKey(2, 5)
	keys.insertAt(3, key)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5WithEarlyDuplicate(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	key := newMockKey(0, 5)
	keys.insertAt(1, key)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5WithDuplicateID(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	key := newMockKey(2, 2)
	tree.Insert(keys...)
	tree.Insert(key)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInsert3_4_5MiddleQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(2, 0))
	result := iter.exhaust()

	assert.Equal(t, keys[2:], result)
}

func TestTreeInsert3_4_5LateQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(4, 0))
	result := iter.exhaust()

	assert.Equal(t, keys[4:], result)
}

func TestTreeInsert3_4_5AfterQuery(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(5)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(5, 0))
	result := iter.exhaust()

	assert.Len(t, result, 0)
}

func TestTreeInternalNodeSplit(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(10)

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Equal(t, keys, result)
}

func TestTreeInternalNodeSplitReverseOrder(t *testing.T) {
	tree := newBTree(4)
	keys := constructMockKeys(10)
	keys.reverse()

	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()
	keys.reverse()

	assert.Equal(t, keys, result)
}

func TestTreeInternalNodeSplitRandomOrder(t *testing.T) {
	ids := []uint64{6, 2, 9, 0, 3, 4, 7, 1, 8, 5}
	keys := make(keys, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, newMockKey(id, id))
	}

	tree := newBTree(4)
	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(0, 0))
	result := iter.exhaust()

	assert.Len(t, result, 10)
	for i, key := range result {
		assert.Equal(t, newMockKey(uint64(i), uint64(i)), key)
	}
}

func TestTreeRandomOrderQuery(t *testing.T) {
	ids := []uint64{6, 2, 9, 0, 3, 4, 7, 1, 8, 5}
	keys := make(keys, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, newMockKey(id, id))
	}

	tree := newBTree(4)
	tree.Insert(keys...)

	iter := tree.Iterate(newMockKey(4, 4))
	result := iter.exhaust()

	assert.Len(t, result, 6)
	for i, key := range result {
		assert.Equal(t, newMockKey(uint64(i)+4, uint64(i)+4), key)
	}
}
