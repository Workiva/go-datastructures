package augmentedtree

import "sync"

var intervalsPool = sync.Pool{
	New: func() interface{} {
		return make(Intervals, 0, 10)
	},
}

// Intervals represents a list of Intervals.
type Intervals []Interval

// Dispose will free any consumed resources and allow this list to be
// re-allocated.
func (ivs *Intervals) Dispose() {
	for i := 0; i < len(*ivs); i++ {
		(*ivs)[i] = nil
	}

	*ivs = (*ivs)[:0]
	intervalsPool.Put(*ivs)
}
