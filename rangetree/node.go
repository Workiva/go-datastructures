package rangetree

type nodes []*node

type node struct {
	value        int64
	entry        Entry
	orderedNodes orderedNodes
}

func newNode(value int64, entry Entry, needNextDimension bool) *node {
	n := &node{}
	n.value = value
	if needNextDimension {
		n.orderedNodes = make(orderedNodes, 0, 10)
	} else {
		n.entry = entry
	}

	return n
}
