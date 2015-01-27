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

type nodes []*node

func (ns nodes) reset() {
	for i := range ns {
		ns[i] = nil
	}
}

type node struct {
	balance  int8 // bounded, |balance| should be <= 1
	children [2]*node
	entry    Entry
}

// copy returns a copy of this node with pointers to the original
// children.
func (n *node) copy() *node {
	return &node{
		balance:  n.balance,
		children: [2]*node{n.children[0], n.children[1]},
		entry:    n.entry,
	}
}

// newNode returns a new node for the provided entry.  A nil
// entry is used to represent the dummy node.
func newNode(entry Entry) *node {
	return &node{
		entry:    entry,
		children: [2]*node{},
	}
}
