/*
Copyright (c) 2016, Theodore Butler
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package dtrie provides an implementation of the dtrie data structure, which
// is a persistent hash trie that dynamically expands or shrinks to provide
// efficient memory allocation. This data structure is based on the papers
// Ideal Hash Trees by Phil Bagwell and Optimizing Hash-Array Mapped Tries for
// Fast and Lean Immutable JVM Collections by Michael J. Steindorfer and
// Jurgen J. Vinju
package dtrie

// Dtrie is a persistent hash trie that dynamically expands or shrinks
// to provide efficient memory allocation.
type Dtrie struct {
	root   *node
	hasher func(v interface{}) uint32
}

type entry struct {
	hash  uint32
	key   interface{}
	value interface{}
}

func (e *entry) KeyHash() uint32 {
	return e.hash
}

func (e *entry) Key() interface{} {
	return e.key
}

func (e *entry) Value() interface{} {
	return e.value
}

// New creates an empty DTrie with the given hashing function.
// If nil is passed in, the default hashing function will be used.
func New(hasher func(v interface{}) uint32) *Dtrie {
	if hasher == nil {
		hasher = defaultHasher
	}
	return &Dtrie{
		root:   emptyNode(0, 32),
		hasher: hasher,
	}
}

// Size returns the number of entries in the Dtrie.
func (d *Dtrie) Size() (size int) {
	for _ = range iterate(d.root, nil) {
		size++
	}
	return size
}

// Get returns the Entry for the associated key or returns nil if the
// key does not exist.
func (d *Dtrie) Get(key interface{}) Entry {
	return get(d.root, d.hasher(key), key)
}

// Insert adds an entry to the Dtrie, replacing the existing value if
// the key already exists and returns the resulting Dtrie.
func (d *Dtrie) Insert(key, value interface{}) *Dtrie {
	root := insert(d.root, &entry{d.hasher(key), key, value})
	return &Dtrie{root, d.hasher}
}

// Remove deletes the value for the associated key if it exists and returns
// the resulting Dtrie.
func (d *Dtrie) Remove(key interface{}) *Dtrie {
	root := remove(d.root, d.hasher(key), key)
	return &Dtrie{root, d.hasher}
}

// Iterator returns a read-only channel of Entries from the Dtrie. If a stop
// channel is provided, closing it will terminate and close the iterator
// channel. Note that if a cancel channel is not used and not every entry is
// read from the iterator, a goroutine will leak.
func (d *Dtrie) Iterator(stop <-chan struct{}) <-chan Entry {
	return iterate(d.root, stop)
}
