package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/bitarray"
)

func checkLayers(tb testing.TB, expected, result []Nodes) {
	if len(expected) != len(result) {
		tb.Errorf(
			`Expected len: %d, received: %d`, len(expected), len(result),
		)
		return
	}

	for i, nodes := range expected {
		checkNodes(tb, result[i], nodes)
	}
}

func TestGetSimpleSubgraph(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 2

	n1 := newTestNode(0)
	n2 := newTestNode(1)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(n1.ID())
			return []bitarray.BitArray{ba}
		} else {
			return []bitarray.BitArray{ba}
		}
	}

	g := FromNodes(dp, Nodes{n1, n2})
	eg := g.GetSubgraph(dp, Nodes{n1, n2})

	assert.Equal(t, 2, eg.Size())
	checkLayers(t, []Nodes{Nodes{n1}, Nodes{n2}}, eg.toApply)
}

func TestGetSubgraphSingleLayer(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 2

	n1 := newTestNode(0)
	n2 := newTestNode(1)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		return []bitarray.BitArray{bitarray.NewBitArray(dp.MaxNode())}
	}

	g := FromNodes(dp, Nodes{n1, n2})
	eg := g.GetSubgraph(dp, Nodes{n1, n2})

	assert.Equal(t, 2, eg.Size())
	if !assert.Len(t, eg.toApply, 1) {
		return
	}
	assert.Contains(t, eg.toApply[0], n1)
	assert.Contains(t, eg.toApply[0], n2)
}

func TestGetSubgraphExternalDependency(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 2

	n1 := newTestNode(0)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		ba.SetBit(1)
		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1})
	eg := g.GetSubgraph(dp, Nodes{n1})

	assert.Equal(t, 1, eg.Size())
	checkLayers(t, []Nodes{Nodes{n1}}, eg.toApply)
}

func TestGetSubgraphWithAllCircularDependencies(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 2

	n1 := newTestNode(0)
	n2 := newTestNode(1)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n1 {
			ba.SetBit(1)
		} else {
			ba.SetBit(0)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2})
	eg := g.GetSubgraph(dp, Nodes{n1, n2})

	assert.Equal(t, 2, eg.Size())
	assert.Len(t, eg.circulars, 2)
	assert.Contains(t, eg.circulars, n1)
	assert.Contains(t, eg.circulars, n2)
	assert.True(t, n1.IsCircular())
	assert.True(t, n2.IsCircular())
}

func TestGetSubgraphWithMixedCirculars(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 3

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(2)
		} else if node == n3 {
			ba.SetBit(1)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2, n3})
	eg := g.GetSubgraph(dp, Nodes{n1, n2, n3})

	assert.Equal(t, 3, eg.Size())
	checkLayers(t, []Nodes{Nodes{n1}}, eg.toApply)
	assert.Len(t, eg.circulars, 2)
	assert.Contains(t, eg.circulars, n2)
	assert.Contains(t, eg.circulars, n3)
	assert.True(t, n2.IsCircular())
	assert.True(t, n3.IsCircular())
}

func TestGetSubgraphResetsCircularState(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 2

	n1 := newTestNode(0)
	n1.SetCircular(true)
	n2 := newTestNode(1)
	n2.SetCircular(true)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		return []bitarray.BitArray{bitarray.NewBitArray(dp.MaxNode())}
	}

	g := FromNodes(dp, Nodes{n1, n2})
	eg := g.GetSubgraph(dp, Nodes{n1, n2})

	assert.Equal(t, 2, eg.Size())
	if !assert.Len(t, eg.toApply, 1) {
		return
	}
	assert.Contains(t, eg.toApply[0], n1)
	assert.Contains(t, eg.toApply[0], n2)
	assert.False(t, n1.IsCircular())
	assert.False(t, n2.IsCircular())
}

func TestGetSubgraphResetsMixedCircular(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 3

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n2.SetCircular(true)
	n3 := newTestNode(2)
	n3.SetCircular(true)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2, n3})
	eg := g.GetSubgraph(dp, Nodes{n1, n2, n3})

	assert.Equal(t, 3, eg.Size())
	checkLayers(t, []Nodes{Nodes{n1}, Nodes{n2}, Nodes{n3}}, eg.toApply)
	assert.False(t, n2.IsCircular())
	assert.False(t, n3.IsCircular())
}

func BenchmarkGetSubgraph(b *testing.B) {
	numItems := uint64(1000)

	nodes := make(Nodes, 0, numItems)
	dependencyMap := make(map[uint64][]bitarray.BitArray, numItems)
	dp := newTestDependencyProvider()
	dp.maxNode = numItems
	dp.baFactory = func(node INode) []bitarray.BitArray {
		return dependencyMap[node.ID()]
	}

	for i := uint64(0); i < numItems; i++ {
		n := newTestNode(i)
		ba := bitarray.NewSparseBitArray()

		for j := uint64(0); j < i; j++ {
			ba.SetBit(j)
		}

		dependencyMap[n.ID()] = []bitarray.BitArray{ba}
		nodes = append(nodes, n)
	}

	g := FromNodes(dp, nodes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.GetSubgraph(dp, nodes)
	}
}

func TestPositionsFlatten(t *testing.T) {
	f1 := newTestNode(0)
	f2 := newTestNode(1)
	f3 := newTestNode(2)
	f4 := newTestNode(3)
	positions := positions{
		&bundle{INode: f1, position: -1},
		&bundle{INode: f2, position: 2},
		&bundle{INode: f3, position: 4},
		&bundle{INode: f4, position: 3},
	}

	layers, circulars := positions.flatten(Nodes{f1, f2, f3, f4})
	assert.Equal(t, Nodes{f1}, circulars)
	assert.Len(t, layers, 3)
	l := layers[0]
	assert.Equal(t, &layer{position: 2, nodes: Nodes{f2}}, l)
	l = layers[1]
	assert.Equal(t, &layer{position: 3, nodes: Nodes{f4}}, l)
	l = layers[2]
	assert.Equal(t, &layer{position: 4, nodes: Nodes{f3}}, l)

	layers, circulars = positions.flatten(Nodes{f2, nil})
	assert.Len(t, circulars, 0)
	assert.Len(t, layers, 1)
}

func TestGraphFromNodes(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
			ba.SetBit(3)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	graph := FromNodes(dp, Nodes{n1, n2, n3, n4})
	positions := positions{
		&bundle{INode: n1, position: 0},
		&bundle{INode: n2, position: 1},
		&bundle{INode: n3, position: -1},
		&bundle{INode: n4, position: -1},
	}
	assert.Equal(t, positions, graph.positions)

	graph = FromNodes(dp, nil)
	assert.Len(t, graph.positions, 0)
}

// Ensures we examine the graph's position index to return
// correct subgraph layering.
func TestGetEntireSubgraphFromGraph(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
			ba.SetBit(3)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	graph := FromNodes(dp, Nodes{n1, n2, n3, n4})

	eg := graph.GetSubgraph(dp, Nodes{n1, n2, n3, n4})
	if !assert.Len(t, eg.toApply, 2) {
		return
	}

	assert.Equal(t, Nodes{n1}, eg.toApply[0])
	assert.Equal(t, Nodes{n2}, eg.toApply[1])
	assert.Contains(t, eg.circulars, n3)
	assert.Contains(t, eg.circulars, n4)
	assert.Len(t, eg.circulars, 2)
	assert.Equal(t, 4, eg.size)
}

func TestGetSmallSubgraphFromGraph(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
			ba.SetBit(3)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	graph := FromNodes(dp, Nodes{n1, n2, n3, n4})

	eg := graph.GetSubgraph(dp, Nodes{n1})
	if !assert.Len(t, eg.toApply, 1) {
		return
	}

	assert.Equal(t, Nodes{n1}, eg.toApply[0])
	assert.Equal(t, 1, eg.size)
	assert.Len(t, eg.circulars, 0)
}

func TestAddNodesWithNoDependents(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
			ba.SetBit(3)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, nil)
	eg := g.AddNodes(dp, Nodes{n1, n2, n3, n4})

	positions := positions{
		&bundle{INode: n1, position: 0},
		&bundle{INode: n2, position: 1},
		&bundle{INode: n3, position: -1},
		&bundle{INode: n4, position: -1},
	}
	assert.Equal(t, positions, g.positions)
	assert.Equal(t, 4, eg.size)
	if !assert.Len(t, eg.toApply, 2) {
		return
	}
	assert.Equal(t, Nodes{n1}, eg.toApply[0])
	assert.Equal(t, Nodes{n2}, eg.toApply[1])
	if !assert.Len(t, eg.circulars, 2) {
		return
	}
	assert.Contains(t, eg.circulars, n3)
	assert.Contains(t, eg.circulars, n4)
}

func TestAddNodesWithDependents(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		return []bitarray.BitArray{bitarray.NewBitArray(dp.maxNode)}
	}

	g := FromNodes(dp, Nodes{n2})
	dp.dependents = Nodes{n2}
	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.MaxNode())
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
			ba.SetBit(3)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	eg := g.AddNodes(dp, Nodes{n1, n3, n4})
	positions := positions{
		&bundle{INode: n1, position: 0},
		&bundle{INode: n2, position: 1},
		&bundle{INode: n3, position: -1},
		&bundle{INode: n4, position: -1},
	}
	assert.Equal(t, positions, g.positions)
	assert.Equal(t, 4, eg.size)
	if !assert.Len(t, eg.toApply, 2) {
		return
	}
	assert.Equal(t, Nodes{n1}, eg.toApply[0])
	assert.Equal(t, Nodes{n2}, eg.toApply[1])
	assert.Len(t, eg.circulars, 2)
	assert.Contains(t, eg.circulars, n3)
	assert.Contains(t, eg.circulars, n4)
}

func TestAddNodesWithCircularDependents(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.maxNode)
		if node == n1 {
			ba.SetBit(1)
		} else {
			ba.SetBit(0)
		}
		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2})
	assert.Equal(t, -1, g.positions[0].position)
	dp.dependents = Nodes{n2}
	eg := g.AddNodes(dp, Nodes{n3})
	assert.Len(t, eg.toApply, 1)
}

func TestPositionsExtractLayers(t *testing.T) {
	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	positions := positions{
		&bundle{INode: n1, position: -1},
		&bundle{INode: n2, position: 3},
		&bundle{INode: n3, position: 3},
		&bundle{INode: n4, position: 4},
	}

	nodes := positions.extractLayers(map[int64]bool{-1: true, 3: true})
	expected := Nodes{n1, n2, n3}
	assert.Equal(t, expected, nodes)
}

func TestPositionsLowest(t *testing.T) {
	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	positions := positions{
		&bundle{INode: n1, position: -1},
		&bundle{INode: n2, position: 3},
		&bundle{INode: n3, position: 3},
		&bundle{INode: n4, position: 4},
	}

	lowest := positions.lowest(Nodes{n3, n4})
	assert.Equal(t, 3, lowest)

	lowest = positions.lowest(Nodes{n2, n1})
	assert.Equal(t, -1, lowest)
}

func TestRemoveNodes(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.maxNode)
		if node == n2 || node == n3 {
			ba.SetBit(0)
		} else if node == n4 {
			ba.SetBit(1)
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2, n3, n4})
	assert.Equal(t, 2, g.maxLayer)

	eg := g.RemoveNodes(dp, Nodes{n3})
	assert.Equal(t, 2, g.maxLayer)
	positions := positions{
		&bundle{INode: n1, position: 0},
		&bundle{INode: n2, position: 1},
		nil,
		&bundle{INode: n4, position: 2},
	}
	assert.Equal(t, positions, g.positions)
	if !assert.Equal(t, 2, eg.Size()) {
		return
	}

	assert.Equal(t, Nodes{n2}, eg.toApply[0])
	assert.Equal(t, Nodes{n4}, eg.toApply[1])
	assert.Len(t, eg.circulars, 0)
}

func TestRemoveNodesWithCircular(t *testing.T) {
	dp := newTestDependencyProvider()
	dp.maxNode = 4

	n1 := newTestNode(0)
	n2 := newTestNode(1)
	n3 := newTestNode(2)
	n4 := newTestNode(3)

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.maxNode)
		if node == n1 {
			ba.SetBit(3)
		} else if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
		} else if node == n4 {
			ba.SetBit(2)
		}

		return []bitarray.BitArray{ba}
	}

	g := FromNodes(dp, Nodes{n1, n2, n3, n4})
	assert.Equal(t, -1, g.maxLayer)
	assert.Equal(t, -1, g.positions.highestSeen())

	dp.baFactory = func(node INode) []bitarray.BitArray {
		ba := bitarray.NewBitArray(dp.maxNode)
		if node == n2 {
			ba.SetBit(0)
		} else if node == n3 {
			ba.SetBit(1)
		}

		return []bitarray.BitArray{ba}
	}

	eg := g.RemoveNodes(dp, Nodes{n4})
	positions := positions{
		&bundle{INode: n1, position: 0},
		&bundle{INode: n2, position: 1},
		&bundle{INode: n3, position: 2},
		nil,
	}
	assert.Equal(t, positions, g.positions)
	assert.Equal(t, 2, g.maxLayer)
	if !assert.Len(t, eg.toApply, 3) {
		return
	}

	assert.Equal(t, Nodes{n1}, eg.toApply[0])
	assert.Equal(t, Nodes{n2}, eg.toApply[1])
	assert.Equal(t, Nodes{n3}, eg.toApply[2])
	assert.Equal(t, 3, eg.size)
	assert.Len(t, eg.circulars, 0)
}

func BenchmarkPreIndexedFlattening(b *testing.B) {
	numItems := uint64(1000)

	nodes := make(Nodes, 0, numItems)
	dependencyMap := make(map[uint64][]bitarray.BitArray, numItems)
	dp := newTestDependencyProvider()
	dp.maxNode = numItems
	dp.baFactory = func(node INode) []bitarray.BitArray {
		return dependencyMap[node.ID()]
	}

	for i := uint64(0); i < numItems; i++ {
		n := newTestNode(i)
		ba := bitarray.NewSparseBitArray()

		for j := uint64(0); j < i; j++ {
			ba.SetBit(j)
		}

		dependencyMap[n.ID()] = []bitarray.BitArray{ba}
		nodes = append(nodes, n)
	}

	g := FromNodes(dp, nodes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.positions.flatten(nodes)
	}
}

func BenchmarkFlattening(b *testing.B) {
	numItems := uint64(1000)

	nodes := make(Nodes, 0, numItems)
	dependencyMap := make(map[uint64][]bitarray.BitArray, numItems)
	dp := newTestDependencyProvider()
	dp.maxNode = numItems
	dp.baFactory = func(node INode) []bitarray.BitArray {
		return dependencyMap[node.ID()]
	}

	for i := uint64(0); i < numItems; i++ {
		n := newTestNode(i)
		ba := bitarray.NewSparseBitArray()

		for j := uint64(0); j < i; j++ {
			ba.SetBit(j)
		}

		dependencyMap[n.ID()] = []bitarray.BitArray{ba}
		nodes = append(nodes, n)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		flatten(dp, nodes)
	}
}
