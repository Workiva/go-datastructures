package yfast

// Entry defines items that can be added to the rangetree.
type Entry interface {
	// ValueAtDimension returns the value of this entry
	// at the specified dimension.
	ValueAtDimension(dimension uint64) uint64
}

type Entries []Entry

// Interval describes the methods required to query the rangetree.
type Interval interface {
	// LowAtDimension returns an integer representing the lower bound
	// at the requested dimension.
	LowAtDimension(dimension uint64) uint64
	// HighAtDimension returns an integer representing the higher bound
	// at the request dimension.
	HighAtDimension(dimension uint64) uint64
}

type RangeTree interface {
	Add(entries ...Entry) Entries
}
