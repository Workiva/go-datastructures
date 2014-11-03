package graph

import (
	"testing"
)

func checkNodes(tb testing.TB, result, expected Nodes) {
	if len(result) != len(expected) {
		tb.Errorf(`Expected len: %d, received: %d`, len(expected), len(result))
		return // prevent panic
	}

	for i, node := range result {
		if node == nil {
			if expected[i] != nil {
				tb.Errorf(`Expected nil at: %d`, i)
			}
			continue
		}
		if !expected.Exists(node) {
			tb.Errorf(`Expected node: %+v not found.`, node)
		}
	}
}

func TestRemove(t *testing.T) {
	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)

	nodes := Nodes{n1, n2, n3}

	nodes = nodes.Remove(n1, n3)

	checkNodes(t, nodes, Nodes{n2})

	nodes = nodes.Remove(n2, n3)
	checkNodes(t, nodes, Nodes{})
}

func TestExists(t *testing.T) {
	n1 := newTestNode(0)
	n2 := newTestNode(1)

	nodes := Nodes{n1}

	if !nodes.Exists(n1) {
		t.Errorf(`Node does exist.`)
	}

	if nodes.Exists(n2) {
		t.Errorf(`Node does not exist.`)
	}
}

func TestNodesConstructor(t *testing.T) {
	n1 := newTestNode(0)
	n2 := newTestNode(1)

	nodes := NewNodes(n1, n2)

	checkNodes(t, nodes, Nodes{n1, n2})
}

func TestNodesDispose(t *testing.T) {
	n1 := newTestNode(0)

	nodes := NewNodes(n1)
	nodes.Dispose()

	if len(nodes) != 0 {
		t.Errorf(`Expected len: %d, received: %d`, 0, len(nodes))
	}
}
