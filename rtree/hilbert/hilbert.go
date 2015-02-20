package hilbert

import (
	"runtime"
	"sync"

	h "github.com/Workiva/go-datastructures/numerics/hilbert"
	"github.com/Workiva/go-datastructures/rtree"
)

func getCenter(rect rtree.Rectangle) (int32, int32) {
	xlow, ylow := rect.LowerLeft()
	xhigh, yhigh := rect.UpperRight()

	return (xhigh + xlow) / 2, (yhigh + ylow) / 2
}

type hilbertBundle struct {
	hilbert hilbert
	rect    rtree.Rectangle
}

func bundlesFromRects(rects ...rtree.Rectangle) []*hilbertBundle {
	chunks := chunkRectangles(rects, int64(runtime.NumCPU()))
	bundleChunks := make([][]*hilbertBundle, len(chunks))
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for i := 0; i < runtime.NumCPU(); i++ {
		if len(chunks[i]) == 0 {
			bundleChunks[i] = []*hilbertBundle{}
			wg.Done()
			continue
		}
		go func(i int) {
			bundles := make([]*hilbertBundle, 0, len(chunks[i]))
			for _, r := range chunks[i] {
				h := h.Encode(getCenter(r))
				bundles = append(bundles, &hilbertBundle{hilbert(h), r})
			}
			bundleChunks[i] = bundles
			wg.Done()
		}(i)
	}

	wg.Wait()

	bundles := make([]*hilbertBundle, 0, len(rects))
	for _, bc := range bundleChunks {
		bundles = append(bundles, bc...)
	}

	return bundles
}

// chunkRectangles takes a slice of rtree.Rectangle values and chunks it into `numParts` subslices.
func chunkRectangles(slice rtree.Rectangles, numParts int64) []rtree.Rectangles {
	parts := make([]rtree.Rectangles, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = slice[i*int64(len(slice))/numParts : (i+1)*int64(len(slice))/numParts]
	}
	return parts
}
