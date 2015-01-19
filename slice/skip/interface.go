package skip

// Entry defines items that can be inserted into the skip list.
// This will also be the type returned from a query.
type Entry interface {
	// Key defines this entry's place in the skip list.
	Key() uint64
}

// Entries is a typed list of interface Entry.
type Entries []Entry
