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
