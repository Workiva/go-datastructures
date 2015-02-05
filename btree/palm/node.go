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

package palm

import "log"

func getParent(parent *node, key Key) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		n = parent.searchNode(key)
		parent = n
	}

	return parent
}

type nodes []*node

func (ns *nodes) push(n *node) {
	*ns = append(*ns, n)
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

type node struct {
	keys          Keys
	nodes         nodes
	isLeaf        bool
	parent, right *node
}

func (n *node) needsSplit(ary uint64) bool {
	return uint64(len(n.keys)) >= ary
}

func (n *node) splitLeaf() (Key, *node, *node) {
	i := (len(n.keys) / 2)
	key := n.keys[i]
	_, rightKeys := n.keys.splitAt(i)
	nn := &node{
		keys:   rightKeys,
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
	for _, nd := range rightNodes {
		nd.parent = nn
	}

	n.keys = n.keys[:i]
	n.nodes = n.nodes[:i+1]

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
