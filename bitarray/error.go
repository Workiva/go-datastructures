package bitarray

import "fmt"

// OutOfRangeError is an error caused by trying to access a bitarray past the end of its
// capacity.
type OutOfRangeError uint64

// Error returns a human readable description of the out-of-range error.
func (err OutOfRangeError) Error() string {
	return fmt.Sprintf(`Index %d is out of range.`, err)
}
