package skip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeSearch(t *testing.T) {
	n0 := newNode(newMockEntry(1), 1)
	n1 := newNode(newMockEntry(1), 1)
	n2 := newNode(newMockEntry(5), 1)
	n3 := newNode(newMockEntry(5), 1)
	n4 := newNode(newMockEntry(5), 1)
	n5 := newNode(newMockEntry(10), 1)

	nodes := nodes{n0, n1, n2, n3, n4, n5, nil}
	low, high := 0, len(nodes)-1

	assert.Equal(t, -1, nodes.search(0, low, high))
	assert.Equal(t, 1, nodes.search(1, low, high))
	assert.Equal(t, 1, nodes.search(2, low, high))
	assert.Equal(t, 4, nodes.search(6, low, high))
	assert.Equal(t, 5, nodes.search(10, low, high))
	assert.Equal(t, 5, nodes.search(11, low, high))
	assert.Equal(t, 4, nodes.search(5, low, high))
}
