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

func binarySearchHashMaps(layers [64]map[uint64]*node, key uint64) (int, *node) {
	low, high := 1, 63
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
	layers [64]map[uint64]*node
	root   *node
	num    uint64
	cache  [64]*node // we'll not need this for the leaf node
}

func (xft *XFastTrie) init() {
	xft.layers = [64]map[uint64]*node{}
	for i := uint64(0); i < 64; i++ {
		xft.layers[i] = make(map[uint64]*node, 50) // we can obviously be more intelligent about this.
	}
	xft.num = 0
	xft.cache = [64]*node{}
	xft.root = newNode(nil, nil)
}

// Exists returns a bool indicating if the provided
// key exists in the trie.  This is typically an
// O(1) operation.
func (xft *XFastTrie) Exists(key uint64) bool {
	// the bottom hashmap of the trie has every entry
	// in it.
	_, ok := xft.layers[63][key]
	return ok
}

// Len returns the number of items in this trie.
func (xft *XFastTrie) Len() uint64 {
	return xft.num
}

func (xft *XFastTrie) insert(entry Entry) {
	key := entry.Key() // so we aren't calling this interface method over and over, fucking Go
	n := xft.layers[63][key]
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

	for i := uint64(0); i < 63; i++ {
		// on 0th, this will be root
		xft.cache[i] = n
		// find out if we need to go left or right
		leftOrRight = (key & positions[i]) >> (63 - i)
		if n.children[leftOrRight] == nil || isLeaf(n.children[leftOrRight]) {
			nn := newNode(n, nil)
			n.children[leftOrRight] = nn
			xft.layers[i][key&masks[i]] = nn // prefix for this layer
		}
		n = n.children[leftOrRight]
	}

	// we are left with next to last possible node
	leftOrRight = key & positions[63] // this will just be 1 or 0
	if n.children[leftOrRight] == nil || isLeaf(n.children[leftOrRight]) {
		leaf := newNode(n, entry)
		n.children[leftOrRight] = leaf
		xft.layers[63][key] = leaf
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
	i := 0
	n := successor.parent
	for n != nil && xft.cache[63-i] != n {
		if isLeaf(n.children[0]) {
			n.children[0] = node
		}
		i++
		if i > 63 {
			break
		}
		n = n.parent
	}
}

func (xft *XFastTrie) walkUpPredecessor(node, predecessor *node) {
	i := 0
	n := predecessor.parent
	for n != nil && xft.cache[63-i] != n {
		if isLeaf(n.children[1]) {
			n.children[1] = node
		}
		i++
		if i > 63 {
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

	n := xft.layers[63][key]
	if n != nil {
		return n
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
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
			log.Printf(`n.children[1].entry: %+v`, n.children[0].entry)
			if n.children[0].entry.Key() >= key {
				return n.children[0]
			}
		}

		return n.children[1].children[0]
	}

	return nil
}

func (xft *XFastTrie) successor(key uint64) *node {
	log.Printf(`SUCCESSOR CALL`)
	if xft.root == nil { // no successor if no nodes
		return nil
	}

	n := xft.layers[63][key]
	if n != nil {
		return n
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
	if n == nil && layer > 1 {
		return nil
	} else if n == nil {
		n = xft.root
	}

	log.Printf(`N0: %+v, N1: %+v`, n.children[0], n.children[1])

	if isLeaf(n.children[0]) {
		if n.children[0].entry.Key() >= key {
			return n.children[0]
		} else if isLeaf(n.children[1]) {
			println(`THIS FUCKING WAY`)
			log.Printf(`n.children[1].entry: %+v`, n.children[1].entry)
			if n.children[1].entry.Key() >= key {
				return n.children[1]
			}
		}

		return n.children[0].children[1]
	}

	return nil

	/*
		println(`HIT THIS`)
		leftOrRight := (key & positions[layer])

		for leftOrRight > 0 {
			layer--
			n = xft.layers[layer][key&masks[layer]] // there should always be n
			leftOrRight = (key & positions[layer])
		}

		n = n.children[1] // now grab the right child

		if n == nil { // it's possible the right child doesn't exist, there may only be one entry
			return nil
		}

		// then grab the cheapest child
		for n.entry == nil {
			if n.children[0] != nil {
				n = n.children[0]
				continue
			}

			n = n.children[1]
		}

		return n*/
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

func New() *XFastTrie {
	xft := &XFastTrie{}
	xft.init()
	return xft
}
