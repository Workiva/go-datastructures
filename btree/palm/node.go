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

	"github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/slice/skip"
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

func (ns *nodes) splitAt(i uint64) (*nodes, *nodes) {
	i++
	//log.Printf(`SPLITTING AT: %+v`, i)

	right := make([]*node, uint64(len(ns.list))-i, cap(ns.list))
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
	list *skip.SkipList
}

func (ks *keys) splitAt(i uint64) (*keys, *keys) {
	_, right := ks.list.SplitAt(i)
	return ks, &keys{list: right}
}

func (ks *keys) len() uint64 {
	return ks.list.Len()
}

func (ks *keys) byPosition(i uint64) common.Comparator {
	k, ok := ks.list.ByPosition(i).(common.Comparator)
	if !ok {
		return nil
	}

	return k
}

func (ks *keys) delete(k common.Comparator) {
	ks.list.Delete(k)
}

func (ks *keys) search(key common.Comparator) uint64 {
	n, i := ks.list.GetWithPosition(key)
	if n == nil {
		return ks.list.Len()
	}

	return i
}

func (ks *keys) insert(key common.Comparator) common.Comparator {
	old := ks.list.Insert(key)[0]
	if old == nil {
		return nil
	}

	return old
}

func (ks *keys) last() common.Comparator {
	return ks.list.ByPosition(ks.list.Len() - 1)
}

func (ks *keys) insertAt(i uint64, k common.Comparator) {
	ks.list.InsertAtPosition(i, k)
}

func (ks *keys) withPosition(k common.Comparator) (common.Comparator, uint64) {
	key, pos := ks.list.GetWithPosition(k)
	if key == nil {
		return nil, pos
	}

	return key, pos
}

func newKeys() *keys {
	return &keys{
		list: skip.New(uint32(0)),
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

func (n *node) splitLeaf() (common.Comparator, *node, *node) {
	i := n.keys.len() / 2
	key := n.keys.byPosition(i)
	_, rightKeys := n.keys.splitAt(i)
	nn := &node{
		keys:   rightKeys,
		nodes:  newNodes(uint64(cap(n.nodes.list))),
		isLeaf: true,
	}
	n.right = nn
	return key, n, nn
}

func (n *node) splitInternal() (common.Comparator, *node, *node) {
	i := n.keys.len() / 2
	key := n.keys.byPosition(i)
	n.keys.delete(key)

	_, rightKeys := n.keys.splitAt(i - 1)
	_, rightNodes := n.nodes.splitAt(i)

	nn := newNode(false, rightKeys, rightNodes)
	for _, n := range rightNodes.list {
		n.parent = nn
	}

	return key, n, nn
}

func (n *node) split() (common.Comparator, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf()
	}

	return n.splitInternal()
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
	for iter := n.keys.list.IterAtPosition(0); iter.Next(); {
		k := iter.Value()
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
