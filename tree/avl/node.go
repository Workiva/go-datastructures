package avl

type nodes []*node

func (ns nodes) reset() {
	for i := range ns {
		ns[i] = nil
	}
}

type node struct {
	balance  int8
	children [2]*node
	entry    Entry
}

func (n *node) copy() *node {
	return &node{
		balance:  n.balance,
		children: [2]*node{n.children[0], n.children[1]},
		entry:    n.entry,
	}
}

func newNode(entry Entry) *node {
	return &node{
		entry:    entry,
		children: [2]*node{},
	}
}
