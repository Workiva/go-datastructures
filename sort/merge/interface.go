package merge

// Comparators defines a typed list of type Comparator.
type Comparators []Comparator

// Less returns a bool indicating if the comparator at index i
// is less than the comparator at index j.
func (c Comparators) Less(i, j int) bool {
	return c[i].Compare(c[j]) < 0
}

// Len returns an int indicating the length of this list
// of comparators.
func (c Comparators) Len() int {
	return len(c)
}

// Swap swaps the values at positions i and j.
func (c Comparators) Swap(i, j int) {
	c[j], c[i] = c[i], c[j]
}

// Comparator defines items that can be sorted.  It contains
// a single method allowing the compare logic to compare one
// comparator to another.
type Comparator interface {
	// Compare will return a value indicating how this comparator
	// compares with the provided comparator.  A negative number
	// indicates this comparator is less than the provided comparator,
	// a 0 indicates equality, and a positive number indicates this
	// comparator is greater than the provided comparator.
	Compare(Comparator) int
}
