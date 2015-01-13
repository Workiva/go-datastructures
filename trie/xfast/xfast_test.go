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

package xfast

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/slice"
)

func checkTrie(t *testing.T, xft *XFastTrie) {
	checkSuccessor(t, xft)
	checkPredecessor(t, xft)
	checkNodes(t, xft)
}

func checkSuccessor(t *testing.T, xft *XFastTrie) {
	n := xft.min
	var side int
	var successor *node
	for n != nil {
		successor = n.children[1]
		hasSuccesor := successor != nil
		immediateSuccessor := false
		if hasSuccesor {
			assert.Equal(t, n, successor.children[0])
			if n.parent == successor.parent {
				immediateSuccessor = true
			}
		}

		for n.parent != nil {
			side = whichSide(n, n.parent)
			if isInternal(n.parent.children[1]) && isInternal(n.parent.children[0]) {
				break
			}
			if immediateSuccessor && n.parent == successor.parent {
				assert.Equal(t, successor, n.parent.children[1])
				break
			}
			if side == 0 && !isInternal(n.parent.children[1]) && hasSuccesor {
				assert.Equal(t, successor, n.parent.children[1])
			}
			n = n.parent
		}
		n = successor
	}
}

func checkPredecessor(t *testing.T, xft *XFastTrie) {
	n := xft.max
	var side int
	var predecessor *node
	for n != nil {
		predecessor = n.children[0]
		hasPredecessor := predecessor != nil
		immediatePredecessor := false
		if hasPredecessor {
			assert.Equal(t, n, predecessor.children[1])
			if n.parent == predecessor.parent {
				immediatePredecessor = true
			}
		}
		for n.parent != nil {
			side = whichSide(n, n.parent)
			if isInternal(n.parent.children[0]) && isInternal(n.parent.children[1]) {
				break
			}
			if immediatePredecessor && n.parent == predecessor.parent {
				assert.Equal(t, predecessor, n.parent.children[0])
				break
			}
			if side == 1 && !isInternal(n.parent.children[0]) && hasPredecessor {
				assert.Equal(t, predecessor, n.parent.children[0])
			}
			n = n.parent
		}
		n = predecessor
	}
}

func checkNodes(t *testing.T, xft *XFastTrie) {
	count := uint64(0)
	n := xft.min
	for n != nil {
		count++
		checkNode(t, xft, n)
		n = n.children[1]
	}

	assert.Equal(t, count, xft.Len())
}

func checkNode(t *testing.T, xft *XFastTrie, n *node) {
	if n.entry == nil {
		assert.Fail(t, `Expected non-nil entry`)
		return
	}
	key := n.entry.Key()
	bits := make([]int, 0, xft.bits)
	for i := uint8(0); i < xft.bits; i++ {
		leftOrRight := (key & positions[xft.diff+i]) >> (xft.bits - 1 - i)
		bits = append(bits, int(leftOrRight))
	}

	checkPattern(t, n, bits)
}

func dumpNode(t *testing.T, n *node) {
	for n != nil {
		t.Logf(`NODE: %+v, %p`, n, n)
		n = n.parent
	}
}

func checkPattern(t *testing.T, n *node, pattern []int) {
	i := len(pattern) - 1
	bottomNode := n
	for n.parent != nil {
		if !assert.False(t, i < 0, fmt.Sprintf(`Too many parents. NODE: %+v, PATTERN: %+v`, bottomNode, pattern)) {
			dumpNode(t, bottomNode)
			break // so we don't panic on the next line
		}
		assert.Equal(t, pattern[i], whichSide(n, n.parent))
		i--
		n = n.parent
	}

	assert.Equal(t, -1, i)
}

func TestEmptyMinMax(t *testing.T) {
	xft := New(uint8(0))

	assert.Nil(t, xft.Min())
	assert.Nil(t, xft.Max())
}

func TestMask(t *testing.T) {
	assert.Equal(t, uint64(math.MaxUint64), masks[63])
}

func TestInsert(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	assert.True(t, xft.Exists(5))
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e1, xft.Max())
	checkTrie(t, xft)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	assert.True(t, xft.Exists(20))
	assert.Equal(t, uint64(2), xft.Len())
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)
}

func TestGet(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	assert.Equal(t, e1, xft.Get(5))
	assert.Nil(t, xft.Get(6))
}

func TestInsertOverwrite(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	e2 := newMockEntry(5)
	xft.Insert(e2)
	checkTrie(t, xft)

	iter := xft.Iter(5)
	assert.Equal(t, Entries{e2}, iter.exhaust())
}

func TestInsertBetween(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	assert.True(t, xft.Exists(10))
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e1, xft.Max())
	checkTrie(t, xft)

	e2 := newMockEntry(20)
	xft.Insert(e2)
	checkTrie(t, xft)

	assert.True(t, xft.Exists(20))
	assert.Equal(t, uint64(2), xft.Len())
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())

	assert.Equal(t, e2, xft.Successor(15))

	e3 := newMockEntry(15)
	xft.Insert(e3)

	assert.True(t, xft.Exists(15))
	assert.Equal(t, uint64(3), xft.Len())
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)

	iter := xft.Iter(0)
	entries := iter.exhaust()
	assert.Equal(t, Entries{e1, e3, e2}, entries)

	iter = xft.Iter(11)
	entries = iter.exhaust()
	assert.Equal(t, Entries{e3, e2}, entries)

	iter = xft.Iter(16)
	entries = iter.exhaust()
	assert.Equal(t, Entries{e2}, entries)
}

func TestSuccessorDoesNotExist(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Successor(6)
	assert.Nil(t, result)
}

func TestSuccessorIsExactValue(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorGreaterThanKey(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(math.MaxUint8)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorCloseToKey(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	result := xft.Successor(5)
	assert.Equal(t, e1, result)
}

func TestSuccessorBetweenTwoKeys(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	for i := uint64(11); i < 20; i++ {
		result := xft.Successor(i)
		assert.Equal(t, e2, result)
	}

	for i := uint64(21); i < 100; i++ {
		result := xft.Successor(i)
		assert.Nil(t, result)
	}
}

func TestPredecessorDoesNotExist(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(4)
	assert.Nil(t, result)
}

func TestPredecessorIsExactValue(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(5)
	assert.Equal(t, e1, result)
}

func TestPredecessorLessThanKey(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	xft.Insert(e1)

	result := xft.Predecessor(math.MaxUint64)
	assert.Equal(t, e1, result)
}

func TestPredecessorCloseToKey(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(5)
	xft.Insert(e1)

	result := xft.Predecessor(10)
	assert.Equal(t, e1, result)
}

func TestPredecessorBetweenTwoKeys(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	for i := uint64(11); i < 20; i++ {
		result := xft.Predecessor(i)
		assert.Equal(t, e1, result)
	}

	for i := uint64(0); i < 10; i++ {
		result := xft.Predecessor(i)
		assert.Nil(t, result)
	}
}

func TestInsertPredecessor(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(5)
	xft.Insert(e2)
	checkTrie(t, xft)

	assert.Equal(t, e2, xft.Min())
	assert.Equal(t, e1, xft.Max())

	iter := xft.Iter(2)
	assert.Equal(t, Entries{e2, e1}, iter.exhaust())

	iter = xft.Iter(5)
	assert.Equal(t, Entries{e2, e1}, iter.exhaust())

	iter = xft.Iter(6)
	assert.Equal(t, Entries{e1}, iter.exhaust())

	iter = xft.Iter(11)
	assert.Equal(t, Entries{}, iter.exhaust())
}

func TestDeleteOnlyBranch(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	xft.Delete(10)
	checkTrie(t, xft)
	assert.Equal(t, uint64(0), xft.Len())
	assert.Nil(t, xft.Min())
	assert.Nil(t, xft.Max())
	for _, hm := range xft.layers {
		assert.Len(t, hm, 0)
	}

	assert.NotNil(t, xft.root)
	assert.Nil(t, xft.root.children[0])
	assert.Nil(t, xft.root.children[1])

	iter := xft.Iter(0)
	assert.False(t, iter.Next())
}

func TestDeleteLargeBranch(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	e2 := newMockEntry(math.MaxUint8)

	xft.Insert(e1, e2)
	checkTrie(t, xft)

	xft.Delete(0)
	assert.Equal(t, e2, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)

	assert.Nil(t, xft.root.children[0])

	n := xft.max
	for n != nil {
		assert.Nil(t, n.children[0])
		n = n.parent
	}

	iter := xft.Iter(0)
	assert.Equal(t, Entries{e2}, iter.exhaust())
}

func TestDeleteLateBranching(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	e2 := newMockEntry(1)

	xft.Insert(e1, e2)
	checkTrie(t, xft)

	xft.Delete(1)
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e1, xft.Max())
	checkTrie(t, xft)

	n := xft.min
	for n != nil {
		assert.Nil(t, n.children[1])
		n = n.parent
	}

	iter := xft.Iter(0)
	assert.Equal(t, Entries{e1}, iter.exhaust())
}

func TestDeleteLateBranchingMin(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	e2 := newMockEntry(1)

	xft.Insert(e1, e2)
	checkTrie(t, xft)

	xft.Delete(0)
	assert.Equal(t, e2, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)

	assert.Nil(t, xft.min.children[0])
	n := xft.min.parent
	assert.Nil(t, n.children[0])
	n = n.parent
	for n != nil {
		assert.Nil(t, n.children[1])
		n = n.parent
	}

	iter := xft.Iter(0)
	assert.Equal(t, Entries{e2}, iter.exhaust())
}

func TestDeleteMiddleBranch(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	e2 := newMockEntry(math.MaxUint8)
	e3 := newMockEntry(64) // [0, 1, 0, 0, 0, 0, 0, 0]

	xft.Insert(e1, e2, e3)
	checkTrie(t, xft)

	xft.Delete(64)
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)

	iter := xft.Iter(0)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())
}

func TestDeleteMiddleBranchOtherSide(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(0)
	e2 := newMockEntry(math.MaxUint8)
	e3 := newMockEntry(128) // [1, 0, 0, 0, 0, 0, 0, 0]

	xft.Insert(e1, e2, e3)
	checkTrie(t, xft)

	xft.Delete(128)
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)

	iter := xft.Iter(0)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())
}

func TestDeleteNotFound(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(64)
	xft.Insert(e1)
	checkTrie(t, xft)

	xft.Delete(128)
	assert.Equal(t, e1, xft.Max())
	assert.Equal(t, e1, xft.Min())
	checkTrie(t, xft)
}

func BenchmarkSuccessor(b *testing.B) {
	numItems := 10000
	xft := New(uint64(0))

	for i := uint64(0); i < uint64(numItems); i++ {
		xft.Insert(newMockEntry(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		xft.Successor(uint64(i))
	}
}

func BenchmarkDelete(b *testing.B) {
	xs := make([]*XFastTrie, 0, b.N)

	for i := 0; i < b.N; i++ {
		x := New(uint8(0))
		x.Insert(newMockEntry(0))
		xs = append(xs, x)
	}

	// this is actually a pretty bad case for the x-fast
	// trie as the entire branch will have to be walked.
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		xs[i].Delete(0)
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xft := New(uint64(0))
		e := newMockEntry(uint64(i))
		xft.Insert(e)
	}
}

// benchmarked against a flat list
func BenchmarkListInsert(b *testing.B) {
	numItems := 100000

	s := make(slice.Int64Slice, 0, numItems)
	for j := int64(0); j < int64(numItems); j++ {
		s = append(s, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Insert(int64(i))
	}
}

func BenchmarkListSearch(b *testing.B) {
	numItems := 1000000

	s := make(slice.Int64Slice, 0, numItems)
	for j := int64(0); j < int64(numItems); j++ {
		s = append(s, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Search(int64(i))
	}
}
