package hilbert

import (
	"log"
	"sort"

	"github.com/Workiva/go-datastructures/rtree"
)

func init() {
	log.Printf(`I HATE THIS`)
}

func getCenter(rect rtree.Rectangle) (int32, int32) {
	xlow, ylow := rect.LowerLeft()
	xhigh, yhigh := rect.UpperRight()

	return (xhigh + xlow) / 2, (yhigh + ylow) / 2
}

func getParent(root *node, key *hilbertBundle) *node {
	for !root.isLeaf {
		i := searchBundles(root.children, key)
		if i == len(root.children) {
			i--
		}
		root = root.children[i].(*node)
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

	return xhigh2 >= xlow1 && xlow2 <= xhigh1 && yhigh2 >= ylow1 && ylow2 <= yhigh1
}

func equals(rect1, rect2 rtree.Rectangle) bool {
	xlow1, ylow1 := rect1.LowerLeft()
	xhigh2, yhigh2 := rect2.UpperRight()

	xhigh1, yhigh1 := rect1.UpperRight()
	xlow2, ylow2 := rect2.LowerLeft()

	return xlow1 == xlow2 && xhigh1 == xhigh2 && ylow1 == ylow2 && yhigh1 == yhigh2
}

type tree struct {
	root        *node
	number, ary uint64
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

	n := getParent(t.root, hb)
	log.Printf(`ROOT: %+v, %p, PARENT: %+v, %p`, t.root, t.root, n, n)
	n.insert(hb)

	var key withHilbert
	var left, right *node
	for n != nil {
		if key != nil {
			n.insert(key)
		}
		if n.needsSplit(t.ary) {
			log.Printf(`N: %+v`, n)
			key, left, right = n.split(uint64(len(n.children) / 2))
			if n.isLeaf {
				key = right
			}
			log.Printf(`AFTER SPLIT KEY: %+v, LEFT: %+v, RIGHT: %+v`, key, left, right)
		} else {
			key = nil
			n.adjust(hb)
		}
		n = n.parent
	}

	if key != nil {
		println(`CREATING ROOT`)
		n := newNode(t.ary)
		n.insert(left)
		n.insert(right)
		t.root = n
		n.mbr = newRectangle(n.children)
		log.Printf(`LEFT: %+v, RIGHT: %+v`, left, right)
		n.maxHilbert = right.maxHilbert
		left.parent, right.parent = n, n
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

func (t *tree) Len() uint64 {
	return t.number
}

func newTree(ary uint64) *tree {
	return &tree{ary: ary}
}
