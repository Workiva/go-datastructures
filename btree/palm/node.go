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

import (
	"log"
	"sort"

	"github.com/Workiva/go-datastructures/common"
)

func getParent(parent *node, key common.Comparator) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		n = parent.searchNode(key)
		parent = n
	}

	return parent
}

type nodes struct {
	list []*node
}

func (ns *nodes) push(n *node) {
	ns.list = append(ns.list, n)
}

func (ns *nodes) splitAt(i, capacity uint64) (*nodes, *nodes) {
	i++
	//log.Printf(`SPLITTING AT: %+v`, i)

	right := make([]*node, uint64(len(ns.list))-i, capacity)
	copy(right, ns.list[i:])
	//log.Printf(`NS.LIST: %+v, RIGHT: %+v`, ns.list, right)
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

	return ns.list[pos]
}

func (ns *nodes) insertAt(i uint64, n *node) {
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
		list: make([]*node, 0, size),
	}
}

type keys struct {
	list common.Comparators
}

func (ks *keys) splitAt(i, capacity uint64) (*keys, *keys) {
	i++
	right := make(common.Comparators, uint64(len(ks.list))-i, capacity)
	copy(right, ks.list[i:])
	for j := i; j < uint64(len(ks.list)); j++ {
		ks.list[j] = nil
	}
	ks.list = ks.list[:i]
	return ks, &keys{list: right}
}

func (ks *keys) len() uint64 {
	return uint64(len(ks.list))
}

func (ks *keys) byPosition(i uint64) common.Comparator {
	if i >= uint64(len(ks.list)) {
		return nil
	}
	return ks.list[i]
}

func (ks *keys) delete(k common.Comparator) {
	i := ks.search(k)
	if i >= uint64(len(ks.list)) {
		return
	}
	copy(ks.list[i:], ks.list[i+1:])
	ks.list[len(ks.list)-1] = nil // GC
	ks.list = ks.list[:len(ks.list)-1]
}

func (ks *keys) search(key common.Comparator) uint64 {
	i := sort.Search(len(ks.list), func(i int) bool {
		return ks.list[i].Compare(key) > -1
	})

	return uint64(i)
}

func (ks *keys) insert(key common.Comparator) common.Comparator {
	i := ks.search(key)
	if i == uint64(len(ks.list)) {
		ks.list = append(ks.list, key)
		return nil
	}

	old := ks.list[i]
	if ks.list[i].Compare(key) == 0 {
		ks.list[i] = key
	} else {
		ks.insertAt(i, key)
	}

	return old
}

func (ks *keys) last() common.Comparator {
	return ks.list[len(ks.list)-1]
}

func (ks *keys) insertAt(i uint64, k common.Comparator) {
	ks.list = append(ks.list, nil)
	copy(ks.list[i+1:], ks.list[i:])
	ks.list[i] = k
}

func (ks *keys) withPosition(k common.Comparator) (common.Comparator, uint64) {
	i := ks.search(k)
	if i == uint64(len(ks.list)) {
		return nil, i
	}
	if ks.list[i].Compare(k) == 0 {
		return ks.list[i], i
	}

	return nil, i
}

func newKeys(size uint64) *keys {
	return &keys{
		list: make(common.Comparators, 0, size),
	}
}

type node struct {
	keys          *keys
	nodes         *nodes
	isLeaf        bool
	parent, right *node
}

func (n *node) needsSplit(ary uint64) bool {
	return n.keys.len() >= ary
}

func (n *node) splitLeaf(i, capacity uint64) (common.Comparator, *node, *node) {
	key := n.keys.byPosition(i)
	_, rightKeys := n.keys.splitAt(i, capacity)
	nn := &node{
		keys:   rightKeys,
		nodes:  newNodes(uint64(cap(n.nodes.list))),
		isLeaf: true,
	}
	n.right = nn
	return key, n, nn
}

func (n *node) splitInternal(i, capacity uint64) (common.Comparator, *node, *node) {
	key := n.keys.byPosition(i)
	n.keys.delete(key)

	_, rightKeys := n.keys.splitAt(i-1, capacity)
	_, rightNodes := n.nodes.splitAt(i, capacity)

	nn := newNode(false, rightKeys, rightNodes)
	for _, n := range rightNodes.list {
		n.parent = nn
	}

	return key, n, nn
}

func (n *node) split(i, capacity uint64) (common.Comparator, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf(i, capacity)
	}

	return n.splitInternal(i, capacity)
}

func (n *node) search(key common.Comparator) uint64 {
	return n.keys.search(key)
}

func (n *node) searchNode(key common.Comparator) *node {
	i := n.search(key)

	return n.nodes.byPosition(uint64(i))
}

func (n *node) key() common.Comparator {
	return n.keys.last()
}

func (n *node) print(output *log.Logger) {
	output.Printf(`NODE: %+v, %p`, n, n)
	for _, k := range n.keys.list {
		output.Printf(`KEY: %+v`, k)
	}
	if !n.isLeaf {
		for _, n := range n.nodes.list {
			if n == nil {
				output.Println(`NIL NODE`)
				continue
			}

			n.print(output)
		}
	}
}

// Compare is required by the skip.Entry interface but nodes are always
// added by position so while this method is required it doesn't
// need to return anything useful.
func (n *node) Compare(e common.Comparator) int {
	return 0
}

func newNode(isLeaf bool, keys *keys, ns *nodes) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
		nodes:  ns,
	}
}
