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

	g.ParallelRecursivelyApply(func(INode) {
		atomic.AddInt64(&numCalls, 1)
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

	g.ParallelRecursivelyApply(func(INode) {
		atomic.AddInt64(&numCalls, 1)
	})

	if numCalls != int64(numItems) {
		t.Errorf(`Expected calls: %d, received: %d`, numItems, numCalls)
	}
}

func TestRecursivelyApply(t *testing.T) {
	g := newExecutionGraph()

	g.toApply = []Nodes{Nodes{newTestNode(0), newTestNode(1)}}

	numCalls := 0

	g.RecursivelyApply(func(node INode) {
		numCalls++
	})

	if numCalls != 2 {
		t.Errorf(`Expected num calls: %d, received: %d`, 2, numCalls)
	}
}

func TestApplyWithCirculars(t *testing.T) {
	g := newExecutionGraph()

	g.toApply = []Nodes{Nodes{newTestNode(0), newTestNode(1)}}
	g.circulars = Nodes{newTestNode(2)}

	numCalls := 0

	g.RecursivelyApply(func(node INode) {
		numCalls++
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

	g.ParallelRecursivelyApply(func(INode) {
		atomic.AddInt64(&numCalls, 1)
	})

	assert.Equal(t, expectedNumCalls, numCalls)
}

func TestCirculars(t *testing.T) {
	g := newExecutionGraph()
	nodes := Nodes{newTestNode(1)}
	g.circulars = nodes

	assert.Equal(t, nodes, g.Circulars())
}

func TestLayer(t *testing.T) {
	eg := newExecutionGraph()
	n1 := Nodes{newTestNode(1)}
	n2 := Nodes{newTestNode(2)}

	eg.toApply = []Nodes{n1, n2}

	result, err := eg.Layer(0)
	assert.Equal(t, n1, result)
	assert.Nil(t, err)

	result, err = eg.Layer(2)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestNumLayers(t *testing.T) {
	eg := newExecutionGraph()
	n1 := Nodes{newTestNode(1)}
	n2 := Nodes{newTestNode(2)}

	assert.Equal(t, 0, eg.NumLayers())

	eg.toApply = []Nodes{n1, n2}

	assert.Equal(t, 2, eg.NumLayers())
}

func TestApplyLayer(t *testing.T) {
	eg := newExecutionGraph()
	n0 := newTestNode(0)
	n1 := newTestNode(1)

	eg.toApply = []Nodes{Nodes{n0}, Nodes{n1}}

	var called INode

	err := eg.ParallelApplyLayer(0, func(n INode) {
		called = n
	})

	assert.Nil(t, err)
	assert.Equal(t, n0, called)
	called = nil

	err = eg.ParallelApplyLayer(3, func(n INode) {
		called = n
	})

	assert.Nil(t, called)
	assert.NotNil(t, err)
}
