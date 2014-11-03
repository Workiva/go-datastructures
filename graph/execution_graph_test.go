package graph

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelRecursivelyApply(t *testing.T) {
	numItems := numNodesBeforeSplit + 1 // to ensure this happens in parallel
	g := newExecutionGraph()

	g.toApply = make([]Nodes, 0, 2)
	g.size = numItems * 2
	var numNodes int64
	var expectedNumCalls int64

	for i := 0; i < 2; i++ {
		if i == 0 {
			numNodes = numItems
		} else {
			numNodes = numItems - 2
		}
		nodes := make(Nodes, 0, numItems)
		for i := int64(0); i < numNodes; i++ {
			node := newTestNode(uint64(i))
			nodes = append(nodes, node)
		}
		g.toApply = append(g.toApply, nodes)
		expectedNumCalls += int64(numNodes)
	}

	numCalls := int64(0)

	g.ParallelRecursivelyApply(func(INode) bool {
		atomic.AddInt64(&numCalls, 1)
		return true
	})

	if numCalls != expectedNumCalls {
		t.Errorf(`Expected calls: %d, received: %d`, expectedNumCalls, numCalls)
	}
}

func TestParallelRecursivelyApplyOneLayer(t *testing.T) {
	numItems := numNodesBeforeSplit - 1 // to ensure this happens in parallel
	g := newExecutionGraph()

	nodes := make(Nodes, 0, numItems)
	for i := int64(0); i < numItems; i++ {
		node := newTestNode(uint64(i))
		nodes = append(nodes, node)
	}

	g.toApply = []Nodes{nodes}
	numCalls := int64(0)

	g.ParallelRecursivelyApply(func(INode) bool {
		atomic.AddInt64(&numCalls, 1)
		return true
	})

	if numCalls != int64(numItems) {
		t.Errorf(`Expected calls: %d, received: %d`, numItems, numCalls)
	}
}

func TestRecursivelyApply(t *testing.T) {
	g := newExecutionGraph()

	g.toApply = []Nodes{Nodes{newTestNode(0), newTestNode(1)}}

	numCalls := 0

	g.RecursivelyApply(func(node INode) bool {
		numCalls++
		return true
	})

	if numCalls != 2 {
		t.Errorf(`Expected num calls: %d, received: %d`, 2, numCalls)
	}

	numCalls = 0
	g.RecursivelyApply(func(node INode) bool {
		numCalls++
		return false
	})

	if numCalls != 1 {
		t.Errorf(`Expected num calls: %d, received: %d`, 1, numCalls)
	}
}

func TestApplyWithCirculars(t *testing.T) {
	g := newExecutionGraph()

	g.toApply = []Nodes{Nodes{newTestNode(0), newTestNode(1)}}
	g.circulars = Nodes{newTestNode(2)}

	numCalls := 0

	g.RecursivelyApply(func(node INode) bool {
		numCalls++
		return true
	})

	assert.Equal(t, 3, numCalls)
}

func TestParallelApplyWithCirculars(t *testing.T) {
	numItems := numNodesBeforeSplit + 1 // to ensure this happens in parallel
	g := newExecutionGraph()

	g.toApply = make([]Nodes, 0, 2)
	g.size = numItems * 2
	var numNodes int64
	var expectedNumCalls int64

	for i := 0; i < 2; i++ {
		if i == 0 {
			numNodes = numItems
		} else {
			numNodes = numItems - 2
		}
		nodes := make(Nodes, 0, numItems)
		for i := int64(0); i < numNodes; i++ {
			node := newTestNode(uint64(i))
			nodes = append(nodes, node)
		}
		g.toApply = append(g.toApply, nodes)
		expectedNumCalls += int64(numNodes)
	}

	g.circulars = Nodes{newTestNode(0)}
	expectedNumCalls += 1

	numCalls := int64(0)

	g.ParallelRecursivelyApply(func(INode) bool {
		atomic.AddInt64(&numCalls, 1)
		return true
	})

	assert.Equal(t, expectedNumCalls, numCalls)
}
