package xfast

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Workiva/go-datastructures/slice"
)

func init() {
	log.Printf(`I HATE THIS`)
}

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
		if hasSuccesor {
			assert.Equal(t, n, successor.children[0])
		}

		for n.parent != nil {
			side = whichSide(n, n.parent)
			if isInternal(n.parent.children[1]) && isInternal(n.parent.children[0]) {
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
		if hasPredecessor {
			assert.Equal(t, n, predecessor.children[1])
		}
		for n.parent != nil {
			side = whichSide(n, n.parent)
			if isInternal(n.parent.children[0]) && isInternal(n.parent.children[1]) {
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

	println(`SHIT STARTS HERE`)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	assert.True(t, xft.Exists(20))
	assert.Equal(t, uint64(2), xft.Len())
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e2, xft.Max())
	checkTrie(t, xft)
}

func TestInsertBetween(t *testing.T) {
	xft := New(uint8(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	assert.True(t, xft.Exists(10))
	assert.Equal(t, e1, xft.Min())
	assert.Equal(t, e1, xft.Max())
	checkPattern(t, xft.min, []int{0, 0, 0, 0, 1, 0, 1, 0})
	checkTrie(t, xft)

	e2 := newMockEntry(20)
	xft.Insert(e2)
	checkPattern(t, xft.max, []int{0, 0, 0, 1, 0, 1, 0, 0})
	checkPattern(t, xft.min, []int{0, 0, 0, 0, 1, 0, 1, 0})
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

	for i := uint64(16); i < 17; i++ {
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
	xft := New(uint64(0))
	e1 := newMockEntry(10)
	xft.Insert(e1)

	e2 := newMockEntry(20)
	xft.Insert(e2)

	for i := uint64(16); i < 17; i++ {
		result := xft.Predecessor(i)
		assert.Equal(t, e1, result)
	}

	for i := uint64(0); i < 10; i++ {
		result := xft.Predecessor(i)
		assert.Nil(t, result)
	}
}

func BenchmarkSuccessor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xft := New(uint64(0))
		e := newMockEntry(uint64(i))
		xft.Insert(e)
		xft.Successor(0)
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
