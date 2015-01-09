package xfast

import (
	"log"
)

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
}

func newNode(entry Entry) *node {
	return &node{
		children: [2]*node{},
		entry:    entry,
	}
}

type leaf struct {
	entry Entry
}

func newLeaf(entry Entry) *leaf {
	return &leaf{
		entry: entry,
	}
}

func binarySearchHashMaps(layers [64]map[uint64]*node, key uint64) (int, *node) {
	low, high := 1, 63
	var mid int
	var node *node
	for low <= high {
		mid = (low + high) / 2
		n, ok := layers[mid][key*masks[mid]]
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
	layers [64]map[uint64]*node
	root   *node
	num    uint64
	cache  [63]*node // we'll not need this for the leaf node
}

func (xft *XFastTrie) init() {
	xft.layers = [64]map[uint64]*node{}
	for i := uint64(0); i < 64; i++ {
		xft.layers[i] = make(map[uint64]*node, 50) // we can obviously be more intelligent about this.
	}
	xft.num = 0
	xft.cache = [63]*node{}
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
	if xft.root == nil {
		xft.root = newNode(nil)
	} else {
		n := xft.layers[63][key]
		if n != nil {
			n.entry = entry
			return
		}
	}

	node := xft.root
	var leftOrRight uint64
	key := entry.Key() // so we aren't calling this interface method over and over, fucking Go
	for i := uint64(0); i < 63; i++ {
		xft.cache[i] = node
		leftOrRight = (key & positions[i]) >> (63 - i) // find out if we need to go left or right
		node = node.children[leftOrRight]
		if node == nil {
			node = newNode(nil)
			xft.cache[i].children[leftOrRight] = node
			xft.layers[i+1][key&masks[i]] = node // prefix for this layer
		}
	}

	// we are left with next to last possible node
	leftOrRight = key & positions[63] // this will just be 1 or 0
	if node.children[leftOrRight] == nil {
		leaf := newNode(entry)
		node.children[leftOrRight] = leaf
		xft.layers[63][key] = leaf
		xft.num++
	} else {
		log.Printf(`WE HAVE A PROBLEM HERE WITH KEY: %+v`, key)
	}
}

func (xft *XFastTrie) Insert(entries ...Entry) {
	for _, e := range entries {
		xft.insert(e)
	}
}

func (xft *XFastTrie) successor(key uint64) *node {
	n := xft.layers[63][key]
	if n != nil {
		return n.children[1]
	}

	layer, n := binarySearchHashMaps(xft.layers, key)
	leftOrRight := (key & positions[layer])
}

func (xft *XFastTrie) Successor(key uint64) uint64 {

}

func New() *XFastTrie {
	xft := &XFastTrie{}
	xft.init()
	return xft
}
