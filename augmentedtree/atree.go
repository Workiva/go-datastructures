package augmentedtree

import "math"

func overlaps(high, otherHigh, low, otherLow int64) bool {
	return high > otherLow && low < otherHigh
}

// compare returns an int indicating which direction the node
// should go.
func compare(nodeLow, ivLow int64, nodeID, ivID uint64) int {
	if ivLow > nodeLow {
		return 1
	}

	if ivLow < nodeLow {
		return 0
	}

	return intFromBool(ivID > nodeID)
}

type node struct {
	interval            Interval
	low, high, max, min int64    // max value held by children
	children            [2]*node // array to hold left/right
	red                 bool     // indicates if this node is red
	id                  uint64   // we store the id locally to reduce the number of calls to the method on the interface
}

func (n *node) query(low, high int64, fn func(node *node)) {
	if n.children[0] != nil && overlaps(n.children[0].max, high, n.children[0].min, low) {
		n.children[0].query(low, high, fn)
	}

	if overlaps(n.high, high, n.low, low) {
		fn(n)
	}

	if n.children[1] != nil && overlaps(n.children[1].max, high, n.children[1].min, low) {
		n.children[1].query(low, high, fn)
	}
}

func (n *node) adjustRanges() {
	for i := 0; i <= 1; i++ {
		if n.children[i] != nil {
			n.children[i].adjustRanges()
		}
	}

	n.adjustRange()
}

func (n *node) adjustRange() {
	setMin(n)
	setMax(n)
}

func newDummy() node {
	return node{
		children: [2]*node{},
	}
}

func newNode(interval Interval, min, max int64, dimension uint64) *node {
	itn := &node{
		interval: interval,
		min:      min,
		max:      max,
		red:      true,
		children: [2]*node{},
	}
	if interval != nil {
		itn.id = interval.ID()
		itn.low = interval.LowAtDimension(dimension)
		itn.high = interval.HighAtDimension(dimension)
	}

	return itn
}

type tree struct {
	root              *node
	dimension, number uint64
	dummy             node
}

func (tree *tree) resetDummy() {
	tree.dummy.children[0], tree.dummy.children[1] = nil, nil
	tree.dummy.red = false
}

// Max is an extremely quick method for determining the rightmost
// bound in the tree.
func (tree *tree) Max(dimension uint64) int64 {
	if dimension != 1 {
		return 0
	}

	if tree.root == nil {
		return 0
	}

	return tree.root.max
}

// Min returns the minimum bound seen in the tree.
func (tree *tree) Min(dimension uint64) int64 {
	if dimension != 1 {
		return 0
	}

	if tree.root == nil {
		return 0
	}

	return tree.root.min
}

// Len returns the number of items in this tree.
func (tree *tree) Len() uint64 {
	return tree.number
}

// add will add the provided interval to the tree.
func (tree *tree) add(iv Interval) {
	if tree.root == nil {
		tree.root = newNode(
			iv, iv.LowAtDimension(tree.dimension),
			iv.HighAtDimension(tree.dimension),
			tree.dimension,
		)
		tree.root.red = false
		tree.number++
		return
	}

	tree.resetDummy()
	var (
		dummy               = tree.dummy
		parent, grandParent *node
		node                = tree.root
		dir, last           int
		otherLast           = 1
		id                  = iv.ID()
		max                 = iv.HighAtDimension(tree.dimension)
		ivLow               = iv.LowAtDimension(tree.dimension)
		helper              = &dummy
	)

	// set this AFTER clearing dummy
	helper.children[1] = tree.root
	for {
		if node == nil {
			node = newNode(iv, ivLow, max, tree.dimension)
			parent.children[dir] = node
			tree.number++
		} else if isRed(node.children[0]) && isRed(node.children[1]) {
			node.red = true
			node.children[0].red = false
			node.children[1].red = false
		}
		if max > node.max {
			node.max = max
		}

		if ivLow < node.min {
			node.min = ivLow
		}

		if isRed(parent) && isRed(node) {
			localDir := intFromBool(helper.children[1] == grandParent)

			if node == parent.children[last] {
				helper.children[localDir] = rotate(grandParent, otherLast)
			} else {
				helper.children[localDir] = doubleRotate(grandParent, otherLast)
			}
		}

		if node.id == id {
			break
		}

		last = dir
		otherLast = takeOpposite(last)
		dir = compare(node.low, ivLow, node.id, id)

		if grandParent != nil {
			helper = grandParent
		}
		grandParent, parent, node = parent, node, node.children[dir]
	}

	tree.root = dummy.children[1]
	tree.root.red = false
}

// Add will add the provided intervals to this tree.
func (tree *tree) Add(intervals ...Interval) {
	for _, iv := range intervals {
		tree.add(iv)
	}
}

// delete will remove the provided interval from the tree.
func (tree *tree) delete(iv Interval) {
	if tree.root == nil {
		return
	}

	tree.resetDummy()
	var (
		dummy                      = tree.dummy
		found, parent, grandParent *node
		last, otherDir, otherLast  int // keeping track of last direction
		id                         = iv.ID()
		dir                        = 1
		node                       = &dummy
		ivLow                      = iv.LowAtDimension(tree.dimension)
	)

	node.children[1] = tree.root
	for node.children[dir] != nil {
		last = dir
		otherLast = takeOpposite(last)

		grandParent, parent, node = parent, node, node.children[dir]

		dir = compare(node.low, ivLow, node.id, id)
		otherDir = takeOpposite(dir)

		if node.id == id {
			found = node
		}

		if !isRed(node) && !isRed(node.children[dir]) {
			if isRed(node.children[otherDir]) {
				parent.children[last] = rotate(node, dir)
				parent = parent.children[last]
			} else if !isRed(node.children[otherDir]) {
				t := parent.children[otherLast]

				if t != nil {
					if !isRed(t.children[otherLast]) && !isRed(t.children[last]) {
						parent.red = false
						node.red = true
						t.red = true
					} else {
						localDir := intFromBool(grandParent.children[1] == parent)

						if isRed(t.children[last]) {
							grandParent.children[localDir] = doubleRotate(
								parent, last,
							)
						} else if isRed(t.children[otherLast]) {
							grandParent.children[localDir] = rotate(
								parent, last,
							)
						}

						node.red = true
						grandParent.children[localDir].red = true
						grandParent.children[localDir].children[0].red = false
						grandParent.children[localDir].children[1].red = false
					}
				}
			}
		}
	}

	if found != nil {
		tree.number--
		found.interval, found.max, found.min, found.low, found.high, found.id = node.interval, node.max, node.min, node.low, node.high, node.id
		parentDir := intFromBool(parent.children[1] == node)
		childDir := intFromBool(node.children[0] == nil)

		parent.children[parentDir] = node.children[childDir]
	}

	tree.root = dummy.children[1]
	if tree.root != nil {
		tree.root.red = false
	}
}

// Insert will shift intervals in the tree based on the specified
// index and the specified count.  Dimension specifies where to
// apply the shift.  Returned is a list of intervals impacted and
// list of intervals deleted.  Intervals are deleted if the shift
// makes the interval size zero or less, ie, min >= max.  These
// intervals are automatically removed from the tree.  The tree
// does not alter the ranges on the intervals themselves, the consumer
// is expected to do that.
func (tree *tree) Insert(dimension uint64,
	index, count int64) (Intervals, Intervals) {

	if dimension != 1 { // consumer made a mistake
		return nil, nil
	}

	if tree.root == nil { // nothing to do
		return nil, nil
	}

	modified, deleted := intervalsPool.Get().(Intervals), intervalsPool.Get().(Intervals)

	tree.root.query(math.MinInt64, math.MaxInt64, func(n *node) {
		if n.max <= index { // won't change min or max in this case
			return
		}

		n.max += count
		if n.min >= index {
			n.min += count
		}

		mod := false
		if n.high > index {
			n.high += count
			if n.high < index {
				n.high = index
			}
			mod = true
		}
		if n.low > index {
			n.low += count
			if n.low < index {
				n.low = index
			}
			mod = true
		}

		if n.low >= n.high {
			deleted = append(deleted, n.interval)
		} else if mod {
			modified = append(modified, n.interval)
		}
	})

	tree.Delete(deleted...)

	return modified, deleted
}

// Delete will remove the provided intervals from this tree.
func (tree *tree) Delete(intervals ...Interval) {
	for _, iv := range intervals {
		tree.delete(iv)
	}
	if tree.root != nil {
		tree.root.adjustRanges()
	}
}

// Query will return a list of intervals that intersect the provided
// interval.  The provided interval's ID method is ignored so the
// provided ID is irrelevant.
func (tree *tree) Query(interval Interval) Intervals {
	if tree.root == nil {
		return nil
	}

	var (
		Intervals = intervalsPool.Get().(Intervals)
		ivLow     = interval.LowAtDimension(tree.dimension)
		ivHigh    = interval.HighAtDimension(tree.dimension)
	)

	tree.root.query(ivLow, ivHigh, func(node *node) {
		Intervals = append(Intervals, node.interval)
	})

	return Intervals
}

func (tree *tree) apply(interval Interval, fn func(*node)) {
	if tree.root == nil {
		return
	}

	low, high := interval.LowAtDimension(tree.dimension), interval.HighAtDimension(tree.dimension)
	tree.root.query(low, high, fn)
}

func isRed(node *node) bool {
	return node != nil && node.red
}

func setMax(parent *node) {
	parent.max = parent.high

	if parent.children[0] != nil && parent.children[0].max > parent.max {
		parent.max = parent.children[0].max
	}

	if parent.children[1] != nil && parent.children[1].max > parent.max {
		parent.max = parent.children[1].max
	}
}

func setMin(parent *node) {
	parent.min = parent.low
	if parent.children[0] != nil && parent.children[0].min < parent.min {
		parent.min = parent.children[0].min
	}

	if parent.children[1] != nil && parent.children[1].min < parent.min {
		parent.min = parent.children[1].min
	}

	if parent.low < parent.min {
		parent.min = parent.low
	}
}

func rotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	child := parent.children[otherDir]
	parent.children[otherDir] = child.children[dir]
	child.children[dir] = parent
	parent.red = true
	child.red = false
	child.max = parent.max
	setMax(child)
	setMax(parent)
	setMin(child)
	setMin(parent)

	return child
}

func doubleRotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	parent.children[otherDir] = rotate(parent.children[otherDir], otherDir)
	return rotate(parent, dir)
}

func intFromBool(value bool) int {
	if value {
		return 1
	}

	return 0
}

func takeOpposite(value int) int {
	return 1 - value
}

func newTree(dimension uint64) *tree {
	return &tree{
		dimension: dimension,
		dummy:     newDummy(),
	}
}

// New constructs and returns a new interval tree with the max
// dimensions provided.
func New(dimensions uint64) Tree {
	if dimensions == 1 {
		return newTree(1)
	}

	return newMultiDimensionalTree(dimensions)
}
