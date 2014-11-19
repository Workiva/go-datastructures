package slice

import "sort"

// Int64Slice is a slice that fulfills the sort.Interface interface.
type Int64Slice []int64

// Len returns the len of this slice.  Required by sort.Interface.
func (s Int64Slice) Len() int {
	return len(s)
}

// Less returns a bool indicating if the value at position i
// is less than at position j.  Required by sort.Interface.
func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

// Search will search this slice and return an index that corresponds
// to the lowest position of that value.  You'll need to check
// separately if the value at that position is equal to x.  The
// behavior of this method is undefinited if the slice is not sorted.
func (s Int64Slice) Search(x int64) int {
	return sort.Search(len(s), func(i int) bool {
		return s[i] >= x
	})
}

// Sort will in-place sort this list of int64s.
func (s Int64Slice) Sort() {
	sort.Sort(s)
}

// Swap will swap the elements at positions i and j.  This is required
// by sort.Interface.
func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Exists returns a bool indicating if the provided value exists
// in this list.  This has undefined behavior if the list is not
// sorted.
func (s Int64Slice) Exists(x int64) bool {
	i := s.Search(x)
	if i == len(s) {
		return false
	}

	return s[i] == x
}

// Insert will insert x into the sorted position in this list
// and return a list with the value added.  If this slice has not
// been sorted Insert's behavior is undefined.
func (s Int64Slice) Insert(x int64) Int64Slice {
	i := s.Search(x)
	if i == len(s) {
		return append(s, x)
	}

	if s[i] == x {
		return s
	}

	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = x
	return s
}
