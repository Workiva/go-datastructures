package rangetree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedAdd(t *testing.T) {
	nodes := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	nodes.add(n1)
	nodes.add(n2)

	assert.Equal(t, orderedNodes{n2, n1}, nodes)
}

func TestOrderedDelete(t *testing.T) {
	nodes := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	nodes.add(n1)
	nodes.add(n2)

	nodes.delete(n2.value)

	assert.Equal(t, orderedNodes{n1}, nodes)

	nodes.delete(n1.value)

	assert.Len(t, nodes, 0)
}

func TestApply(t *testing.T) {
	ns := make(orderedNodes, 0)

	n1 := newNode(4, constructMockEntry(1, 4), false)
	n2 := newNode(1, constructMockEntry(2, 1), false)

	ns.add(n1)
	ns.add(n2)

	results := make(nodes, 0, 2)

	ns.apply(1, 2, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n2}, results)

	results = results[:0]

	ns.apply(0, 1, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(2, 4, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(4, 5, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n1}, results)
	results = results[:0]

	ns.apply(0, 5, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Equal(t, nodes{n2, n1}, results)
	results = results[:0]

	ns.apply(5, 10, func(n *node) bool {
		results = append(results, n)
		return true
	})

	assert.Len(t, results, 0)
	results = results[:0]

	ns.apply(0, 100, func(n *node) bool {
		results = append(results, n)
		return false
	})

	assert.Equal(t, nodes{n2}, results)
}
