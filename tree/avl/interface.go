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

package avl

// Entries is a list of type Entry.
type Entries []Entry

// Entry represents all items that can be placed into the AVL tree.
// They must implement a Compare method that can be used to determine
// the Entry's correct place in the tree.  Any object can implement
// Compare.
type Entry interface {
	// Compare should return a value indicating the relationship
	// of this Entry to the provided Entry.  A -1 means this entry
	// is less than, 0 means equality, and 1 means greater than.
	Compare(Entry) int
}
