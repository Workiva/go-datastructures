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

type widths []uint64

type nodes []*node

type node struct {
	// forward denotes the forward pointing pointers in this
	// node.
	forward nodes
	// widths keeps track of the distance between this pointer
	// and the forward pointers so we can access skip list
	// values by position in logarithmic time.
	widths widths
	// entry is the associated value with this node.
	entry Entry
}

func (n *node) Compare(e Entry) int {
	return n.entry.Compare(e)
}

// newNode will allocate and return a new node with the entry
// provided.  maxLevels will determine the length of the forward
// pointer list associated with this node.
func newNode(entry Entry, maxLevels uint8) *node {
	return &node{
		entry:   entry,
		forward: make(nodes, maxLevels),
		widths:  make(widths, maxLevels),
	}
}
