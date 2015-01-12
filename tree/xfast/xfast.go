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
	entry    Entry
	children [2]*node
	parent   *node // i hate this, but it is really the best way
	// to walk up successor and predecessor threads
}

func newNode(parent *node, entry Entry) *node {
	return &node{
		children: [2]*node{},
		entry:    entry,
		parent:   parent,
	}
}

func binarySearchHashMaps(layers []map[uint64]*node, key uint64) (int, *node) {
	low, high := 1, len(layers)-1
	var mid int
	var node *node
	for low <= high {
		mid = (low + high) / 2
		n, ok := layers[mid][key&masks[mid]]
		if ok {
			node = n
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return low - 1, node
}

type XFastTrie struct {
	layers []map[uint64]*node
	root   *node
	num    uint64
	cache  []*node // we'll not need this for the leaf node
	bits   uint8
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
	for i := uint8(0); i < bits; i++ {
		xft.layers[i] = make(map[uint64]*node, 50) // we can obviously be more intelligent about this.
	}
	xft.num = 0
	xft.cache = make([]*node, bits)
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

func (xft *XFastTrie) insert(entry Entry) {
	key := entry.Key() // so we aren't calling this interface method over and over, fucking Go
	n := xft.layers[xft.bits-1][key]
	if n != nil {
		n.entry = entry
		return
	}

	var predecessor *node
	successor := xft.successor(key)
	if successor == nil {
		predecessor = xft.predecessor(key)
		if predecessor != nil {
			log.Printf(`PREDECESSOR: %+v`, predecessor.entry)
		} else {
			log.Printf(`NO PREDECCESOR FOUND`)
		}
	} else {
		log.Printf(`SUCCESSOR: %+v`, successor)
	}

	n = xft.root
	var leftOrRight uint64

	for i := uint8(0); i < xft.bits-1; i++ {
		// on 0th, this will be root
		xft.cache[i] = n
		// find out if we need to go left or right
		leftOrRight = (key & positions[i]) >> (xft.bits - 1 - i)
		if n.children[leftOrRight] == nil || isLeaf(n.children[leftOrRight]) {
			nn := newNode(n, nil)
			n.children[leftOrRight] = nn
			xft.layers[i][key&masks[i]] = nn // prefix for this layer
		}
		n = n.children[leftOrRight]
	}

	// we are left with next to last possible node
	leftOrRight = key & positions[xft.bits-1] // this will just be 1 or 0
	if n.children[leftOrRight] == nil || isLeaf(n.children[leftOrRight]) {
		leaf := newNode(n, entry)
		n.children[leftOrRight] = leaf
		xft.layers[xft.bits-1][key] = leaf
		xft.num++
		n = leaf
	} else {
		log.Printf(`WE HAVE A PROBLEM HERE WITH KEY: %+v`, key)
	}

	if successor != nil { // we have to walk predecessor and successor threads
		predecessor = successor.children[0]
		if predecessor != nil {
			predecessor.children[1] = n
			n.children[0] = predecessor
		}
		n.children[1] = successor
		successor.children[0] = n
	} else {
		if predecessor != nil {
			n.children[0] = predecessor
			predecessor.children[1] = n
		}
	}

	if successor != nil {
		xft.walkUpSuccessor(n, successor)
	}

	if predecessor != nil {
		xft.walkUpPredecessor(n, predecessor)
	}

	xft.walkUpNode(n)
}

func (xft *XFastTrie) walkUpSuccessor(node, successor *node) {
	i := uint8(0)
	n := successor.parent
	for n != nil && xft.cache[xft.bits-1-i] != n {
		if isLeaf(n.children[0]) {
			n.children[0] = node
		}
		i++
		if i > xft.bits-1 {
			break
		}
		n = n.parent
	}
}

func (xft *XFastTrie) walkUpPredecessor(node, predecessor *node) {
	i := uint8(0)
	n := predecessor.parent
	for n != nil && xft.cache[xft.bits-1-i] != n {
		if isLeaf(n.children[1]) {
			n.children[1] = node
		}
		i++
		if i > xft.bits-1 {
			break
		}
		n = n.parent
	}
}

func (xft *XFastTrie) walkUpNode(node *node) {
	n := node.parent
	for n != nil {
		if n.children[0] == nil {
			n.children[0] = node
		}
		if n.children[1] == nil {
			n.children[1] = node
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
	if xft.root == nil { // no successor if no nodes
		return nil
	}

	n := xft.layers[xft.bits-1][key]
	if n != nil {
		return n
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
	log.Printf(`LAYER: %+v, N: %+v`, layer, n)
	if n == nil && layer > 1 {
		return nil
	} else if n == nil {
		n = xft.root
	}

	log.Printf(`N0: %+v, N1: %+v`, n.children[0], n.children[1])

	if isLeaf(n.children[1]) {
		if n.children[1].entry.Key() <= key {
			return n.children[1]
		} else if isLeaf(n.children[0]) {
			println(`THIS FUCKING WAY`)
			log.Printf(`n.children[1].entry: %+v`, n.children[1].entry)
			if n.children[0].entry.Key() <= key {
				return n.children[0]
			}
		}

		return n.children[1].children[0]
	} else if isLeaf(n.children[0]) {
		println(`HERE`)
		log.Printf(`n.children[0].entry: %+v`, n.children[0].entry)
		if n.children[0].entry.Key() <= key {
			return n.children[0]
		}
	}

	return nil
}

func (xft *XFastTrie) successor(key uint64) *node {
	if xft.root == nil { // no successor if no nodes
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

	if isLeaf(n.children[0]) {
		if n.children[0].entry.Key() >= key {
			return n.children[0]
		} else if isLeaf(n.children[1]) {
			if n.children[1].entry.Key() >= key {
				return n.children[1]
			}
		}

		return n.children[0].children[1]
	} else if isLeaf(n.children[1]) {
		if n.children[1].entry.Key() >= key {
			return n.children[1]
		}
	}

	return nil
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

// New will construct a new X-Fast Trie with the given "size,"
// that is the size of the universe of the trie.  This expects
// a uint of some sort, ie, uint8, uint16, etc.  The size of the
// universe will be 2^n-1 and will affect the speed of all operations.
func New(ifc interface{}) *XFastTrie {
	xft := &XFastTrie{}
	xft.init(ifc)
	return xft
}
