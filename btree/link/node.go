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

func search(parent *node, key Key) Key {
	parent = getParent(parent, nil, key)
	parent.lock.RLock()
	println(`FOUND`)
	log.Printf(`BEFORE MOVE: %+v`, parent)
	parent = moveRight(parent, key, false)
	defer parent.lock.RUnlock()

	log.Printf(`PARENT: %+v`, parent)
	i := parent.search(key)
	if i == len(parent.keys) {
		return nil
	}

	return parent.keys[i]
}

func getParent(parent *node, stack *nodes, key Key) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		parent.lock.RLock()
		log.Printf(`PARENT: %+v, %p`, parent, parent)
		parent = moveRight(parent, key, false) // if this happens on the root this should always just return
		n = parent.searchNode(key)

		if stack != nil {
			stack.push(parent)
		}

		parent.lock.RUnlock()
		parent = n
	}

	return parent
}

func insert(tree *blink, parent *node, stack *nodes, key Key) Key {
	parent = getParent(parent, stack, key)

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

func split(tree *blink, n *node, stack *nodes) {
	var l, r *node
	var k Key
	var parent *node
	for n.needsSplit() {
		k, l, r = n.split()
		parent = stack.pop()
		if parent == nil {
			parent = newNode(false, make(Keys, 0, tree.ary), make(nodes, 0, tree.ary+1))
			parent.maxSeen = r.max()
			parent.keys.insert(k)
			parent.nodes.push(l)
			parent.nodes.push(r)
			tree.lock.Lock()
			tree.root = parent
			tree.lock.Unlock()
			n.lock.Unlock()
			return
		}

		parent.lock.Lock()
		parent = moveRight(parent, r.key(), true)
		i := parent.search(k)
		parent.keys.insertAt(k, i)
		parent.nodes[i] = l
		parent.nodes.insertAt(r, i+1)

		n.lock.Unlock()
		n = parent
	}

	n.lock.Unlock()
}

func moveRight(node *node, key Key, getLock bool) *node {
	for {
		if len(node.keys) == 0 || node.right == nil { // this is either the node or the rightmost node
			return node
		}
		if key.Compare(node.max()) < 1 {
			return node
		}

		if getLock {
			node.right.lock.Lock()
			node.lock.Unlock()
		} else {
			node.right.lock.RLock()
			node.lock.RUnlock()
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

func (ns *nodes) insertAt(n *node, i int) {
	if i == len(*ns) {
		*ns = append(*ns, n)
		return
	}

	*ns = append(*ns, nil)
	copy((*ns)[i+1:], (*ns)[i:])
	(*ns)[i] = n
}

func (ns *nodes) splitAt(i int) (nodes, nodes) {
	length := len(*ns) - i
	right := make(nodes, length, cap(*ns))
	copy(right, (*ns)[i+1:])
	for j := i + 1; j < len(*ns); j++ {
		(*ns)[j] = nil
	}
	*ns = (*ns)[:i+1]
	return *ns, right
}

type node struct {
	keys    Keys
	nodes   nodes
	right   *node
	lock    sync.RWMutex
	isLeaf  bool
	maxSeen Key
}

func (n *node) key() Key {
	return n.keys.first()
}

func (n *node) insert(key Key) Key {
	if !n.isLeaf {
		panic(`Can't only insert key in an internal node.`)
	}

	overwritten := n.keys.insert(key)
	return overwritten
}

func (n *node) insertNode(other *node) {
	key := other.key()
	i := n.keys.search(key)
	n.keys.insertAt(key, i)
	n.nodes.insertAt(other, i)
}

func (n *node) needsSplit() bool {
	return n.keys.needsSplit()
}

func (n *node) max() Key {
	if n.isLeaf {
		return n.keys.last()
	}

	return n.maxSeen
}

func (n *node) splitLeaf() (Key, *node, *node) {
	i := (len(n.keys) / 2)
	key := n.keys[i]
	_, rightKeys := n.keys.splitAt(i - 1)
	nn := &node{
		keys:   rightKeys,
		right:  n.right,
		isLeaf: true,
	}
	n.right = nn
	return key, n, nn
}

func (n *node) splitInternal() (Key, *node, *node) {
	i := (len(n.keys) / 2)
	key := n.keys[i]

	rightKeys := make(Keys, len(n.keys)-1-i, cap(n.keys))
	rightNodes := make(nodes, len(rightKeys)+1, cap(n.nodes))

	copy(rightKeys, n.keys[i+1:])
	copy(rightNodes, n.nodes[i+1:])

	// for garbage collection
	for j := i + 1; j < len(n.nodes); j++ {
		if j != len(n.keys) {
			n.keys[j] = nil
		}
		n.nodes[j] = nil
	}

	nn := newNode(false, rightKeys, rightNodes)
	nn.maxSeen = n.max()

	n.maxSeen = key
	n.keys = n.keys[:i]
	n.nodes = n.nodes[:i+1]
	n.right = nn

	return key, n, nn
}

func (n *node) split() (Key, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf()
	}

	return n.splitInternal()
}

func (n *node) search(key Key) int {
	return n.keys.search(key)
}

func (n *node) searchNode(key Key) *node {
	i := n.search(key)

	if i < len(n.keys) && n.keys[i].Compare(key) == 0 {
		i++
	}

	return n.nodes[i]
}

func (n *node) print(output *log.Logger) {
	output.Printf(`NODE: %+v, %p`, n, n)
	if !n.isLeaf {
		for _, n := range n.nodes {
			if n == nil {
				output.Println(`NIL NODE`)
				continue
			}
			n.print(output)
		}
	}
}

func newNode(isLeaf bool, keys Keys, ns nodes) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
		nodes:  ns,
	}
}
