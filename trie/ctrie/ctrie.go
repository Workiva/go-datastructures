/*
Copyright 2015 Workiva, LLC

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

/*
Package ctrie provides an implementation of the Ctrie data structure, which is
a concurrent, lock-free hash trie. This data structure was originally presented
in the paper Concurrent Tries with Efficient Non-Blocking Snapshots:

https://axel22.github.io/resources/docs/ctries-snapshot.pdf

TODO: Add snapshot support.
*/
package ctrie

import (
	"bytes"
	"hash"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/Workiva/go-datastructures/list"
)

const (
	// w controls the number of branches at a node (2^w branches).
	w = 5

	// exp2 is 2^w, which is the hashcode space.
	exp2 = 32
)

// Ctrie is a concurrent, lock-free hash trie. By default, keys are hashed
// using FNV-1a, but the hashing function used can be set with SetHash.
type Ctrie struct {
	root *iNode
	h    hash.Hash32
	hMu  sync.Mutex
}

// iNode is an indirection node. I-nodes remain present in the Ctrie even as
// nodes above and below change. Thread-safety is achieved in part by
// performing CAS operations on the I-node instead of the internal node array.
type iNode struct {
	main *mainNode
}

// mainNode is either a cNode, tNode, or lNode which makes up an I-node.
type mainNode struct {
	cNode *cNode
	tNode *tNode
	lNode *lNode
}

// cNode is an internal main node containing a bitmap and the array with
// references to branch nodes. A branch node is either another I-node or a
// singleton S-node.
type cNode struct {
	bmp   uint32
	array []branch
}

// newMainNode is a recursive constructor which creates a new mainNode. This
// mainNode will consist of cNodes as long as the hashcode chunks of the two
// keys are equal at the given level. If the level exceeds 2^w, an lNode is
// created.
func newMainNode(x *sNode, xhc uint32, y *sNode, yhc uint32, lev uint) *mainNode {
	if lev < exp2 {
		xidx := (xhc >> lev) & 0x1f
		yidx := (yhc >> lev) & 0x1f
		bmp := uint32((1 << xidx) | (1 << yidx))

		if xidx == yidx {
			// Recurse when indexes are equal.
			main := newMainNode(x, xhc, y, yhc, lev+w)
			iNode := &iNode{main}
			return &mainNode{cNode: &cNode{bmp, []branch{iNode}}}
		}
		if xidx < yidx {
			return &mainNode{cNode: &cNode{bmp, []branch{x, y}}}
		}
		return &mainNode{cNode: &cNode{bmp, []branch{y, x}}}
	}
	l := list.Empty.Add(x).Add(y)
	return &mainNode{lNode: &lNode{l}}
}

// inserted returns a copy of this cNode with the new entry at the given
// position.
func (c *cNode) inserted(pos, flag uint32, br branch) *cNode {
	length := uint32(len(c.array))
	bmp := c.bmp
	array := make([]branch, length+1)
	copy(array, c.array)
	array[pos] = br
	for i, x := pos, uint32(0); x < length-pos; i++ {
		array[i+1] = c.array[i]
		x++
	}
	ncn := &cNode{bmp: bmp | flag, array: array}
	return ncn
}

// updated returns a copy of this cNode with the entry at the given index
// updated.
func (c *cNode) updated(pos uint32, br branch) *cNode {
	array := make([]branch, len(c.array))
	copy(array, c.array)
	array[pos] = br
	ncn := &cNode{bmp: c.bmp, array: array}
	return ncn
}

// removed returns a copy of this cNode with the entry at the given index
// removed.
func (c *cNode) removed(pos, flag uint32) *cNode {
	length := uint32(len(c.array))
	bmp := c.bmp
	array := make([]branch, length-1)
	for i := uint32(0); i < pos; i++ {
		array[i] = c.array[i]
	}
	for i, x := pos, uint32(0); x < length-pos-1; i++ {
		array[i] = c.array[i+1]
		x++
	}
	ncn := &cNode{bmp: bmp ^ flag, array: array}
	return ncn
}

// tNode is tomb node which is a special node used to ensure proper ordering
// during removals.
type tNode struct {
	*sNode
}

// untombed returns the S-node contained by the T-node.
func (t *tNode) untombed() *sNode {
	return &sNode{&entry{key: t.key, hash: t.hash, value: t.value}}
}

// lNode is a list node which is a leaf node used to handle hashcode
// collisions by keeping such keys in a persistent list.
type lNode struct {
	list.PersistentList
}

// entry returns the first S-node contained in the L-node.
func (l *lNode) entry() *sNode {
	head, _ := l.Head()
	return head.(*sNode)
}

// lookup returns the value at the given entry in the L-node or returns false
// if it's not contained.
func (l *lNode) lookup(e *entry) (interface{}, bool) {
	found, ok := l.Find(func(sn interface{}) bool {
		return bytes.Equal(e.key, sn.(*sNode).key)
	})
	if !ok {
		return nil, false
	}
	return found.(*sNode).value, true
}

// inserted creates a new L-node with the added entry.
func (l *lNode) inserted(entry *entry) *lNode {
	return &lNode{l.Add(&sNode{entry})}
}

// removed creates a new L-node with the entry removed.
func (l *lNode) removed(e *entry) *lNode {
	idx := l.FindIndex(func(sn interface{}) bool {
		return bytes.Equal(e.key, sn.(*sNode).key)
	})
	if idx < 0 {
		return l
	}
	nl, _ := l.Remove(uint(idx))
	return &lNode{nl}
}

// length returns the L-node list length.
func (l *lNode) length() uint {
	return l.Length()
}

// branch is either an iNode or sNode.
type branch interface{}

// entry contains a Ctrie entry, which is also a technique used to cache the
// hashcode of the key.
type entry struct {
	key   []byte
	hash  uint32
	value interface{}
}

// sNode is a singleton node which contains a single key and value.
type sNode struct {
	*entry
}

// New creates an empty Ctrie, defaulting to FNV-1a for key hashing. Use
// SetHash to change the hash function.
func New() *Ctrie {
	root := &iNode{main: &mainNode{cNode: &cNode{}}}
	return &Ctrie{root: root, h: fnv.New32a()}
}

// SetHash sets the hash function used by the Ctrie. Existing entries are not
// rehashed when this is set, so this should be called on a newly created
// Ctrie.
func (c *Ctrie) SetHash(hash hash.Hash32) {
	c.hMu.Lock()
	c.h = hash
	c.hMu.Unlock()
}

// Insert adds the key-value pair to the Ctrie, replacing the existing value if
// the key already exists.
func (c *Ctrie) Insert(key []byte, value interface{}) {
	c.insert(&entry{
		key:   key,
		hash:  c.hash(key),
		value: value,
	})
}

// Lookup returns the value for the associated key or returns false if the key
// doesn't exist.
func (c *Ctrie) Lookup(key []byte) (interface{}, bool) {
	return c.lookup(&entry{key: key, hash: c.hash(key)})
}

// Remove deletes the value for the associated key, returning true if it was
// removed or false if the entry doesn't exist.
func (c *Ctrie) Remove(key []byte) (interface{}, bool) {
	return c.remove(&entry{key: key, hash: c.hash(key)})
}

func (c *Ctrie) insert(entry *entry) {
	rootPtr := (*unsafe.Pointer)(unsafe.Pointer(&c.root))
	root := (*iNode)(atomic.LoadPointer(rootPtr))
	if !iinsert(root, entry, 0, nil) {
		c.insert(entry)
	}
}

func (c *Ctrie) lookup(entry *entry) (interface{}, bool) {
	rootPtr := (*unsafe.Pointer)(unsafe.Pointer(&c.root))
	root := (*iNode)(atomic.LoadPointer(rootPtr))
	result, exists, ok := ilookup(root, entry, 0, nil)
	for !ok {
		return c.lookup(entry)
	}
	return result, exists
}

func (c *Ctrie) remove(entry *entry) (interface{}, bool) {
	rootPtr := (*unsafe.Pointer)(unsafe.Pointer(&c.root))
	root := (*iNode)(atomic.LoadPointer(rootPtr))
	result, exists, ok := iremove(root, entry, 0, nil)
	for !ok {
		return c.remove(entry)
	}
	return result, exists
}

func (c *Ctrie) hash(k []byte) uint32 {
	c.hMu.Lock()
	c.h.Write(k)
	hash := c.h.Sum32()
	c.h.Reset()
	c.hMu.Unlock()
	return hash
}

func iinsert(i *iNode, entry *entry, lev uint, parent *iNode) bool {
	mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&i.main))
	main := (*mainNode)(atomic.LoadPointer(mainPtr))
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the relevant bit is not in the bitmap, then a copy of the
			// cNode with the new entry is created. The linearization point is
			// a successful CAS.
			ncn := &mainNode{cNode: cn.inserted(pos, flag, &sNode{entry})}
			return atomic.CompareAndSwapPointer(
				mainPtr, unsafe.Pointer(main), unsafe.Pointer(ncn))
		}
		// If the relevant bit is present in the bitmap, then its corresponding
		// branch is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, then iinsert is called recursively.
			return iinsert(branch.(*iNode), entry, lev+w, i)
		case *sNode:
			sn := branch.(*sNode)
			if !bytes.Equal(sn.key, entry.key) {
				// If the branch is an S-node and its key is not equal to the
				// key being inserted, then the Ctrie has to be extended with
				// an additional level. The C-node is replaced with its updated
				// version, created using the updated function that adds a new
				// I-node at the respective position. The new Inode has its
				// main node pointing to a C-node with both keys. The
				// linearization point is a successful CAS.
				nsn := &sNode{entry}
				nin := &iNode{newMainNode(sn, sn.hash, nsn, nsn.hash, lev+w)}
				ncn := &mainNode{cNode: cn.updated(pos, nin)}
				return atomic.CompareAndSwapPointer(
					mainPtr, unsafe.Pointer(main), unsafe.Pointer(ncn))
			}
			// If the key in the S-node is equal to the key being inserted,
			// then the C-node is replaced with its updated version with a new
			// S-node. The linearization point is a successful CAS.
			ncn := &mainNode{cNode: cn.updated(pos, &sNode{entry})}
			return atomic.CompareAndSwapPointer(
				mainPtr, unsafe.Pointer(main), unsafe.Pointer(ncn))
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w)
		return false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.inserted(entry)}
		return atomic.CompareAndSwapPointer(
			mainPtr, unsafe.Pointer(main), unsafe.Pointer(nln))
	default:
		panic("Ctrie is in an invalid state")
	}
}

func ilookup(i *iNode, entry *entry, lev uint, parent *iNode) (interface{}, bool, bool) {
	mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&i.main))
	// Linearization point.
	main := (*mainNode)(atomic.LoadPointer(mainPtr))
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the bitmap does not contain the relevant bit, a key with the
			// required hashcode prefix is not present in the trie.
			return nil, false, true
		}
		// Otherwise, the relevant branch at index pos is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, the ilookup procedure is called
			// recursively at the next level.
			return ilookup(branch.(*iNode), entry, lev+w, i)
		case *sNode:
			// If the branch is an S-node, then the key within the S-node is
			// compared with the key being searched – these two keys have the
			// same hashcode prefixes, but they need not be equal. If they are
			// equal, the corresponding value from the S-node is
			// returned and a NOTFOUND value otherwise.
			sn := branch.(*sNode)
			if bytes.Equal(sn.key, entry.key) {
				return sn.value, true, true
			}
			return nil, false, true
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w)
		return nil, false, false
	case main.lNode != nil:
		// Hash collisions are handled using L-nodes, which are essentially
		// persistent linked lists.
		val, ok := main.lNode.lookup(entry)
		return val, ok, true
	default:
		panic("Ctrie is in an invalid state")
	}
}

func iremove(i *iNode, entry *entry, lev uint, parent *iNode) (interface{}, bool, bool) {
	mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&i.main))
	// Linearization point.
	main := (*mainNode)(atomic.LoadPointer(mainPtr))
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the bitmap does not contain the relevant bit, a key with the
			// required hashcode prefix is not present in the trie.
			return nil, false, true
		}
		// Otherwise, the relevant branch at index pos is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, the iremove procedure is called
			// recursively at the next level.
			return iremove(branch.(*iNode), entry, lev+w, i)
		case *sNode:
			// If the branch is an S-node, its key is compared against the key
			// being removed.
			sn := branch.(*sNode)
			if !bytes.Equal(sn.key, entry.key) {
				// If the keys are not equal, the NOTFOUND value is returned.
				return nil, false, true
			}
			//  If the keys are equal, a copy of the current node without the
			//  S-node is created. The contraction of the copy is then created
			//  using the toContracted procedure. A successful CAS will
			//  substitute the old C-node with the copied C-node, thus removing
			//  the S-node with the given key from the trie – this is the
			//  linearization point
			ncn := cn.removed(pos, flag)
			cntr := toContracted(ncn, lev)
			if atomic.CompareAndSwapPointer(
				mainPtr, unsafe.Pointer(main), unsafe.Pointer(cntr)) {
				if parent != nil {
					main = (*mainNode)(atomic.LoadPointer(mainPtr))
					if main.tNode != nil {
						cleanParent(parent, i, entry.hash, lev-w)
					}
				}
				return sn.value, true, true
			}
			return nil, false, false
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w)
		return nil, false, false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.removed(entry)}
		if nln.lNode.length() == 1 {
			nln = entomb(nln.lNode.entry())
		}
		if atomic.CompareAndSwapPointer(
			mainPtr, unsafe.Pointer(main), unsafe.Pointer(nln)) {
			val, ok := main.lNode.lookup(entry)
			return val, ok, true
		}
		return nil, false, true
	default:
		panic("Ctrie is in an invalid state")
	}
}

// toContracted ensures that every I-node except the root points to a C-node
// with at least one branch. If a given C-Node has only a single S-node below
// it and is not at the root level, a T-node which wraps the S-node is
// returned.
func toContracted(cn *cNode, lev uint) *mainNode {
	if lev > 0 && len(cn.array) == 1 {
		branch := cn.array[0]
		switch branch.(type) {
		case *sNode:
			return entomb(branch.(*sNode))
		default:
			return &mainNode{cNode: cn}
		}
	}
	return &mainNode{cNode: cn}
}

// toCompressed compacts the C-node as a performance optimization.
func toCompressed(cn *cNode, lev uint) *mainNode {
	bmp := cn.bmp
	i := 0
	arr := cn.array
	tmpArray := make([]branch, len(arr))
	for i < len(arr) {
		sub := arr[i]
		switch sub.(type) {
		case *iNode:
			inode := sub.(*iNode)
			mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&inode.main))
			main := (*mainNode)(atomic.LoadPointer(mainPtr))
			tmpArray[i] = resurrect(inode, main)
		case *sNode:
			tmpArray[i] = sub
		default:
			panic("Ctrie is in an invalid state")
		}
		i++
	}

	return toContracted(&cNode{bmp: bmp, array: tmpArray}, lev)
}

func entomb(m *sNode) *mainNode {
	return &mainNode{tNode: &tNode{m}}
}

func resurrect(iNode *iNode, main *mainNode) branch {
	if main.tNode != nil {
		return main.tNode.untombed()
	}
	return iNode
}

func clean(i *iNode, lev uint) bool {
	mainPtr := (*unsafe.Pointer)(unsafe.Pointer(&i.main))
	main := (*mainNode)(atomic.LoadPointer(mainPtr))
	if main.cNode != nil {
		return atomic.CompareAndSwapPointer(mainPtr, unsafe.Pointer(main),
			unsafe.Pointer(toCompressed(main.cNode, lev)))
	}
	return true
}

func cleanParent(p, i *iNode, hc uint32, lev uint) {
	var (
		mainPtr  = (*unsafe.Pointer)(unsafe.Pointer(&i.main))
		main     = (*mainNode)(atomic.LoadPointer(mainPtr))
		pMainPtr = (*unsafe.Pointer)(unsafe.Pointer(&p.main))
		pMain    = (*mainNode)(atomic.LoadPointer(pMainPtr))
	)
	if pMain.cNode != nil {
		flag, pos := flagPos(hc, lev, pMain.cNode.bmp)
		if pMain.cNode.bmp&flag != 0 {
			sub := pMain.cNode.array[pos]
			if sub == i && main.tNode != nil {
				ncn := pMain.cNode.updated(pos, resurrect(i, main))
				if !atomic.CompareAndSwapPointer(pMainPtr, unsafe.Pointer(pMain),
					unsafe.Pointer(toContracted(ncn, lev))) {
					cleanParent(p, i, hc, lev)
				}
			}
		}
	}
}

func flagPos(hashcode uint32, lev uint, bmp uint32) (uint32, uint32) {
	idx := (hashcode >> lev) & 0x1f
	flag := uint32(1) << uint32(idx)
	mask := uint32(flag - 1)
	pos := bitCount(bmp & mask)
	return flag, pos
}

func bitCount(x uint32) uint32 {
	x = ((x >> 1) & 0x55555555) + (x & 0x55555555)
	x = ((x >> 2) & 0x33333333) + (x & 0x33333333)
	x = ((x >> 4) & 0x0f0f0f0f) + (x & 0x0f0f0f0f)
	x = ((x >> 8) & 0x00ff00ff) + (x & 0x00ff00ff)
	return ((x >> 16) & 0x0000ffff) + (x & 0x0000ffff)
}
