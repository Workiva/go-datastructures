package plus

type Key interface {
	// Compare should return an int indicating how this key relates
	// to the provided key.  -1 will indicate less than, 0 will indicate
	// equality, and 1 will indicate greater than.  Duplicate keys
	// are allowed, but duplicate IDs are not.
	Compare(Key) int
	// ID should be a unique identifier for this key.  Keys that are
	// identical including IDs will not be inserted.
	ID() uint64
}

// Iterator will be called with matching keys until either false is
// returned or we run out of keys to iterate.
type Iterator interface {
	Next() bool
	Value() Key
}
