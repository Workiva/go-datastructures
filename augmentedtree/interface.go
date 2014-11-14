package augmentedtree

// Interval is the interface that must be implemented by any
// item added to the interval tree.
type Interval interface {
	// LowAtDimension returns an integer representing the lower bound
	// at the requested dimension.
	LowAtDimension(uint64) int64
	// HighAtDimension returns an integer representing the higher bound
	// at the requested dimension.
	HighAtDimension(uint64) int64
	// OverlapsAtDimension should return a bool indicating if the provided
	// interval overlaps this interval at the dimension requested.
	OverlapsAtDimension(Interval, uint64) bool
	// ID should be a unique ID representing this interval.  This
	// is used to identify which interval to delete from the tree if
	// there are duplicates.
	ID() uint64
}

// Tree defines the object that is returned from the
// tree constructor.  We use a Tree interface here because
// the returned tree could be a single dimension or many
// dimensions.
type Tree interface {
	// Add will add the provided intervals to the tree.
	Add(intervals ...Interval)
	// Len returns the number of intervals in the tree.
	Len() uint64
	// Max returns the rightmost bound in the tree at the provided dimension.
	Max(dimension uint64) int64
	// Min returns the leftmost bound in the tree at the provided dimension.
	Min(dimension uint64) int64
	// Delete will remove the provided intervals from the tree.
	Delete(intervals ...Interval)
	// Query will return a list of intervals that intersect the provided
	// interval.  The provided interval's ID method is ignored so the
	// provided ID is irrelevant.
	Query(interval Interval) Intervals
	// Insert will shift intervals in the tree based on the specified
	// index and the specified count.  Dimension specifies where to
	// apply the shift.  Returned is a list of intervals impacted and
	// list of intervals deleted.  Intervals are deleted if the shift
	// makes the interval size zero or less, ie, min >= max.  These
	// intervals are automatically removed from the tree.  The tree
	// does not alter the ranges on the intervals themselves, the consumer
	// is expected to do that.
	//Insert(dimension uint64, index, count int64) (Intervals, Intervals)
}
