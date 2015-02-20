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

package hilbert

import (
	"log"
	"sort"

	"github.com/Workiva/go-datastructures/rtree"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

type hilbert int64

type hilberts []hilbert

func getParent(parent *node, key hilbert) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		n = parent.searchNode(key)
		parent = n
	}

	return parent
}

type nodes struct {
	list rtree.Rectangles
}

func (ns *nodes) push(n *node) {
	ns.list = append(ns.list, n)
}

func (ns *nodes) splitAt(i, capacity uint64) (*nodes, *nodes) {
	i++
	right := make(rtree.Rectangles, uint64(len(ns.list))-i, capacity)
	copy(right, ns.list[i:])
	for j := i; j < uint64(len(ns.list)); j++ {
		ns.list[j] = nil
	}
	ns.list = ns.list[:i]
	return ns, &nodes{list: right}
}

func (ns *nodes) byPosition(pos uint64) *node {
	if pos >= uint64(len(ns.list)) {
		return nil
	}

	return ns.list[pos].(*node)
}

func (ns *nodes) insertAt(i uint64, n rtree.Rectangle) {
	ns.list = append(ns.list, nil)
	copy(ns.list[i+1:], ns.list[i:])
	ns.list[i] = n
}

func (ns *nodes) replaceAt(i uint64, n *node) {
	ns.list[i] = n
}

func (ns *nodes) len() uint64 {
	return uint64(len(ns.list))
}

func newNodes(size uint64) *nodes {
	return &nodes{
		list: make(rtree.Rectangles, 0, size),
	}
}

type keys struct {
	list hilberts
}

func (ks *keys) splitAt(i, capacity uint64) (*keys, *keys) {
	i++
	right := make(hilberts, uint64(len(ks.list))-i, capacity)
	copy(right, ks.list[i:])
	ks.list = ks.list[:i]
	return ks, &keys{list: right}
}

func (ks *keys) len() uint64 {
	return uint64(len(ks.list))
}

func (ks *keys) byPosition(i uint64) hilbert {
	if i >= uint64(len(ks.list)) {
		return -1
	}
	return ks.list[i]
}

func (ks *keys) delete(k hilbert) hilbert {
	i := ks.search(k)
	if i >= uint64(len(ks.list)) {
		return -1
	}

	if ks.list[i] != k {
		return -1
	}
	old := ks.list[i]

	copy(ks.list[i:], ks.list[i+1:])
	ks.list = ks.list[:len(ks.list)-1]
	return old
}

func (ks *keys) search(key hilbert) uint64 {
	i := sort.Search(len(ks.list), func(i int) bool {
		return ks.list[i] >= key
	})

	return uint64(i)
}

func (ks *keys) insert(key hilbert) (hilbert, uint64) {
	i := ks.search(key)
	if i == uint64(len(ks.list)) {
		ks.list = append(ks.list, key)
		return -1, i
	}

	var old hilbert
	if ks.list[i] == key {
		old = ks.list[i]
		ks.list[i] = key
	} else {
		ks.insertAt(i, key)
	}

	return old, i
}

func (ks *keys) last() hilbert {
	return ks.list[len(ks.list)-1]
}

func (ks *keys) insertAt(i uint64, k hilbert) {
	ks.list = append(ks.list, -1)
	copy(ks.list[i+1:], ks.list[i:])
	ks.list[i] = k
}

func (ks *keys) withPosition(k hilbert) (hilbert, uint64) {
	i := ks.search(k)
	if i == uint64(len(ks.list)) {
		return -1, i
	}
	if ks.list[i] == k {
		return ks.list[i], i
	}

	return -1, i
}

func newKeys(size uint64) *keys {
	return &keys{
		list: make(hilberts, 0, size),
	}
}

type node struct {
	keys          *keys
	nodes         *nodes
	isLeaf        bool
	parent, right *node
	mbr           *rectangle
	maxHilbert    hilbert
}

func (n *node) insert(kb *keyBundle) rtree.Rectangle {
	i := n.keys.search(kb.key.hilbert)
	log.Printf(`I: %+v`, i)
	if n.isLeaf && i != n.keys.len() { // we can have multiple keys with the same hilbert number
		for n.keys.list[i] == kb.key.hilbert {
			if equal(n.nodes.list[i], kb.key.rect) {
				old := n.nodes.list[i]
				n.nodes.list[i] = kb.key.rect
				return old
			}
			i++
		}
	}

	if i == n.keys.len() {
		n.maxHilbert = kb.key.hilbert
	}

	n.keys.insertAt(i, kb.key.hilbert)
	if n.isLeaf {
		n.nodes.insertAt(i, kb.key.rect)
	} else {
		n.nodes.replaceAt(i, kb.left)
		n.nodes.insertAt(i+1, kb.right)
		n.mbr.adjust(kb.key.rect)
	}

	return nil
}

func (n *node) LowerLeft() (int32, int32) {
	return n.mbr.xlow, n.mbr.ylow
}

func (n *node) UpperRight() (int32, int32) {
	return n.mbr.xhigh, n.mbr.yhigh
}

func (n *node) needsSplit(ary uint64) bool {
	return n.keys.len() >= ary
}

func (n *node) splitLeaf(i, capacity uint64) (hilbert, *node, *node) {
	key := n.keys.byPosition(i)
	_, rightKeys := n.keys.splitAt(i, capacity)
	nn := &node{
		keys:   rightKeys,
		nodes:  newNodes(uint64(cap(n.nodes.list))),
		isLeaf: true,
		right:  n.right,
	}
	n.right = nn
	return key, n, nn
}

func (n *node) splitInternal(i, capacity uint64) (hilbert, *node, *node) {
	key := n.keys.byPosition(i)
	n.keys.delete(key)

	_, rightKeys := n.keys.splitAt(i-1, capacity)
	_, rightNodes := n.nodes.splitAt(i, capacity)

	nn := newNode(false, rightKeys, rightNodes)
	for _, n := range rightNodes.list {
		n.(*node).parent = nn
	}

	return key, n, nn
}

func (n *node) split(i, capacity uint64) (hilbert, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf(i, capacity)
	}

	return n.splitInternal(i, capacity)
}

func (n *node) search(key hilbert) uint64 {
	return n.keys.search(key)
}

func (n *node) searchNode(key hilbert) *node {
	i := n.search(key)

	return n.nodes.byPosition(uint64(i))
}

func (n *node) searchRects(r *rectangle) rtree.Rectangles {
	rects := make(rtree.Rectangles, 0, n.nodes.len())
	for _, child := range n.nodes.list {
		if intersect(r, child) {
			rects = append(rects, child)
		}
	}

	return rects
}

func (n *node) key() hilbert {
	return n.keys.last()
}

func newNode(isLeaf bool, keys *keys, ns *nodes) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
		nodes:  ns,
	}
}
