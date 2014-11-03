package rangetree

import "fmt"

// NoEntriesError is returned from an operation that requires
// existing entries when none are found.
type NoEntriesError struct{}

func (nee NoEntriesError) Error() string {
	return `No entries in this tree.`
}

// OutOfDimensionError is returned when a requested operation
// doesn't meet dimensional requirements.
type OutOfDimensionError struct {
	provided, max uint64
}

func (oode OutOfDimensionError) Error() string {
	return fmt.Sprintf(`Provided dimension: %d is 
		greater than max dimension: %d`,
		oode.provided, oode.max,
	)
}
