package plus

import "sort"

type node interface {
	insert(tree *btree, key Key) bool
}

type nodes []node

type inode struct {
	keys  keys
	nodes nodes
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

	if uint64(len(lnode.keys)) == tree.nodeSize { // all the magic happens here

	}

	return true
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
