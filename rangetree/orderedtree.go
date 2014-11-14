package rangetree

func isLastDimension(value, test uint64) bool {
	return test >= value
}

type nodeBundle struct {
	list  *orderedNodes
	index int
}

type orderedTree struct {
	top        orderedNodes
	number     uint64
	dimensions uint64
	path       []*nodeBundle
}

func (ot *orderedTree) resetPath() {
	ot.path = ot.path[:0]
}

func (ot *orderedTree) needNextDimension() bool {
	return ot.dimensions > 1
}

func (ot *orderedTree) add(entry Entry) {
	var node *node
	list := &ot.top

	for i := uint64(1); i <= ot.dimensions; i++ {
		if isLastDimension(ot.dimensions, i) {
			overwritten := list.add(
				newNode(entry.ValueAtDimension(i), entry, false),
			)
			if !overwritten {
				ot.number++
			}
			break
		}
		node, _ = list.getOrAdd(entry, i, ot.dimensions)
		list = &node.orderedNodes
	}
}

// Add will add the provided entries to the tree.
func (ot *orderedTree) Add(entries ...Entry) {
	for _, entry := range entries {
		if entry == nil {
			continue
		}

		ot.add(entry)
	}
}

func (ot *orderedTree) delete(entry Entry) {
	ot.resetPath()
	var index int
	var node *node
	list := &ot.top

	for i := uint64(1); i <= ot.dimensions; i++ {
		value := entry.ValueAtDimension(i)
		node, index = list.get(value)
		if node == nil { // there's nothing to delete
			return
		}

		nb := &nodeBundle{list: list, index: index}
		ot.path = append(ot.path, nb)

		list = &node.orderedNodes
	}

	ot.number--

	for i := len(ot.path) - 1; i >= 0; i-- {
		nb := ot.path[i]
		nb.list.deleteAt(nb.index)
		if len(*nb.list) > 0 {
			break
		}
	}
}

// Delete will remove the entries from the tree.
func (ot *orderedTree) Delete(entries ...Entry) {
	for _, entry := range entries {
		ot.delete(entry)
	}
}

// Len returns the number of items in the tree.
func (ot *orderedTree) Len() uint64 {
	return ot.number
}

func (ot *orderedTree) apply(list orderedNodes, interval Interval,
	dimension uint64, fn func(*node) bool) bool {

	low, high := interval.LowAtDimension(dimension), interval.HighAtDimension(dimension)

	if isLastDimension(ot.dimensions, dimension) {
		if !list.apply(low, high, fn) {
			return false
		}
	} else {
		if !list.apply(low, high, func(n *node) bool {
			if !ot.apply(n.orderedNodes, interval, dimension+1, fn) {
				return false
			}
			return true
		}) {
			return false
		}
		return true
	}

	return true
}

// Apply will call (in order) the provided function to every
// entry that falls within the provided interval.  Any alteration
// the the entry that would result in different answers to the
// interface methods results in undefined behavior.
func (ot *orderedTree) Apply(interval Interval, fn func(Entry) bool) {
	ot.apply(ot.top, interval, 1, func(n *node) bool {
		return fn(n.entry)
	})
}

// Query will return an ordered list of results in the given
// interval.
func (ot *orderedTree) Query(interval Interval) Entries {
	entries := NewEntries()

	ot.apply(ot.top, interval, 1, func(n *node) bool {
		entries = append(entries, n.entry)
		return true
	})

	return entries
}

// InsertAtDimension will increment items at and above the given index
// by the number provided.  Provide a negative number to to decrement.
// Returned are two lists.  The first list is a list of entries that
// were moved.  The second is a list entries that were deleted.  These
// lists are exclusive.
func (ot *orderedTree) InsertAtDimension(dimension uint64,
	index, number int64) (Entries, Entries) {

	// TODO: perhaps return an error here?
	if dimension > ot.dimensions {
		return nil, nil
	}

	modified := make(Entries, 0, 100)
	deleted := make(Entries, 0, 100)

	ot.top.insert(dimension, 1, ot.dimensions,
		index, number, &modified, &deleted,
	)

	return modified, deleted
}

func newOrderedTree(dimensions uint64) *orderedTree {
	return &orderedTree{
		dimensions: dimensions,
		path:       make([]*nodeBundle, 0, dimensions),
	}
}

// New is the constructor to create a new rangetree with
// the provided number of dimensions.
func New(dimensions uint64) RangeTree {
	return newOrderedTree(dimensions)
}
