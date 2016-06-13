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
Package btree provides a very specific set implementation for k/v
lookup.  This is based on a B PALM tree as described here:
http://irvcvcs01.intel-research.net/publications/palm.pdf

This tree is best interacted with in batches.  Insertions and deletions
are optimized for dealing with large amounts of data.

Future work includes:

1) Optimization
2) Range scans

Usage:

rt := New(config)
mutable := rt.AsMutable()
... operations

rt, err := mutable.Commit() // saves all mutated nodes

.. rt reading/operations

Once a mutable has been committed, its further operations are undefined.
*/
package btree

// Tree describes the common functionality of both the read-only and mutable
// forms of a btree.
type Tree interface {
	// Apply takes a range and applies the provided function to every value
	// in that range in order.  If a key could not be found, it is
	// skipped.
	Apply(fn func(item *Item), keys ...interface{}) error
	// ID returns the identifier for this tree.
	ID() ID
	// Len returns the number of items in the tree.
	Len() int
}

// ReadableTree represents the operations that can be performed on a read-only
// version of the tree.  All reads of the readable tree are threadsafe and
// an indefinite number of mutable trees can be created from a single readable
// tree with the caveat that no mutable trees reflect any mutations to any other
// mutable tree.
type ReadableTree interface {
	Tree
	// AsMutable returns a mutable version of this tree.  The mutable version
	// has common mutations and you can create as many mutable versions of this
	// tree as you'd like.  However, the returned mutable is not threadsafe.
	AsMutable() MutableTree
}

// MutableTree represents a mutable version of the btree.  This interface
// is not threadsafe.
type MutableTree interface {
	Tree
	// Commit commits all mutated nodes to persistence and returns a
	// read-only version of this tree.  An error is returned if nodes
	// could not be committed to persistence.
	Commit() (ReadableTree, error)
	// AddItems adds the provided items to the btree.  Any existing items
	// are overwritten.  An error is returned if the tree could not be
	// traversed due to an error in the persistence layer.
	AddItems(items ...*Item) ([]*Item, error)
	// DeleteItems removes all provided keys and returns them.
	// An error is returned if the tree could not be traversed.
	DeleteItems(keys ...interface{}) ([]*Item, error)
}

// Comparator is used to determine ordering in the tree.  If item1
// is less than item2, a negative number should be returned and
// vice versa.  If equal, 0 should be returned.
type Comparator func(item1, item2 interface{}) int

// Payload is very basic and simply contains a key and a payload.
type Payload struct {
	Key     []byte
	Payload []byte
}

// Perister describes the interface of the different implementations.
// Given that we expect that datastrutures are immutable, we never
// have the need to delete.
type Persister interface {
	Save(items ...*Payload) error
	Load(keys ...[]byte) ([]*Payload, error)
}
