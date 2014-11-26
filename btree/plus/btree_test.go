package plus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestTreeInsert(t *testing.T) {
	tree := newBTree(3)
	keys := constructMockKeys(4)

	tree.Insert(keys...)

	assert.Len(t, tree.root.(*inode).keys, 2)
	assert.Len(t, tree.root.(*inode).nodes, 3)
	assert.IsType(t, &inode{}, tree.root)
}
