package xfast

import (
	"log"
)

// isInternal returns a bool indicating if the provided
// node is an internal node, that is, non-leaf node.
func isInternal(n *node) bool {
	if n == nil {
		return false
	}
	return n.entry == nil
}

// isLeaf returns a bool indicating if the provided node
// is a leaf node, that is, has a valid entry value.
func isLeaf(n *node) bool {
	if n == nil {
		return false
	}
	return !isInternal(n)
}

// Entry defines items that can be inserted into the x-fast
// trie.
type Entry interface {
	// Key is the key for this entry.
	Key() uint64
}

func init() {
	log.Printf(`I HATE THIS`)
}

var masks = func() [64]uint64 { // we don't technically need the last mask, this is just to be consistent
	masks := [64]uint64{}
	mask := uint64(0)
	for i := uint64(0); i < 64; i++ {
		mask = mask | 1<<(63-i)
		masks[i] = mask
	}
	return masks
}()

var positions = func() [64]uint64 {
	positions := [64]uint64{}
	for i := uint64(0); i < 64; i++ {
		positions[i] = uint64(1 << (63 - i))
	}
	return positions
}()

type node struct {
	// entry will
	entry    Entry
	children [2]*node
	// i hate this, but it is really the best way
	// to walk up successor and predecessor threads
	parent *node
}

func newNode(parent *node, entry Entry) *node {
	return &node{
		children: [2]*node{},
		entry:    entry,
		parent:   parent,
	}
}

func binarySearchHashMaps(layers []map[uint64]*node, key uint64) (int, *node) {
	low, high := 0, len(layers)-1
	diff := 64 - len(layers)
	var mid int
	var node *node
	for low <= high {
		mid = (low + high) / 2
		n, ok := layers[mid][key&masks[diff+mid]]
		if ok {
			node = n
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low, node
}

type XFastTrie struct {
	layers     []map[uint64]*node
	root       *node
	num        uint64
	bits, diff uint8
	min, max   *node
}

func (xft *XFastTrie) init(intType interface{}) {
	bits := uint8(0)
	switch intType.(type) {
	case uint8:
		bits = 8
	case uint16:
		bits = 16
	case uint32:
		bits = 32
	case uint, uint64:
		bits = 64
	default:
		panic(`Invalid universe size provided.`)
	}

	xft.layers = make([]map[uint64]*node, bits)
	xft.bits = bits
	xft.diff = 64 - bits
	for i := uint8(0); i < bits; i++ {
		xft.layers[i] = make(map[uint64]*node, 50) // we can obviously be more intelligent about this.
	}
	xft.num = 0
	xft.root = newNode(nil, nil)
}

// Exists returns a bool indicating if the provided
// key exists in the trie.  This is typically an
// O(1) operation.
func (xft *XFastTrie) Exists(key uint64) bool {
	// the bottom hashmap of the trie has every entry
	// in it.
	_, ok := xft.layers[xft.bits-1][key]
	return ok
}

// Len returns the number of items in this trie.
func (xft *XFastTrie) Len() uint64 {
	return xft.num
}

func (xft *XFastTrie) Max() Entry {
	if xft.max == nil {
		return nil
	}

	return xft.max.entry
}

func (xft *XFastTrie) Min() Entry {
	if xft.min == nil {
		return nil
	}

	return xft.min.entry
}

func (xft *XFastTrie) insert(entry Entry) {
	key := entry.Key() // so we aren't calling this interface method over and over, fucking Go
	n := xft.layers[xft.bits-1][key]
	if n != nil {
		n.entry = entry
		return
	}

	var predecessor, successor *node
	if xft.min != nil && key < xft.min.entry.Key() {
		successor = xft.min
	} else {
		successor = xft.successor(key)
	}

	if successor == nil {
		if xft.max != nil && key > xft.max.entry.Key() {
			predecessor = xft.max
		} else {
			predecessor = xft.predecessor(key)
		}
	}

	layer, root := binarySearchHashMaps(xft.layers, key)
	if root == nil {
		n = xft.root
		layer = 0
	} else {
		n = root
	}

	var leftOrRight uint64

	for i := uint8(layer); i < xft.bits; i++ {
		// on 0th, this will be root
		// find out if we need to go left or right
		leftOrRight = (key & positions[xft.diff+i]) >> (xft.bits - 1 - i)
		if n.children[leftOrRight] == nil || isLeaf(n.children[leftOrRight]) {
			var nn *node
			if i < xft.bits-1 {
				nn = newNode(n, nil)
			} else {
				nn = newNode(n, entry)
				xft.num++
			}

			n.children[leftOrRight] = nn
			xft.layers[i][key&masks[xft.diff+i]] = nn // prefix for this layer
		}

		n = n.children[leftOrRight]
	}

	if successor != nil { // we have to walk predecessor and successor threads
		predecessor = successor.children[0]
		if predecessor != nil {
			predecessor.children[1] = n
			n.children[0] = predecessor
		}
		n.children[1] = successor
		successor.children[0] = n
	} else if predecessor != nil {
		n.children[0] = predecessor
		predecessor.children[1] = n
	}

	if successor != nil {
		xft.walkUpSuccessor(root, n, successor)
	}

	if predecessor != nil {
		xft.walkUpPredecessor(root, n, predecessor)
	}

	xft.walkUpNode(root, n, predecessor, successor)

	if xft.max == nil || key > xft.max.entry.Key() {
		xft.max = n
	}

	if xft.min == nil || key < xft.min.entry.Key() {
		xft.min = n
	}
}

func (xft *XFastTrie) walkUpSuccessor(root, node, successor *node) {
	n := successor.parent
	for n != nil && n != root {
		if !isInternal(n.children[0]) && n.children[0] != successor {
			n.children[0] = node
		}
		n = n.parent
	}
}

func (xft *XFastTrie) walkUpPredecessor(root, node, predecessor *node) {
	n := predecessor.parent
	for n != nil && n != root {
		if !isInternal(n.children[1]) && n.children[1] != predecessor {
			n.children[1] = node
		}
		n = n.parent
	}
}

func (xft *XFastTrie) walkUpNode(root, node, predecessor, successor *node) {
	n := node.parent
	for n != nil && n != root {
		if !isInternal(n.children[1]) && n.children[1] != successor && n.children[1] != node {
			n.children[1] = successor
		}
		if !isInternal(n.children[0]) && n.children[0] != predecessor && n.children[0] != node {
			n.children[0] = predecessor
		}
		n = n.parent
	}
}

func (xft *XFastTrie) Insert(entries ...Entry) {
	for _, e := range entries {
		xft.insert(e)
	}
}

func (xft *XFastTrie) predecessor(key uint64) *node {
	if xft.root == nil || xft.max == nil { // no successor if no nodes
		return nil
	}

	if key >= xft.max.entry.Key() {
		return xft.max
	}

	if key < xft.min.entry.Key() {
		return nil
	}

	n := xft.layers[xft.bits-1][key]
	if n != nil {
		return n
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
	if n == nil && layer > 1 {
		return nil
	} else if n == nil {
		n = xft.root
	}

	if isInternal(n.children[0]) && isLeaf(n.children[1]) {
		return n.children[1].children[0]
	}
	return n.children[0]
}

func (xft *XFastTrie) successor(key uint64) *node {
	if xft.root == nil || xft.min == nil { // no successor if no nodes
		return nil
	}

	if key <= xft.min.entry.Key() {
		return xft.min
	}

	if key > xft.max.entry.Key() {
		return nil
	}

	n := xft.layers[xft.bits-1][key]
	if n != nil {
		return n
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
	if n == nil && layer > 1 {
		return nil
	} else if n == nil {
		n = xft.root
	}

	if isInternal(n.children[1]) && isLeaf(n.children[0]) {
		return n.children[0].children[1]
	}
	return n.children[1]
}

func (xft *XFastTrie) Successor(key uint64) Entry {
	n := xft.successor(key)
	if n == nil {
		return nil
	}

	return n.entry
}

func (xft *XFastTrie) Predecessor(key uint64) Entry {
	n := xft.predecessor(key)
	if n == nil {
		return nil
	}

	return n.entry
}

func (xft *XFastTrie) Iter(key uint64) *Iterator {
	return &Iterator{
		n:     xft.successor(key),
		first: true,
	}
}

// New will construct a new X-Fast Trie with the given "size,"
// that is the size of the universe of the trie.  This expects
// a uint of some sort, ie, uint8, uint16, etc.  The size of the
// universe will be 2^n-1 and will affect the speed of all operations.
func New(ifc interface{}) *XFastTrie {
	xft := &XFastTrie{}
	xft.init(ifc)
	return xft
}
