package plus

// Keys is a typed list of Key interfaces.
type Keys []Key

type Key interface {
	// Compare should return an int indicating how this key relates
	// to the provided key.  -1 will indicate less than, 0 will indicate
	// equality, and 1 will indicate greater than.  Duplicate keys
	// are allowed, but duplicate IDs are not.
	Compare(Key) int
}

// Iterator will be called with matching keys until either false is
// returned or we run out of keys to iterate.
type Iterator interface {
	// Next will move the iterator to the next position and return
	// a bool indicating if there is a value.
	Next() bool
	// Value returns a Key at the associated iterator position.  Returns
	// nil if the iterator is exhausted or has never been nexted.
	Value() Key
	// exhaust is an internal helper method to iterate this iterator
	// until exhausted and returns the resulting list of keys.
	exhaust() keys
}
