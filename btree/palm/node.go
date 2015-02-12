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
	isLeaf        bool
	parent, right *node
}

func (n *node) needsSplit(ary uint64) bool {
	return n.keys.len() >= ary
}

func (n *node) splitLeaf(i uint64) (Key, *node, *node) {
	key := n.keys.byPosition(i)
	leftKeys, rightKeys := n.keys.splitAt(i - 1)
	nn := &node{
		keys:   leftKeys,
		isLeaf: true,
	}
	n.keys = rightKeys
	nn.right = n
	return key, nn, n
}

func (n *node) max() Key {
	if n.isLeaf {
		return n.keys.last()
	}

	return n.keys.last().(*node).max()
}

func (n *node) splitInternal(i uint64) (Key, *node, *node) {
	key := n.keys.byPosition(i)
	//n.keys.delete(key)

	leftKeys, rightKeys := n.keys.splitAt(i - 1)

	n.keys = rightKeys
	nn := newNode(false, leftKeys)
	return key, nn, n
}

func (n *node) split(i uint64) (Key, *node, *node) {
	if n.isLeaf {
		return n.splitLeaf(i)
	}

	return n.splitInternal(i)
}

func (n *node) search(key Key) uint64 {
	return n.keys.search(key)
}

func (n *node) searchNode(key Key) *node {
	// TODO: add successor search to skiplist to improve performance here
	i := n.search(key)
	if i == n.keys.len() {
		i--
	}
	result := n.keys.byPosition(i)

	if result == nil {
		return nil
	}

	return result.(*node)
}

func (n *node) key() Key {
	return n.keys.last()
}

func (n *node) iter() {
	for iter := n.keys.list.IterAtPosition(0); iter.Next(); {
		log.Printf(`ITER.VALUE: %+v, %p`, iter.Value(), iter.Value())
	}
}

func (n *node) print(output *log.Logger) {
	output.Printf(`NODE: %+v, %p`, n, n)
	for iter := n.keys.list.IterAtPosition(0); iter.Next(); {
		k := iter.Value().(Key)
		output.Printf(`KEY: %+v`, k)
	}
	if !n.isLeaf {
		for iter := n.keys.list.IterAtPosition(0); iter.Next(); {
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
	var key Key
	switch e.(type) {
	case *node:
		if e.(*node).isLeaf {
			key = e.(*node).keys.last()
		} else {
			key = e.(*node).max()
		}
	default:
		key = e.(Key)
	}

	if n.isLeaf {
		return n.keys.last().Compare(key)
	}
	return n.max().Compare(key)
}

func newNode(isLeaf bool, keys *keys) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   keys,
	}
}
