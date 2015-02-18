package hilbert

import (
	"github.com/Workiva/go-datastructures/numerics/hilbert"
	"github.com/Workiva/go-datastructures/rtree"
)

type rectangle struct {
	xlow, xhigh, ylow, yhigh int32
}

func (r *rectangle) adjust(rect rtree.Rectangle) {
	x, y := rect.LowerLeft()
	if x < r.xlow {
		r.xlow = x
	}
	if y < r.ylow {
		r.ylow = y
	}

	x, y = rect.UpperRight()
	if x > r.xhigh {
		r.xhigh = x
	}

	if y > r.yhigh {
		r.yhigh = y
	}
}

func newRectangle(children []withHilbert) *rectangle {
	if len(children) == 0 {
		panic(`Cannot construct rectangle with no dimensions.`)
	}

	xlow, ylow := children[0].LowerLeft()
	xhigh, yhigh := children[0].UpperRight()
	r := &rectangle{
		xlow:  xlow,
		xhigh: xhigh,
		ylow:  ylow,
		yhigh: yhigh,
	}

	for i := 1; i < len(children); i++ {
		r.adjust(children[i])
	}

	return r
}

type withHilbert interface {
	rtree.Rectangle
	hilbert() int64
}

type hilbertBundle struct {
	rtree.Rectangle
	value int64
}

func (hb *hilbertBundle) hilbert() int64 {
	return hb.value
}

func newHilbertBundle(rect rtree.Rectangle) *hilbertBundle {
	x, y := getCenter(rect)
	h := hilbert.Encode(x, y)
	return &hilbertBundle{rect, h}
}

type node struct {
	right      *node
	children   []withHilbert
	isLeaf     bool
	maxHilbert int64
	mbr        *rectangle
	parent     *node
}

func (n *node) adjust(bundle *hilbertBundle) {
	if bundle.value > n.maxHilbert {
		n.maxHilbert = bundle.value
	}

	n.mbr.adjust(bundle)
}

func (n *node) hilbert() int64 {
	return n.maxHilbert
}

func (n *node) LowerLeft() (int32, int32) {
	return n.mbr.xlow, n.mbr.ylow
}

func (n *node) UpperRight() (int32, int32) {
	return n.mbr.xhigh, n.mbr.yhigh
}

func (n *node) search(rect rtree.Rectangle) []withHilbert {
	results := make([]withHilbert, 0)
	for _, child := range n.children {
		if intersect(child, rect) {
			results = append(results, child)
		}
	}

	return results
}

func (n *node) splitLeaf(i uint64) (withHilbert, *node, *node) {
	key := n.children[i]
	right := make([]withHilbert, uint64(len(n.children))-i-1, cap(n.children))
	copy(right, n.children[i+1:])
	for j := i + 1; j < uint64(len(n.children)); j++ {
		n.children[j] = nil // GC
	}
	n.children = n.children[:i+1]
	nn := &node{
		children:   right,
		isLeaf:     true,
		parent:     n.parent,
		mbr:        newRectangle(right),
		maxHilbert: right[len(right)-1].hilbert(),
	}
	n.right = nn
	n.mbr = newRectangle(n.children)
	n.maxHilbert = n.children[len(n.children)-1].hilbert()

	return key, n, nn
}

func (n *node) reset() {
	n.maxHilbert = n.children[len(n.children)-1].hilbert()
	n.mbr = newRectangle(n.children)
}

func (n *node) splitInternal(i uint64) (withHilbert, *node, *node) {
	key := n.children[i]
	right := make([]withHilbert, uint64(len(n.children))-i-1, cap(n.children))
	copy(right, n.children[i+1:])
	for j := i + 1; j < uint64(len(n.children)); j++ {
		n.children[j] = nil
	}
	n.children = n.children[:i+1]
	nn := &node{
		children: right,
		isLeaf:   false,
		parent:   n.parent,
	}
	n.reset()
	nn.reset()

	return key, n, nn
}

func (n *node) split(i uint64) (withHilbert, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf(i)
	}

	return n.splitInternal(i)
}

func (n *node) needsSplit(ary uint64) bool {
	return uint64(len(n.children)) >= ary
}

func (n *node) insert(hb withHilbert) {
	i := searchBundles(n.children, hb)
	if i == len(n.children) {
		n.children = append(n.children, hb)
		return
	}

	n.children = append(n.children, nil)
	copy(n.children[i+1:], n.children[i:])
	n.children[i] = hb
}

func newNode(ary uint64) *node {
	return &node{children: make([]withHilbert, 0, ary)}
}
