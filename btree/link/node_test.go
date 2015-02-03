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
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestNode(isLeaf bool, ary int) *node {
	return &node{
		isLeaf: isLeaf,
		keys:   make(Keys, 0, ary),
		nodes:  make(nodes, 0, ary+1),
	}
}

func checkTree(t testing.TB, tree *blink) {
	if tree.root == nil {
		return
	}

	checkNode(t, tree.root)
}

func checkNode(t testing.TB, n *node) {
	if len(n.keys) == 0 {
		assert.Len(t, n.nodes, 0)
		return
	}

	if n.isLeaf {
		assert.Len(t, n.nodes, 0)
		return
	}

	if !assert.Len(t, n.nodes, len(n.keys)+1) {
		return
	}

	for i := 0; i < len(n.keys); i++ {
		assert.Equal(t, n.keys[i], n.nodes[i+1].key())
	}

	assert.True(t, n.nodes[0].key().Compare(n.keys.first()) < 0)
	for i, child := range n.nodes {
		assert.NotNil(t, child)
		if child != nil {
			checkNode(t, child)
			if i != 0 {
				assert.Equal(t, n.keys[i-1], child.key())
			}
		}
	}
}

func TestSplitInternalNodeOddAry(t *testing.T) {
	parent := newTestNode(false, 3)
	n1 := newTestNode(true, 3)
	n1.keys.insert(mockKey(1))
	n2 := newTestNode(true, 3)
	n2.keys.insert(mockKey(5))
	n3 := newTestNode(true, 3)
	n3.keys.insert(mockKey(10))
	n4 := newTestNode(true, 3)
	n4.keys.insert(mockKey(15))

	parent.nodes = nodes{n1, n2, n3, n4}
	parent.keys = Keys{mockKey(5), mockKey(10), mockKey(15)}

	key, l, r := parent.split()
	assert.Equal(t, mockKey(10), key)
	assert.Equal(t, Keys{mockKey(5)}, l.keys)
	assert.Equal(t, Keys{mockKey(15)}, r.keys)

	assert.Equal(t, nodes{n1, n2}, l.nodes)
	assert.Equal(t, nodes{n3, n4}, r.nodes)
	assert.Equal(t, l.right, r)
	assert.False(t, l.isLeaf)
	assert.False(t, r.isLeaf)
	checkNode(t, l)
	checkNode(t, r)
}

func TestSplitInternalNodeEvenAry(t *testing.T) {
	parent := newTestNode(false, 4)
	n1 := newTestNode(true, 4)
	n1.keys.insert(mockKey(1))
	n2 := newTestNode(true, 4)
	n2.keys.insert(mockKey(5))
	n3 := newTestNode(true, 4)
	n3.keys.insert(mockKey(10))
	n4 := newTestNode(true, 4)
	n4.keys.insert(mockKey(15))
	n5 := newTestNode(true, 4)
	n5.keys.insert(mockKey(20))

	parent.nodes = nodes{n1, n2, n3, n4, n5}
	parent.keys = Keys{mockKey(5), mockKey(10), mockKey(15), mockKey(20)}

	key, l, r := parent.split()
	assert.Equal(t, mockKey(15), key)
	assert.Equal(t, Keys{mockKey(5), mockKey(10)}, l.keys)
	assert.Equal(t, Keys{mockKey(20)}, r.keys)

	assert.Equal(t, nodes{n1, n2, n3}, l.nodes)
	assert.Equal(t, nodes{n4, n5}, r.nodes)
	assert.Equal(t, l.right, r)
	assert.False(t, l.isLeaf)
	assert.False(t, r.isLeaf)
	checkNode(t, l)
	checkNode(t, r)
}

func TestSplitLeafNodeOddAry(t *testing.T) {
	parent := newTestNode(true, 3)
	k1 := mockKey(5)
	k2 := mockKey(15)
	k3 := mockKey(20)

	parent.keys = Keys{k1, k2, k3}
	key, l, r := parent.split()
	assert.Equal(t, k2, key)
	assert.Equal(t, Keys{k1}, l.keys)
	assert.Equal(t, Keys{k2, k3}, r.keys)
	assert.True(t, l.isLeaf)
	assert.True(t, r.isLeaf)
	assert.Equal(t, r, l.right)
	checkNode(t, l)
	checkNode(t, r)
}

func TestSplitLeafNodeEvenAry(t *testing.T) {
	parent := newTestNode(true, 4)
	k1 := mockKey(5)
	k2 := mockKey(15)
	k3 := mockKey(20)
	k4 := mockKey(25)

	parent.keys = Keys{k1, k2, k3, k4}
	key, l, r := parent.split()
	assert.Equal(t, k3, key)
	assert.Equal(t, Keys{k1, k2}, l.keys)
	assert.Equal(t, Keys{k3, k4}, r.keys)
	assert.True(t, l.isLeaf)
	assert.True(t, r.isLeaf)
	assert.Equal(t, r, l.right)
	checkNode(t, l)
	checkNode(t, r)
}
