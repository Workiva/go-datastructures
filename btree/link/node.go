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

package link

import (
	"log"
	"sync"
)

func scan(node *node, key Key) (Key, int) {
	index := node.search(key)
	if index == len(node.keys) {
		right := moveRight(node, key, false)
		index = right.search(key)
		if index == len(right.keys) {
			index--
			return right.keys[index], index
		}

		return right.keys[index], index
	}

	return node.keys[index], index
}

func search(parent *node, key Key) (*node, int) {
	var found Key
	var ok bool
	for parent != nil && !parent.isLeaf {
		found, _ = scan(parent, key)
		parent, ok = found.(*node)
		if !ok {
			break
		}
	}

	parent = moveRight(parent, key, false)
	return parent, parent.search(key)
}

func insert(tree *blink, parent *node, stack nodes, key Key) Key {
	var found Key
	var index int
	var ok bool
	for parent != nil && !parent.isLeaf {
		found, index = scan(parent, key)
		if index < len(parent.keys) {
			stack.push(parent)
		}

		parent, ok = found.(*node)
		if !ok {
			break
		}
	}

	parent.lock.Lock()
	parent = moveRight(parent, key, true)

	result := parent.insert(key)
	if result != nil { // overwrite
		parent.lock.Unlock()
		return result
	}

	if !parent.needsSplit() {
		parent.lock.Unlock()
		return nil
	}

	split(tree, parent, stack)
	return nil
}

func split(tree *blink, n *node, stack nodes) {
	var l, r *node
	var parent *node
	for n.needsSplit() {
		l, r = n.split()
		parent = stack.pop()
		if parent == nil {
			parent = newNode(false)
			parent.key = r.key
			parent.insertNode(l)
			parent.insertNode(r)
			tree.lock.Lock()
			tree.root = parent
			tree.lock.Unlock()
			n.lock.Unlock()
			return
		}

		parent.lock.Lock()
		parent = moveRight(parent, r.key, true)
		parent.insertNode(r)
		n.lock.Unlock()
		n = parent
	}

	n.lock.Unlock()
}

func moveRight(node *node, key Key, getLock bool) *node {
	for {
		if node.key == nil || node.key.Compare(key) > -1 || node.right == nil { // this is either the node or the rightmost node
			return node
		}
		if getLock {
			node.right.lock.Lock()
			node.lock.Unlock()
		}
		node = node.right
	}
}

type nodes []*node

func (ns *nodes) push(n *node) {
	*ns = append(*ns, n)
}

func (ns *nodes) pop() *node {
	if len(*ns) == 0 {
		return nil
	}

	n := (*ns)[len(*ns)-1]
	(*ns)[len(*ns)-1] = nil
	*ns = (*ns)[:len(*ns)-1]
	return n
}

type node struct {
	keys   Keys
	key    Key
	right  *node
	lock   sync.RWMutex
	isLeaf bool
}

func (n *node) insertNode(other *node) {
	n.keys.insertNode(other)
}

func (n *node) insert(key Key) Key {
	result := n.keys.insert(key)
	if n.key == nil || key.Compare(n.key) > 0 {
		n.key = key
	}

	return result
}

func (n *node) needsSplit() bool {
	return n.keys.needsSplit()
}

func (n *node) split() (*node, *node) {
	key, _, right := n.keys.split()
	nn := &node{
		keys:   right,
		key:    right.last(),
		right:  n.right,
		isLeaf: n.isLeaf,
	}
	n.right = nn
	n.key = key
	return n, nn
}

func (n *node) search(key Key) int {
	return n.keys.search(key)
}

func (n *node) Compare(key Key) int {
	return n.key.Compare(key)
}

func (n *node) print(output *log.Logger) {
	output.Printf(`NODE: %+v, %p`, n, n)
	for _, k := range n.keys {
		result, ok := k.(*node)
		if ok {
			result.print(output)
		} else {
			output.Printf(`CHILD: %+v, %p`, k, k)
		}
	}
}

func newNode(isLeaf bool) *node {
	return &node{isLeaf: isLeaf}
}
