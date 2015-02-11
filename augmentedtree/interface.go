/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Package augmentedtree is designed to be useful when checking
for intersection of ranges in n-dimensions.  For instance, if you imagine
an xy plane then the augmented tree is for telling you if
plane defined by the points (0, 0) and (10, 10).  The augmented
tree can tell you if that plane overlaps with a plane defined by
(-5, -5) and (5, 5) (true in this case).  You can also check
intersections against a point by constructing a range of encompassed
solely if a single point.

The current tree is a simple top-down red-black binary search tree.

TODO: Add a bottom-up implementation to assist with duplicate
range handling.
*/
package augmentedtree

// Interval is the interface that must be implemented by any
// item added to the interval tree.
type Interval interface {
	// LowAtDimension returns an integer representing the lower bound
	// at the requested dimension.
	LowAtDimension(uint64) int64
	// HighAtDimension returns an integer representing the higher bound
	// at the requested dimension.
	HighAtDimension(uint64) int64
	// OverlapsAtDimension should return a bool indicating if the provided
	// interval overlaps this interval at the dimension requested.
	OverlapsAtDimension(Interval, uint64) bool
	// ID should be a unique ID representing this interval.  This
	// is used to identify which interval to delete from the tree if
	// there are duplicates.
	ID() uint64
}

// Tree defines the object that is returned from the
// tree constructor.  We use a Tree interface here because
// the returned tree could be a single dimension or many
// dimensions.
type Tree interface {
	// Add will add the provided intervals to the tree.
	Add(intervals ...Interval)
	// Len returns the number of intervals in the tree.
	Len() uint64
	// Delete will remove the provided intervals from the tree.
	Delete(intervals ...Interval)
	// Query will return a list of intervals that intersect the provided
	// interval.  The provided interval's ID method is ignored so the
	// provided ID is irrelevant.
	Query(interval Interval) Intervals
	// Insert will shift intervals in the tree based on the specified
	// index and the specified count.  Dimension specifies where to
	// apply the shift.  Returned is a list of intervals impacted and
	// list of intervals deleted.  Intervals are deleted if the shift
	// makes the interval size zero or less, ie, min >= max.  These
	// intervals are automatically removed from the tree.  The tree
	// does not alter the ranges on the intervals themselves, the consumer
	// is expected to do that.
	Insert(dimension uint64, index, count int64) (Intervals, Intervals)
}
