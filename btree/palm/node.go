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

	"github.com/Workiva/go-datastructures/slice/skip"
)

func getParent(parent *node, key Key) *node {
	var n *node
	for parent != nil && !parent.isLeaf {
		n = parent.searchNode(key)
		parent = n
	}

	return parent
}

type nodes struct {
	list *skip.SkipList
}

func (ns *nodes) push(n *node) {
	ns.list.InsertAtPosition(ns.list.Len(), n)
}

func (ns *nodes) splitAt(i uint64) (*nodes, *nodes) {
	_, right := ns.list.SplitAt(i)
	return ns, &nodes{list: right}
}

func (ns *nodes) byPosition(pos uint64) *node {
	n, ok := ns.list.ByPosition(pos).(*node)
	if !ok {
		return nil
	}

	return n
}

func (ns *nodes) insertAt(i uint64, n *node) {
	ns.list.InsertAtPosition(i, n)
}

func (ns *nodes) replaceAt(i uint64, n *node) {
	ns.list.ReplaceAtPosition(i, n)
}

func (ns *nodes) len() uint64 {
	return ns.list.Len()
}

func newNodes() *nodes {
	return &nodes{
		list: skip.New(uint64(0)),
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

func (ks *keys) byPosition(i uint64) Key {
	k, ok := ks.list.ByPosition(i).(Key)
	if !ok {
		return nil
	}

	return k
}

func (ks *keys) delete(k Key) {
	ks.list.Delete(k.(skip.Entry))
}

func (ks *keys) search(key Key) uint64 {
	n, i := ks.list.GetWithPosition(key.(skip.Entry))
	if n == nil {
		return ks.list.Len()
	}

	return i
}

func (ks *keys) insert(key Key) Key {
	old := ks.list.Insert(key)[0]
	if old == nil {
		return nil
	}

	return old.(Key)
}

func (ks *keys) last() Key {
	return ks.list.ByPosition(ks.list.Len() - 1).(Key)
}

func (ks *keys) insertAt(i uint64, k Key) {
	ks.list.InsertAtPosition(i, k.(skip.Entry))
}

func (ks *keys) withPosition(k Key) (Key, uint64) {
	key, pos := ks.list.GetWithPosition(k)
	if key == nil {
		return nil, pos
	}

	return key.(Key), pos
}

func newKeys() *keys {
	return &keys{
		list: skip.New(uint64(0)),
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

func (n *node) splitLeaf() (Key, *node, *node) {
	i := n.keys.len() / 2
	key := n.keys.byPosition(i)
	_, rightKeys := n.keys.splitAt(i)
	nn := &node{
		keys:   rightKeys,
		nodes:  newNodes(),
		isLeaf: true,
	}
	n.right = nn
	return key, n, nn
}

func (n *node) splitInternal() (Key, *node, *node) {
	i := n.keys.len() / 2
	key := n.keys.byPosition(i)
	n.keys.delete(key)

	_, rightKeys := n.keys.splitAt(i - 1)
	_, rightNodes := n.nodes.splitAt(i)

	nn := newNode(false, rightKeys, rightNodes)
	for iter := rightNodes.list.IterAtPosition(0); iter.Next(); {
		nd := iter.Value().(*node)
		nd.parent = nn
	}

	return key, n, nn
}

func (n *node) split() (Key, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf()
	}

	return n.splitInternal()
}

func (n *node) search(key Key) uint64 {
	return n.keys.search(key)
}

func (n *node) searchNode(key Key) *node {
	i := n.search(key)

	return n.nodes.byPosition(uint64(i))
}

func (n *node) key() Key {
	return n.keys.last()
}

func (n *node) print(output *log.Logger) {
	output.Printf(`NODE: %+v, %p`, n, n)
	for iter := n.keys.list.IterAtPosition(0); iter.Next(); {
		k := iter.Value().(Key)
		output.Printf(`KEY: %+v`, k)
	}
	if !n.isLeaf {
		for iter := n.nodes.list.IterAtPosition(0); iter.Next(); {
			n := iter.Value().(*node)
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
func (n *node) Compare(e skip.Entry) int {
	return 0
}

func newNode(isLeaf bool, keys *keys, ns *nodes) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
		nodes:  ns,
	}
}
