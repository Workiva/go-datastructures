package augmentedtree

import (
	"sync"

	"github.com/Workiva/gotable/structures/bitarray"
)

var bitarrayPool = sync.Pool{
	New: func() interface{} {
		return bitarray.NewSparseBitArray()
	},
}

// multidimensionalTree implements an interval tree
// in multiple dimensions.
type multiDimensionalTree struct {
	dimensions trees
	chunks     []trees
}

// Insert will insert the provided intervals into the tree.
func (mdt *multiDimensionalTree) Insert(intervals ...Interval) {
	var wg sync.WaitGroup
	wg.Add(len(mdt.chunks))

	for _, chunk := range mdt.chunks {
		go func(trees trees) {
			for _, tree := range trees {
				tree.Insert(intervals...)
			}

			wg.Done()
		}(chunk)
	}

	wg.Wait()
}

// Delete removes the provided intervals from the tree.
func (mdt *multiDimensionalTree) Delete(intervals ...Interval) {
	var wg sync.WaitGroup
	wg.Add(len(mdt.chunks))

	for _, chunk := range mdt.chunks {
		go func(trees trees) {
			for _, tree := range trees {
				tree.Delete(intervals...)
			}

			wg.Done()
		}(chunk)
	}

	wg.Wait()
}

// Len returns the number of items in the tree.
func (mdt *multiDimensionalTree) Len() uint64 {
	return mdt.dimensions[0].Len()
}

// Max returns the rightmost value at the provided dimension.
func (mdt *multiDimensionalTree) Max(dimension uint64) int64 {
	if dimension > uint64(len(mdt.dimensions)) {
		return 0
	}

	return mdt.dimensions[dimension-1].Max(1)
}

// Min returns the leftmost value at the provided dimension.
func (mdt *multiDimensionalTree) Min(dimension uint64) int64 {
	if dimension > uint64(len(mdt.dimensions)) {
		return 0
	}

	return mdt.dimensions[dimension-1].Min(1)
}

// Query will return a list of intervals that intersect the provided
// interval.  The provided interval's ID method is ignored.
func (mdt *multiDimensionalTree) Query(interval Interval) Intervals {
	bas := make([]bitarray.BitArray, len(mdt.dimensions)-1)
	var wg sync.WaitGroup
	var intervals Intervals
	wg.Add(len(mdt.chunks))

	for _, chunk := range mdt.chunks {
		go func(trees trees) {
			for _, tree := range trees {
				if tree.dimension == uint64(len(mdt.dimensions)) { // i am the last dimension
					intervals = tree.Query(interval)
					continue
				}

				ba := bitarrayPool.Get().(bitarray.BitArray)
				tree.apply(interval, func(n *node) {
					ba.SetBit(n.id)
				})

				bas[tree.dimension-1] = ba
			}

			wg.Done()
		}(chunk)
	}

	wg.Wait()

	result := intervalsPool.Get().(Intervals)
	for _, interval := range intervals {
		intersects := true
		for _, ba := range bas {
			if ok, _ := ba.GetBit(interval.ID()); !ok {
				intersects = false
				break
			}
		}

		if intersects {
			result = append(result, interval)
		}
	}

	go intervals.Dispose()
	go func() {
		for _, ba := range bas {
			ba.Reset()
			bitarrayPool.Put(ba)
		}
	}()

	return result
}

func newMultiDimensionalTree(dimensions uint64) *multiDimensionalTree {
	ts := make(trees, 0, dimensions)
	for i := uint64(0); i < dimensions; i++ {
		ts = append(ts, newTree(i+1))
	}

	split := ts.split()
	// we're going to remove the empty lists in the case with
	// a low number of dimensions
	chunks := make([]trees, 0, len(split))
	for _, chunk := range split {
		if len(chunk) == 0 {
			continue
		}

		chunks = append(chunks, chunk)
	}

	return &multiDimensionalTree{dimensions: ts, chunks: chunks}
}
