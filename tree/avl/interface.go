package avl

type Entries []Entry

type Entry interface {
	// Compare should return a value indicating the relationship
	// of this Entry to the provided Entry.  A -1 means this entry
	// is less than, 0 means equality, and 1 means greater than.
	Compare(Entry) int
}
