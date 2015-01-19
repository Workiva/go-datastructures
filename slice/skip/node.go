package skip

type nodes []*node

func (ns nodes) reset() nodes {
	for i := range ns {
		ns[i] = nil
	}

	return ns
}

type node struct {
	// forward denotes the forward pointing pointers in this
	// node.
	forward nodes
	entry   Entry
}

func (n *node) key() uint64 {
	return n.entry.Key()
}

func newNode(entry Entry, maxLevels uint8) *node {
	return &node{
		entry:   entry,
		forward: make(nodes, maxLevels),
	}
}
