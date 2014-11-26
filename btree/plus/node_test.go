package plus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeInsert(t *testing.T) {
	tree := newBTree(3)
	n := newLeafNode(3)
	key := newMockKey(3, 3)

	n.insert(tree, key)

	assert.Len(t, n.keys, 1)
	assert.Nil(t, n.pointer)
	assert.Equal(t, n.keys[0].(*payload).keys[0], key)
	assert.Equal(t, 0, n.keys[0].Compare(key))
}

func TestDuplicateNodeInsert(t *testing.T) {
	tree := newBTree(3)
	n := newLeafNode(3)
	k1 := newMockKey(3, 3)
	k2 := newMockKey(3, 4)

	n.insert(tree, k1)
	n.insert(tree, k2)
	n.insert(tree, k1)

	assert.Len(t, n.keys, 1)
	assert.Nil(t, n.pointer)
	if !assert.Len(t, n.keys[0].(*payload).keys, 2) {
		return
	}
	assert.Equal(t, n.keys[0].(*payload).keys[0], k1)
	assert.Equal(t, n.keys[0].(*payload).keys[1], k2)
	assert.Equal(t, 0, n.keys[0].Compare(k1))
	assert.Equal(t, 0, n.keys[0].Compare(k2))
}

func TestMultipleNodeInsert(t *testing.T) {

}
