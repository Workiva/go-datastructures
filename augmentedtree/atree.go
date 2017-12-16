/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package augmentedtree

func intervalOverlaps(n *node, low, high int64, interval Interval, maxDimension uint64) bool {
	if !overlaps(n.interval.HighAtDimension(1), high, n.interval.LowAtDimension(1), low) {
		return false
	}

	if interval == nil {
		return true
	}

	for i := uint64(2); i <= maxDimension; i++ {
		if !n.interval.OverlapsAtDimension(interval, i) {
			return false
		}
	}

	return true
}

func overlaps(high, otherHigh, low, otherLow int64) bool {
	return high >= otherLow && low <= otherHigh
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
	interval Interval
	max, min int64    // max value held by children
	children [2]*node // array to hold left/right
	red      bool     // indicates if this node is red
	id       uint64   // we store the id locally to reduce the number of calls to the method on the interface
}

func (n *node) query(low, high int64, interval Interval, maxDimension uint64, fn func(node *node)) {
	if n.children[0] != nil && overlaps(n.children[0].max, high, n.children[0].min, low) {
		n.children[0].query(low, high, interval, maxDimension, fn)
	}

	if intervalOverlaps(n, low, high, interval, maxDimension) {
		fn(n)
	}

	if n.children[1] != nil && overlaps(n.children[1].max, high, n.children[1].min, low) {
		n.children[1].query(low, high, interval, maxDimension, fn)
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
	}

	return itn
}

type tree struct {
	root                 *node
	maxDimension, number uint64
	dummy                node
}

func (t *tree) Traverse(fn func(id Interval)) {
	nodes := []*node{t.root}

	for len(nodes) != 0 {
		c := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		if c != nil {
			fn(c.interval)
			if c.children[0] != nil {
				nodes = append(nodes, c.children[0])
			}
			if c.children[1] != nil {
				nodes = append(nodes, c.children[1])
			}
		}
	}
}

func (tree *tree) resetDummy() {
	tree.dummy.children[0], tree.dummy.children[1] = nil, nil
	tree.dummy.red = false
}

// Len returns the number of items in this tree.
func (tree *tree) Len() uint64 {
	return tree.number
}

// add will add the provided interval to the tree.
func (tree *tree) add(iv Interval) {
	if tree.root == nil {
		tree.root = newNode(
			iv, iv.LowAtDimension(1),
			iv.HighAtDimension(1),
			1,
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
		max                 = iv.HighAtDimension(1)
		ivLow               = iv.LowAtDimension(1)
		helper              = &dummy
	)

	// set this AFTER clearing dummy
	helper.children[1] = tree.root
	for {
		if node == nil {
			node = newNode(iv, ivLow, max, 1)
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
		dir = compare(node.interval.LowAtDimension(1), ivLow, node.id, id)

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
		ivLow                      = iv.LowAtDimension(1)
	)

	node.children[1] = tree.root
	for node.children[dir] != nil {
		last = dir
		otherLast = takeOpposite(last)

		grandParent, parent, node = parent, node, node.children[dir]

		dir = compare(node.interval.LowAtDimension(1), ivLow, node.id, id)
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
		found.interval, found.max, found.min, found.id = node.interval, node.max, node.min, node.id
		parentDir := intFromBool(parent.children[1] == node)
		childDir := intFromBool(node.children[0] == nil)

		parent.children[parentDir] = node.children[childDir]
	}

	tree.root = dummy.children[1]
	if tree.root != nil {
		tree.root.red = false
	}
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
		ivLow     = interval.LowAtDimension(1)
		ivHigh    = interval.HighAtDimension(1)
	)

	tree.root.query(ivLow, ivHigh, interval, tree.maxDimension, func(node *node) {
		Intervals = append(Intervals, node.interval)
	})

	return Intervals
}

func isRed(node *node) bool {
	return node != nil && node.red
}

func setMax(parent *node) {
	parent.max = parent.interval.HighAtDimension(1)

	if parent.children[0] != nil && parent.children[0].max > parent.max {
		parent.max = parent.children[0].max
	}

	if parent.children[1] != nil && parent.children[1].max > parent.max {
		parent.max = parent.children[1].max
	}
}

func setMin(parent *node) {
	parent.min = parent.interval.LowAtDimension(1)
	if parent.children[0] != nil && parent.children[0].min < parent.min {
		parent.min = parent.children[0].min
	}

	if parent.children[1] != nil && parent.children[1].min < parent.min {
		parent.min = parent.children[1].min
	}

	if parent.interval.LowAtDimension(1) < parent.min {
		parent.min = parent.interval.LowAtDimension(1)
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

func newTree(maxDimension uint64) *tree {
	return &tree{
		maxDimension: maxDimension,
		dummy:        newDummy(),
	}
}

// New constructs and returns a new interval tree with the max
// dimensions provided.
func New(dimensions uint64) Tree {
	return newTree(dimensions)
}
