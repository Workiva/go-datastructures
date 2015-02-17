package hilbert

import (
	"sort"

	"github.com/Workiva/go-datastructures/rtree"
)

func getCenter(rect rtree.Rectangle) (int32, int32) {
	xlow, ylow := rect.LowerLeft()
	xhigh, yhigh := rect.UpperRight()

	return (xhigh + xlow) / 2, (yhigh + ylow) / 2
}

func getParent(root *node, key *hilbertBundle, cache []*node) *node {
	cache = append(cache, root)
	for !root.isLeaf {
		i := searchBundles(root.children, key)
		if i == len(root.children) {
			i--
		}
		root = root.children[i].(*node)
		cache = append(cache, root)
	}

	return root
}

func searchBundles(bundles []withHilbert, key withHilbert) int {
	return sort.Search(len(bundles), func(i int) bool {
		return bundles[i].hilbert() >= key.hilbert()
	})
}

func intersect(rect1, rect2 rtree.Rectangle) bool {
	xlow1, ylow1 := rect1.LowerLeft()
	xhigh2, yhigh2 := rect2.UpperRight()

	xhigh1, yhigh1 := rect1.UpperRight()
	xlow2, ylow2 := rect2.LowerLeft()

	return xhigh2 > xlow1 && xlow2 < xhigh1 && yhigh2 > ylow1 && ylow2 < yhigh1
}

type tree struct {
	root        *node
	number, ary uint64
	cache       []*node
}

func (t *tree) resetCache() {
	for i := range t.cache {
		t.cache[i] = nil
	}
	t.cache = t.cache[:0]
}

func (t *tree) Len() uint64 {
	return t.number
}

func (t *tree) insert(rect rtree.Rectangle) {
	hb := newHilbertBundle(rect)
	t.number++
	if t.root == nil {
		n := newNode(t.ary)
		n.isLeaf = true
		n.insert(hb)
		n.mbr = newRectangle(n.children)
		t.root = n
		return
	}

	t.resetCache()
	n := getParent(t.root, hb, t.cache)
	n.insert(hb)

	var key withHilbert
	var left, right *node
	for n != nil {
		if key != nil {
			n.insert(key)
		}
		if n.needsSplit(t.ary) {
			key, left, right = n.split(uint64(len(n.children) / 2))
		} else {
			key = nil
			n.adjust(hb)
		}
		n = n.parent
	}

	if key != nil {
		n := newNode(t.ary)
		n.insert(left)
		n.insert(right)
		t.root = n
	}
}

func (t *tree) Insert(rects ...rtree.Rectangle) {
	for _, r := range rects {
		t.insert(r)
	}
}

func (t *tree) Search(r rtree.Rectangle) []rtree.Rectangle {
	if t.root == nil {
		return []rtree.Rectangle{}
	}

	result := make([]rtree.Rectangle, 0, 10)
	whs := t.root.search(r)
	for len(whs) > 0 {
		wh := whs[0]
		if n, ok := wh.(*node); ok {
			whs = append(whs, n.search(r)...)
		} else {
			result = append(result, wh.(*hilbertBundle).Rectangle)
		}
		whs = whs[1:]
	}

	return result
}

func newTree(ary uint64) *tree {
	return &tree{ary: ary}
}
