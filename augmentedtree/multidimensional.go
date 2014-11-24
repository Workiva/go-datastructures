package augmentedtree

/*
import (
	"sync"

	"github.com/Workiva/go-datastructures/bitarray"
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

// Add will insert the provided intervals into the tree.
func (mdt *multiDimensionalTree) Add(intervals ...Interval) {
	var wg sync.WaitGroup
	wg.Add(len(mdt.chunks))

	for _, chunk := range mdt.chunks {
		go func(trees trees) {
			for _, tree := range trees {
				tree.Add(intervals...)
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

// Insert will shift intervals in the tree based on the specified
// index and the specified count.  Dimension specifies where to
// apply the shift.  Returned is a list of intervals impacted and
// list of intervals deleted.  Intervals are deleted if the shift
// makes the interval size zero or less, ie, min >= max.  These
// intervals are automatically removed from the tree.  The tree
// does not alter the ranges on the intervals themselves, the consumer
// is expected to do that.
func (mdt *multiDimensionalTree) Insert(dimension uint64,
	index, count int64) (Intervals, Intervals) {

	dimension = dimension - 1

	if dimension >= uint64(len(mdt.dimensions)) { // invalid dimension
		return nil, nil
	}

	tree := mdt.dimensions[dimension]
	modified, deleted := tree.Insert(1, index, count)
	for i, tree := range mdt.dimensions {
		if uint64(i) == dimension {
			continue
		}

		tree.Delete(deleted...)
	}

	return modified, deleted
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
}*/

func newMultiDimensionalTree(dimensions uint64) *tree {
	return newTree(dimensions)
}
