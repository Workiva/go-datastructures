package rangetree

// Entry defines items that can be added to the rangetree.
type Entry interface {
	// ValueAtDimension returns the value of this entry
	// at the specified dimension.
	ValueAtDimension(dimension uint64) int64
}

// Interval describes the methods required to query the rangetree.
type Interval interface {
	// LowAtDimension returns an integer representing the lower bound
	// at the requested dimension.
	LowAtDimension(dimension uint64) int64
	// HighAtDimension returns an integer representing the higher bound
	// at the request dimension.
	HighAtDimension(dimension uint64) int64
}

// RangeTree describes the methods available to the rangetree.
type RangeTree interface {
	// Insert will add the provided entries to the tree.
	Insert(entries ...Entry)
	// Len returns the number of entries in the tree.
	Len() uint64
	// Delete will remove the provided entries from the tree.
	Delete(entries ...Entry)
	// Query will return a list of entries that fall within
	// the provided interval.
	Query(interval Interval) Entries
	// Apply will call the provided function with each entry that exists
	// within the provided range, in order.  Return false at any time to
	// cancel iteration.  Altering the entry in such a way that its location
	// changes will result in undefined behavior.
	Apply(interval Interval, fn func(Entry) bool)
}
