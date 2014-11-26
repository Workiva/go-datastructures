package plus

import (
	"log"
	"sort"
)

func init() {
	log.Printf(`I HATE THIS`)
}

func split(tree *btree, parent, child node) node {
	if !child.needsSplit(tree.nodeSize) {
		return parent
	}

	switch child.(type) {
	case *lnode: // we need to split a leaf

	}
	return nil
}

type node interface {
	insert(tree *btree, key Key) bool
	needsSplit(nodeSize uint64) bool
	// key is the median key while left and right nodes
	// represent the left and right nodes respectively
	split() (Key, node, node)
}

type nodes []node

type inode struct {
	keys  keys
	nodes nodes
}

func newInternalNode(size uint64) *inode {
	return &inode{}
}

type lnode struct {
	// points to the left leaf node is there is one
	pointer *lnode
	keys    keys
}

func (lnode *lnode) insert(tree *btree, key Key) bool {
	i := keySearch(lnode.keys, key)
	var inserted bool
	if i == len(lnode.keys) { // simple append will do
		lnode.keys = append(lnode.keys, newPayload(key))
		inserted = true
	} else {
		if lnode.keys[i].Compare(key) == 0 {
			inserted = lnode.keys[i].(*payload).insert(key)
		} else {
			lnode.keys.insertAt(i, newPayload(key))
			inserted = true
		}
	}

	if !inserted {
		return false
	}

	return true
}

func (node *lnode) split() (Key, node, node) {
	if len(node.keys) < 2 {
		return nil, nil, nil
	}
	i := len(node.keys) / 2
	key := node.keys[i].(*payload).key()
	otherKeys := make(keys, i)
	ourKeys := make(keys, len(node.keys)-i)
	// we perform these copies so these slices don't all end up
	// pointing to the same underlying array which may make
	// for some very difficult to debug situations later.
	copy(otherKeys, node.keys[:i])
	copy(ourKeys, node.keys[i:])

	// this should release the original array for GC
	node.keys = ourKeys
	otherNode := &lnode{
		keys:    otherKeys,
		pointer: node,
	}
	return key, otherNode, node
}

func (lnode *lnode) needsSplit(nodeSize uint64) bool {
	return uint64(len(lnode.keys)) >= nodeSize
}

func newLeafNode(size uint64) *lnode {
	return &lnode{
		keys: make(keys, 0, size),
	}
}

type payload struct {
	keys sortedByIDKeys
}

func (payload *payload) insert(key Key) bool {
	return payload.keys.insert(key)
}

func (payload *payload) ID() uint64 {
	return 0
}

func (payload *payload) Compare(key Key) int {
	if len(payload.keys) == 0 {
		panic(`WE HAVE A PAYLOAD WITH NO KEYS`)
	}

	return payload.keys[0].Compare(key)
}

func (payload *payload) key() Key {
	return payload.keys[0]
}

func newPayload(key Key) *payload {
	p := &payload{
		keys: make(sortedByIDKeys, 0, 5),
	}
	p.keys = append(p.keys, key)
	return p
}

type keys []Key

func (keys *keys) insertAt(i int, key Key) {
	if i == len(*keys) {
		*keys = append(*keys, key)
		return
	}

	*keys = append(*keys, nil)
	copy((*keys)[i+1:], (*keys)[i:])
	(*keys)[i] = key
}

type sortedByIDKeys keys

func (sorted sortedByIDKeys) search(id uint64) int {
	return sort.Search(len(sorted), func(i int) bool {
		return sorted[i].ID() >= id
	})
}

func (sorted *sortedByIDKeys) insert(key Key) bool {
	i := sorted.search(key.ID())
	if i == len(*sorted) {
		*sorted = append(*sorted, key)
		return true
	}

	if (*sorted)[i].ID() == key.ID() { // we don't allow duplicates
		return false
	}

	*sorted = append(*sorted, nil)
	copy((*sorted)[i+1:], (*sorted)[i:])
	(*sorted)[i] = key
	return true
}
