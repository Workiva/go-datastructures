package graph

import (
	"runtime"
	"sync"
)

var nodesPool = sync.Pool{
	New: func() interface{} {
		return make(Nodes, 0, 10) // 10 is a guess
	},
}

// Nodes is an alias for a list of nodes with some helpful
// utility functions.
type Nodes []INode

// Remove the given items from the list.
func (nodes Nodes) Remove(items ...INode) Nodes {
	indexes := make([]int64, 0, len(nodes))
	for i, dep := range nodes {
		for _, node := range items {
			if dep == node {
				indexes = append(indexes, int64(i-len(indexes)))
				break
			}
		}

		if len(indexes) == len(nodes) {
			break
		}
	}

	for _, i := range indexes {
		nodes = nodes[:i+int64(copy(nodes[i:], nodes[i+1:]))]
	}

	return nodes
}

// Split splits the list of nodes into smaller and evenly sized chunks.
func (nodes Nodes) Split() []Nodes {
	numParts := runtime.NumCPU()
	parts := make([]Nodes, numParts)
	for i := 0; i < numParts; i++ {
		parts[i] = nodes[i*len(nodes)/numParts : (i+1)*len(nodes)/numParts]
	}
	return parts
}

// Exists returns a bool indicating if the given node is in
// this list of nodes.  BEWARE: at worst this is an O(n) operation.
func (nodes Nodes) Exists(node INode) bool {
	for _, n := range nodes {
		if n == nil {
			continue
		}
		if n.ID() == node.ID() {
			return true
		}
	}

	return false
}

// Dispose will free the resources of this list of nodes to be recycled
// to reduce the number of allocations.
func (nodes *Nodes) Dispose() {
	for i := 0; i < len(*nodes); i++ {
		(*nodes)[i] = nil
	}

	(*nodes) = (*nodes)[:0]
	nodesPool.Put(*nodes)
}

// Highest returns the highest ID seen in this list.
func (nodes Nodes) Highest() uint64 {
	var highest uint64
	for _, node := range nodes {
		if node.ID() > highest {
			highest = node.ID()
		}
	}

	return highest
}

// NewNodes is the constructor to create a typed list of nodes.  This will
// pull from a pool to reduce allocations.
func NewNodes(nodes ...INode) Nodes {
	r := nodesPool.Get()

	var n Nodes
	switch r.(type) {
	case []INode:
		n = Nodes(r.([]INode))
	default:
		n = r.(Nodes)
	}

	n = append(n, nodes...)
	return n
}
