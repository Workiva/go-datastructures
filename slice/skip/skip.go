package skip

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

const p = .5 // the p level defines the probability that a node
// with a value at level i also has a value at i+1.  This number
// is also important in determining max level.  Max level will
// be defined as L(N) where L = log base (1/p) of n where n
// is the number of items in the list and N is the number of possible
// items in the universe.  If p = .5 then maxlevel = 32 is appropriate
// for uint32.

// generator will be the common generator to create random numbers. It
// is seeded with unix nanosecond when this line is executed at runtime,
// and only executed once ensuring all random numbers come from the same
// randomly seeded generator.
var generator = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateLevel(maxLevel uint8) uint8 {
	var level uint8
	for level = uint8(1); level < maxLevel-1; level++ {
		if generator.Float64() >= p {
			return level
		}
	}

	return level
}

type SkipList struct {
	maxLevel, level uint8
	head, tail      *node
	num             uint64
	// a list of nodes that can be reused, should reduce
	// the number of allocations in the insert/delete case.
	cache nodes
}

// init will initialize this skiplist.  The parameter is expected
// to be of some uint type which will set this skiplist's maximum
// level.
func (sl *SkipList) init(ifc interface{}) {
	switch ifc.(type) {
	case uint8:
		sl.maxLevel = 8
	case uint16:
		sl.maxLevel = 16
	case uint32:
		sl.maxLevel = 32
	case uint64, uint:
		sl.maxLevel = 64
	}
	sl.cache = make(nodes, sl.maxLevel)
	sl.head = newNode(nil, sl.maxLevel)
}

func (sl *SkipList) search(key uint64, update []*node) *node {
	if sl.num == 0 { // nothing in the list
		return nil
	}

	var offset uint8
	n := sl.head
	for i := uint8(0); i <= sl.level; i++ {
		offset = sl.level - i
		for n.forward[offset] != nil && n.forward[offset].key() < key {
			n = n.forward[offset]
		}

		if update != nil {
			update[offset] = n
		}
	}

	return n.forward[0]
}

func (sl *SkipList) Get(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))

	var n *node
	for _, key := range keys {
		n = sl.search(key, nil)
		if n == nil {
			entries = append(entries, nil)
		} else {
			entries = append(entries, n.entry)
		}
	}

	return entries
}

func (sl *SkipList) insert(entry Entry) Entry {
	sl.cache.reset()
	n := sl.search(entry.Key(), sl.cache)
	if n != nil && n.key() == entry.Key() { // a simple update in this case
		oldEntry := n.entry
		n.entry = entry
		return oldEntry
	}
	sl.num++

	nodeLevel := generateLevel(sl.maxLevel)
	if nodeLevel > sl.level {
		for i := sl.level; i <= nodeLevel; i++ {
			sl.cache[i] = sl.head
		}
		sl.level = nodeLevel
	}

	nn := newNode(entry, sl.maxLevel)
	for i := uint8(0); i <= nodeLevel; i++ {
		nn.forward[i] = sl.cache[i].forward[i]
		sl.cache[i].forward[i] = nn
	}

	return nil
}

// Insert will insert the provided entries into the list.  Returned
// is a list of entries that were overwritten.  This is expected to
// be an O(log n) operation.
func (sl *SkipList) Insert(entries ...Entry) Entries {
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, sl.insert(e))
	}

	return overwritten
}

func (sl *SkipList) delete(key uint64) Entry {
	sl.cache.reset()
	n := sl.search(key, sl.cache)

	if n == nil || n.entry.Key() != key {
		return nil
	}

	sl.num--

	for i := uint8(0); i <= sl.level; i++ {
		if sl.cache[i].forward[i] != n {
			break
		}

		sl.cache[i].forward[i] = n.forward[i]
	}

	for sl.level > 0 && sl.head.forward[sl.level] == nil {
		sl.level = sl.level - 1
	}

	return n.entry
}

func (sl *SkipList) Delete(keys ...uint64) Entries {
	deleted := make(Entries, 0, len(keys))

	for _, key := range keys {
		deleted = append(deleted, sl.delete(key))
	}

	return deleted
}

// Len returns the number of items in this skiplist.
func (sl *SkipList) Len() uint64 {
	return sl.num
}

func New(ifc interface{}) *SkipList {
	sl := &SkipList{}
	sl.init(ifc)
	return sl
}
