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
Package palm implements parallel architecture-friendly latch-free
modifications (PALM).  Details can be found here:

http://cs.unc.edu/~sewall/palm.pdf

The primary purpose of the tree is to efficiently batch operations
in such a way that locks are not required.  This is most beneficial
for in-memory indices.  Otherwise, the operations have typical B-tree
time complexities.

You primarily see the benefits of multithreading in availability and
bulk operations, below is a benchmark against the B-plus tree in this
package.

BenchmarkBulkAddToExisting-8	200	   8690207 ns/op
BenchmarkBulkAddToExisting-8    100   16778514 ns/op
*/

package palm

// Keys is a typed list of Key interfaces.
type Keys []Key

// Key defines items that can be inserted into or searched for
// in the tree.
type Key interface {
	// Compare should return an int indicating how this key relates
	// to the provided key.  -1 will indicate less than, 0 will indicate
	// equality, and 1 will indicate greater than.  Duplicate keys
	// are allowed, but duplicate IDs are not.
	Compare(Key) int
}

// BTree is the interface returned from this package's constructor.
type BTree interface {
	// Insert will insert the provided keys into the tree.
	Insert(...Key)
	// Get will return a key matching the associated provided
	// key if it exists.
	Get(...Key) Keys
	// Len returns the number of items in the tree.
	Len() uint64
	// Dispose will clean up any resources used by this tree.  This
	// must be called to prevent a memory leak.
	Dispose()
}
