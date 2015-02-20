package hilbert

import (
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

func equal(r1, r2 rtree.Rectangle) bool {
	xlow1, ylow1 := r1.LowerLeft()
	xhigh2, yhigh2 := r2.UpperRight()

	xhigh1, yhigh1 := r1.UpperRight()
	xlow2, ylow2 := r2.LowerLeft()

	return xlow1 == xlow2 && xhigh1 == xhigh2 && ylow1 == ylow2 && yhigh1 == yhigh2
}

func intersect(rect1 *rectangle, rect2 rtree.Rectangle) bool {
	xhigh2, yhigh2 := rect2.UpperRight()
	xlow2, ylow2 := rect2.LowerLeft()

	return xhigh2 >= rect1.xlow && xlow2 <= rect1.xhigh && yhigh2 >= rect1.ylow && ylow2 <= rect1.yhigh
}

func newRectangeFromRect(rect rtree.Rectangle) *rectangle {
	r := &rectangle{}
	x, y := rect.LowerLeft()
	r.xlow = x
	r.ylow = y

	x, y = rect.UpperRight()
	r.xhigh = x
	r.yhigh = y

	return r
}

func newRectangleFromRects(rects rtree.Rectangles) *rectangle {
	if len(rects) == 0 {
		panic(`Cannot construct rectangle with no dimensions.`)
	}

	xlow, ylow := rects[0].LowerLeft()
	xhigh, yhigh := rects[0].UpperRight()
	r := &rectangle{
		xlow:  xlow,
		xhigh: xhigh,
		ylow:  ylow,
		yhigh: yhigh,
	}

	for i := 1; i < len(rects); i++ {
		r.adjust(rects[i])
	}

	return r
}
