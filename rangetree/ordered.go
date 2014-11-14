package rangetree

import "sort"

// orderedNodes represents an ordered list of points living
// at the last dimension.  No duplicates can be inserted here.
type orderedNodes nodes

func (nodes orderedNodes) search(value int64) int {
	return sort.Search(
		len(nodes),
		func(i int) bool { return nodes[i].value >= value },
	)
}

func (nodes *orderedNodes) addAt(i int, node *node) bool {
	if i == len(*nodes) {
		*nodes = append(*nodes, node)
		return false
	}

	if (*nodes)[i].value == node.value {
		// this is a duplicate, there can't be a duplicate
		// point in the last dimension
		(*nodes)[i] = node
		return true
	}

	*nodes = append(*nodes, nil)
	copy((*nodes)[i+1:], (*nodes)[i:])
	(*nodes)[i] = node
	return false
}

func (nodes *orderedNodes) add(node *node) bool {
	i := nodes.search(node.value)
	return nodes.addAt(i, node)
}

func (nodes *orderedNodes) deleteAt(i int) {
	if i == len(*nodes) { // no matching found
		return
	}

	copy((*nodes)[i:], (*nodes)[i+1:])
	(*nodes)[len(*nodes)-1] = nil
	*nodes = (*nodes)[:len(*nodes)-1]
}

func (nodes *orderedNodes) delete(value int64) {
	i := nodes.search(value)

	if (*nodes)[i].value != value || i == len(*nodes) {
		return
	}

	nodes.deleteAt(i)
}

func (nodes orderedNodes) apply(low, high int64, fn func(*node) bool) bool {
	index := nodes.search(low)
	if index == len(nodes) {
		return true
	}

	for ; index < len(nodes); index++ {
		if nodes[index].value >= high {
			break
		}

		if !fn(nodes[index]) {
			return false
		}
	}

	return true
}

func (nodes orderedNodes) get(value int64) (*node, int) {
	i := nodes.search(value)
	if i == len(nodes) {
		return nil, i
	}

	if nodes[i].value == value {
		return nodes[i], i
	}

	return nil, i
}

func (nodes *orderedNodes) getOrAdd(entry Entry,
	dimension, lastDimension uint64) (*node, bool) {

	isLastDimension := isLastDimension(lastDimension, dimension)
	value := entry.ValueAtDimension(dimension)

	i := nodes.search(value)
	if i == len(*nodes) {
		node := newNode(value, entry, !isLastDimension)
		*nodes = append(*nodes, node)
		return node, true
	}

	if (*nodes)[i].value == value {
		return (*nodes)[i], false
	}

	node := newNode(value, entry, !isLastDimension)
	*nodes = append(*nodes, nil)
	copy((*nodes)[i+1:], (*nodes)[i:])
	(*nodes)[i] = node
	return node, true
}
