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

package plus

func keySearch(keys keys, key Key) int {
	low, high := 0, len(keys)-1
	var mid int
	for low <= high {
		mid = (high + low) / 2
		switch keys[mid].Compare(key) {
		case 1:
			low = mid + 1
		case -1:
			high = mid - 1
		case 0:
			return mid
		}
	}
	return low
}

type btree struct {
	root             node
	nodeSize, number uint64
}

func (tree *btree) insert(key Key) {
	if tree.root == nil {
		n := newLeafNode(tree.nodeSize)
		n.insert(tree, key)
		tree.number = 1
		return
	}

	result := tree.root.insert(tree, key)
	if result {
		tree.number++
	}

	if tree.root.needsSplit(tree.nodeSize) {
		tree.root = split(tree, nil, tree.root)
	}
}

// Insert will insert the provided keys into the btree.  This is an
// O(m*log n) operation where m is the number of keys to be inserted
// and n is the number of items in the tree.
func (tree *btree) Insert(keys ...Key) {
	for _, key := range keys {
		tree.insert(key)
	}
}

// Iter returns an iterator that can be used to traverse the b-tree
// starting from the specified key or its successor.
func (tree *btree) Iter(key Key) Iterator {
	if tree.root == nil {
		return nilIterator()
	}

	return tree.root.find(key)
}

func (tree *btree) get(key Key) Key {
	iter := tree.root.find(key)
	if !iter.Next() {
		return nil
	}

	if iter.Value().Compare(key) == 0 {
		return iter.Value()
	}

	return nil
}

// Get will retrieve any keys matching the provided keys in the tree.
// Returns nil in any place of a key that couldn't be found.  Each lookup
// is an O(log n) operation.
func (tree *btree) Get(keys ...Key) Keys {
	results := make(Keys, 0, len(keys))
	for _, k := range keys {
		results = append(results, tree.get(k))
	}

	return results
}

// Len returns the number of items in this tree.
func (tree *btree) Len() uint64 {
	return tree.number
}

func newBTree(nodeSize uint64) *btree {
	return &btree{
		nodeSize: nodeSize,
		root:     newLeafNode(nodeSize),
	}
}
