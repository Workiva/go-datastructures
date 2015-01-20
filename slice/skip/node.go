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

package skip

type nodes []*node

// reset will mark every index in this list of nodes as nil.
func (ns nodes) reset() nodes {
	for i := range ns {
		ns[i] = nil
	}

	return ns
}

// copy will create a deep copy of this list of nodes and return it.
func (ns nodes) copy() nodes {
	cp := make(nodes, len(ns))
	for i, n := range ns {
		if n == nil {
			break
		}

		cp[i] = n.copy()
	}

	return cp
}

type node struct {
	// forward denotes the forward pointing pointers in this
	// node.
	forward nodes
	// entry is the associated value with this node.
	entry Entry
}

// key is a helper method to return the key of the entry associated
// with this node.
func (n *node) key() uint64 {
	return n.entry.Key()
}

// copy will make a shallow copy of this node and return the result.
func (n *node) copy() *node {
	forward := make(nodes, len(n.forward))
	copy(forward, n.forward)
	return &node{
		forward: forward,
		entry:   n.entry,
	}
}

// newNode will allocate and return a new node with the entry
// provided.  maxLevels will determine the length of the forward
// pointer list associated with this node.
func newNode(entry Entry, maxLevels uint8) *node {
	return &node{
		entry:   entry,
		forward: make(nodes, maxLevels),
	}
}
