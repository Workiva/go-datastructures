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

package btree

/*
This file contains logic pertaining to keeping track of the path followed
to find a particular node while descending the tree.
*/

type pathBundle struct {
	// i defines the child index of the n.
	i    int
	n    *Node
	prev *pathBundle
}

// path is simply a linked list of pathBundles.  We only ever
// go in one direction and there's no need to search so a linked list
// makes sense.
type path struct {
	head *pathBundle
	tail *pathBundle
}

func (p *path) append(pb *pathBundle) {
	if p.head == nil {
		p.head = pb
		p.tail = pb
		return
	}

	pb.prev = p.tail
	p.tail = pb
}

// pop removes the last item from the path.  Note that it also nils
// out the returned pathBundle's prev field.  Returns nil if no items
// remain.
func (p *path) pop() *pathBundle {
	if pb := p.tail; pb != nil {
		p.tail = pb.prev
		pb.prev = nil
		return pb
	}

	return nil
}

func (p *path) peek() *pathBundle {
	return p.tail
}
