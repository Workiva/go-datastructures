package graph

import "github.com/Workiva/go-datastructures/bitarray"

type testDependencyProvider struct {
	baFactory  func(node INode) []bitarray.BitArray
	maxNode    uint64
	dependents Nodes
}

func (provider *testDependencyProvider) GetDependencies(node INode) []bitarray.BitArray {
	return provider.baFactory(node)
}

func (provider *testDependencyProvider) MaxNode() uint64 {
	return provider.maxNode
}

func (provider *testDependencyProvider) GetDependents(nodes Nodes) Nodes {
	return provider.dependents
}

func newTestDependencyProvider() *testDependencyProvider {
	return &testDependencyProvider{}
}

type testNode struct {
	isCircular bool
	id         uint64
}

func (node *testNode) SetCircular(value bool) {
	node.isCircular = value
}

func (node *testNode) IsCircular() bool {
	return node.isCircular
}

func (node *testNode) ID() uint64 {
	return node.id
}

func newTestNode(id uint64) *testNode {
	return &testNode{id: id}
}
