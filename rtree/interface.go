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

package rtree

// Rectangles is a typed list of Rectangle.
type Rectangles []Rectangle

// Rectangle describes a two-dimensional bound.
type Rectangle interface {
	// LowerLeft describes the lower left coordinate of this rectangle.
	LowerLeft() (int32, int32)
	// UpperRight describes the upper right coordinate of this rectangle.
	UpperRight() (int32, int32)
}

// RTree defines an object that can be returned from any subpackage
// of this package.
type RTree interface {
	// Search will perform an intersection search of the given
	// rectangle and return any rectangles that intersect.
	Search(Rectangle) Rectangles
	// Len returns in the number of items in the RTree.
	Len() uint64
	// Dispose will clean up any objects used by the RTree.
	Dispose()
	// Delete will remove the provided rectangles from the RTree.
	Delete(...Rectangle)
	// Insert will add the provided rectangles to the RTree.
	Insert(...Rectangle)
}
