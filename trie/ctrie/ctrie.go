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
*/
package ctrie

import (
	"bytes"
	"hash"
	"hash/fnv"
	"sync/atomic"
	"unsafe"

	"github.com/Workiva/go-datastructures/list"
	"github.com/Workiva/go-datastructures/queue"
)

const (
	// w controls the number of branches at a node (2^w branches).
	w = 5

	// exp2 is 2^w, which is the hashcode space.
	exp2 = 32

	// hasherPoolSize is the number of hashers to buffer.
	hasherPoolSize = 16
)

// HashFactory returns a new Hash32 used to hash keys.
type HashFactory func() hash.Hash32

func defaultHashFactory() hash.Hash32 {
	return fnv.New32a()
}

// Ctrie is a concurrent, lock-free hash trie. By default, keys are hashed
// using FNV-1a unless a HashFactory is provided to New.
type Ctrie struct {
	root        *iNode
	readOnly    bool
	hashFactory HashFactory
	hasherPool  *queue.RingBuffer
}

// generation demarcates Ctrie snapshots. We use a heap-allocated reference
// instead of an integer to avoid integer overflows.
type generation struct{}

// iNode is an indirection node. I-nodes remain present in the Ctrie even as
// nodes above and below change. Thread-safety is achieved in part by
// performing CAS operations on the I-node instead of the internal node array.
type iNode struct {
	main *mainNode
	gen  *generation
}

// copyToGen returns a copy of this I-node copied to the given generation.
func (i *iNode) copyToGen(gen *generation, ctrie *Ctrie) *iNode {
	nin := &iNode{gen: gen}
	main := gcasRead(i, ctrie)
	atomic.StorePointer(
		(*unsafe.Pointer)(unsafe.Pointer(&nin.main)), unsafe.Pointer(main))
	return nin
}

// mainNode is either a cNode, tNode, lNode, or failed node which makes up an
// I-node.
type mainNode struct {
	cNode  *cNode
	tNode  *tNode
	lNode  *lNode
	failed *mainNode

	// prev is set as a failed main node when we attempt to CAS and the
	// I-node's generation does not match the root generation. This signals
	// that the GCAS failed and the I-node's main node must be set back to the
	// previous value.
	prev *mainNode
}

// cNode is an internal main node containing a bitmap and the array with
// references to branch nodes. A branch node is either another I-node or a
// singleton S-node.
type cNode struct {
	bmp   uint32
	array []branch
	gen   *generation
}

// newMainNode is a recursive constructor which creates a new mainNode. This
// mainNode will consist of cNodes as long as the hashcode chunks of the two
// keys are equal at the given level. If the level exceeds 2^w, an lNode is
// created.
func newMainNode(x *sNode, xhc uint32, y *sNode, yhc uint32, lev uint, gen *generation) *mainNode {
	if lev < exp2 {
		xidx := (xhc >> lev) & 0x1f
		yidx := (yhc >> lev) & 0x1f
		bmp := uint32((1 << xidx) | (1 << yidx))

		if xidx == yidx {
			// Recurse when indexes are equal.
			main := newMainNode(x, xhc, y, yhc, lev+w, gen)
			iNode := &iNode{main: main, gen: gen}
			return &mainNode{cNode: &cNode{bmp, []branch{iNode}, gen}}
		}
		if xidx < yidx {
			return &mainNode{cNode: &cNode{bmp, []branch{x, y}, gen}}
		}
		return &mainNode{cNode: &cNode{bmp, []branch{y, x}, gen}}
	}
	l := list.Empty.Add(x).Add(y)
	return &mainNode{lNode: &lNode{l}}
}

// inserted returns a copy of this cNode with the new entry at the given
// position.
func (c *cNode) inserted(pos, flag uint32, br branch, gen *generation) *cNode {
	length := uint32(len(c.array))
	bmp := c.bmp
	array := make([]branch, length+1)
	copy(array, c.array)
	array[pos] = br
	for i, x := pos, uint32(0); x < length-pos; i++ {
		array[i+1] = c.array[i]
		x++
	}
	ncn := &cNode{bmp: bmp | flag, array: array, gen: gen}
	return ncn
}

// updated returns a copy of this cNode with the entry at the given index
// updated.
func (c *cNode) updated(pos uint32, br branch, gen *generation) *cNode {
	array := make([]branch, len(c.array))
	copy(array, c.array)
	array[pos] = br
	ncn := &cNode{bmp: c.bmp, array: array, gen: gen}
	return ncn
}

// removed returns a copy of this cNode with the entry at the given index
// removed.
func (c *cNode) removed(pos, flag uint32, gen *generation) *cNode {
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
	ncn := &cNode{bmp: bmp ^ flag, array: array, gen: gen}
	return ncn
}

// renewed returns a copy of this cNode with the I-nodes below it copied to the
// given generation.
func (c *cNode) renewed(gen *generation, ctrie *Ctrie) *cNode {
	array := make([]branch, len(c.array))
	for i, br := range c.array {
		switch br.(type) {
		case *iNode:
			array[i] = br.(*iNode).copyToGen(gen, ctrie)
		default:
			array[i] = br
		}
	}
	return &cNode{bmp: c.bmp, array: array, gen: gen}
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

// New creates an empty Ctrie which uses the provided HashFactory for key
// hashing. If nil is passed in, it will default to FNV-1a hashing.
func New(hashFactory HashFactory) *Ctrie {
	if hashFactory == nil {
		hashFactory = defaultHashFactory
	}
	root := &iNode{main: &mainNode{cNode: &cNode{}}}
	return newCtrie(root, hashFactory, false)
}

func newCtrie(root *iNode, hashFactory HashFactory, readOnly bool) *Ctrie {
	hasherPool := queue.NewRingBuffer(hasherPoolSize)
	for i := 0; i < hasherPoolSize; i++ {
		hasherPool.Put(hashFactory())
	}
	return &Ctrie{
		root:        root,
		hashFactory: hashFactory,
		hasherPool:  hasherPool,
		readOnly:    readOnly,
	}
}

// Insert adds the key-value pair to the Ctrie, replacing the existing value if
// the key already exists.
func (c *Ctrie) Insert(key []byte, value interface{}) {
	c.assertReadWrite()
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
	c.assertReadWrite()
	return c.remove(&entry{key: key, hash: c.hash(key)})
}

// Snapshot returns a stable, point-in-time snapshot of the Ctrie.
func (c *Ctrie) Snapshot() *Ctrie {
	for {
		root := c.readRoot()
		main := gcasRead(root, c)
		if c.rdcssRoot(root, main, root.copyToGen(&generation{}, c)) {
			return newCtrie(root.copyToGen(&generation{}, c), c.hashFactory, c.readOnly)
		}
	}
}

// ReadOnlySnapshot returns a stable, point-in-time snapshot of the Ctrie which
// is read-only. Write operations on a read-only snapshot will panic.
func (c *Ctrie) ReadOnlySnapshot() *Ctrie {
	if c.readOnly {
		return c
	}
	for {
		root := c.readRoot()
		main := gcasRead(root, c)
		if c.rdcssRoot(root, main, root.copyToGen(&generation{}, c)) {
			return newCtrie(root, c.hashFactory, true)
		}
	}
}

func (c *Ctrie) assertReadWrite() {
	if c.readOnly {
		panic("Cannot modify read-only snapshot")
	}
}

func (c *Ctrie) insert(entry *entry) {
	root := c.readRoot()
	if !c.iinsert(root, entry, 0, nil, root.gen) {
		c.insert(entry)
	}
}

func (c *Ctrie) lookup(entry *entry) (interface{}, bool) {
	root := c.readRoot()
	result, exists, ok := c.ilookup(root, entry, 0, nil, root.gen)
	for !ok {
		return c.lookup(entry)
	}
	return result, exists
}

func (c *Ctrie) remove(entry *entry) (interface{}, bool) {
	root := c.readRoot()
	result, exists, ok := c.iremove(root, entry, 0, nil, root.gen)
	for !ok {
		return c.remove(entry)
	}
	return result, exists
}

func (c *Ctrie) hash(k []byte) uint32 {
	hasher, _ := c.hasherPool.Get()
	h := hasher.(hash.Hash32)
	h.Write(k)
	hash := h.Sum32()
	h.Reset()
	c.hasherPool.Put(h)
	return hash
}

// iinsert attempts to insert the entry into the Ctrie. If false is returned,
// the operation should be retried.
func (c *Ctrie) iinsert(i *iNode, entry *entry, lev uint, parent *iNode, startGen *generation) bool {
	// Linearization point.
	main := gcasRead(i, c)
	switch {
	case main.cNode != nil:
		cn := main.cNode
		flag, pos := flagPos(entry.hash, lev, cn.bmp)
		if cn.bmp&flag == 0 {
			// If the relevant bit is not in the bitmap, then a copy of the
			// cNode with the new entry is created. The linearization point is
			// a successful CAS.
			rn := cn
			if cn.gen != i.gen {
				rn = cn.renewed(i.gen, c)
			}
			ncn := &mainNode{cNode: rn.inserted(pos, flag, &sNode{entry}, i.gen)}
			return gcas(i, main, ncn, c)
		}
		// If the relevant bit is present in the bitmap, then its corresponding
		// branch is read from the array.
		branch := cn.array[pos]
		switch branch.(type) {
		case *iNode:
			// If the branch is an I-node, then iinsert is called recursively.
			in := branch.(*iNode)
			if startGen == in.gen {
				return c.iinsert(in, entry, lev+w, i, i.gen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.iinsert(i, entry, lev, parent, startGen)
			}
			return false
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
				rn := cn
				if cn.gen != i.gen {
					rn = cn.renewed(i.gen, c)
				}
				nsn := &sNode{entry}
				nin := &iNode{main: newMainNode(sn, sn.hash, nsn, nsn.hash, lev+w, i.gen), gen: i.gen}
				ncn := &mainNode{cNode: rn.updated(pos, nin, i.gen)}
				return gcas(i, main, ncn, c)
			}
			// If the key in the S-node is equal to the key being inserted,
			// then the C-node is replaced with its updated version with a new
			// S-node. The linearization point is a successful CAS.
			ncn := &mainNode{cNode: cn.updated(pos, &sNode{entry}, i.gen)}
			return gcas(i, main, ncn, c)
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w, c)
		return false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.inserted(entry)}
		return gcas(i, main, nln, c)
	default:
		panic("Ctrie is in an invalid state")
	}
}

// ilookup attempts to fetch the entry from the Ctrie. The first two return
// values are the entry value and whether or not the entry was contained in the
// Ctrie. The last bool indicates if the operation succeeded. False means it
// should be retried.
func (c *Ctrie) ilookup(i *iNode, entry *entry, lev uint, parent *iNode, startGen *generation) (interface{}, bool, bool) {
	// Linearization point.
	main := gcasRead(i, c)
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
			in := branch.(*iNode)
			if c.readOnly || startGen == in.gen {
				return c.ilookup(in, entry, lev+w, i, startGen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.ilookup(i, entry, lev, parent, startGen)
			}
			return nil, false, false
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
		return cleanReadOnly(main.tNode, lev, parent, c, entry)
	case main.lNode != nil:
		// Hash collisions are handled using L-nodes, which are essentially
		// persistent linked lists.
		val, ok := main.lNode.lookup(entry)
		return val, ok, true
	default:
		panic("Ctrie is in an invalid state")
	}
}

// iremove attempts to remove the entry from the Ctrie. The first two return
// values are the entry value and whether or not the entry was contained in the
// Ctrie. The last bool indicates if the operation succeeded. False means it
// should be retried.
func (c *Ctrie) iremove(i *iNode, entry *entry, lev uint, parent *iNode, startGen *generation) (interface{}, bool, bool) {
	// Linearization point.
	main := gcasRead(i, c)
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
			in := branch.(*iNode)
			if startGen == in.gen {
				return c.iremove(in, entry, lev+w, i, startGen)
			}
			if gcas(i, main, &mainNode{cNode: cn.renewed(startGen, c)}, c) {
				return c.iremove(i, entry, lev, parent, startGen)
			}
			return nil, false, false
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
			ncn := cn.removed(pos, flag, i.gen)
			cntr := toContracted(ncn, lev)
			if gcas(i, main, cntr, c) {
				if parent != nil {
					main = gcasRead(i, c)
					if main.tNode != nil {
						cleanParent(parent, i, entry.hash, lev-w, c, startGen)
					}
				}
				return sn.value, true, true
			}
			return nil, false, false
		default:
			panic("Ctrie is in an invalid state")
		}
	case main.tNode != nil:
		clean(parent, lev-w, c)
		return nil, false, false
	case main.lNode != nil:
		nln := &mainNode{lNode: main.lNode.removed(entry)}
		if nln.lNode.length() == 1 {
			nln = entomb(nln.lNode.entry())
		}
		if gcas(i, main, nln, c) {
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
	tmpArray := make([]branch, len(cn.array))
	for i, sub := range cn.array {
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
	}

	return toContracted(&cNode{bmp: cn.bmp, array: tmpArray}, lev)
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

func clean(i *iNode, lev uint, ctrie *Ctrie) bool {
	main := gcasRead(i, ctrie)
	if main.cNode != nil {
		return gcas(i, main, toCompressed(main.cNode, lev), ctrie)
	}
	return true
}

func cleanReadOnly(tn *tNode, lev uint, p *iNode, ctrie *Ctrie, entry *entry) (interface{}, bool, bool) {
	if !ctrie.readOnly {
		clean(p, lev-5, ctrie)
		return nil, false, false
	}
	if tn.hash == entry.hash && bytes.Equal(tn.key, entry.key) {
		return tn.value, true, true
	}
	return nil, false, true
}

func cleanParent(p, i *iNode, hc uint32, lev uint, ctrie *Ctrie, startGen *generation) {
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
				ncn := pMain.cNode.updated(pos, resurrect(i, main), i.gen)
				if !gcas(p, pMain, toContracted(ncn, lev), ctrie) && ctrie.readRoot().gen == startGen {
					cleanParent(p, i, hc, lev, ctrie, startGen)
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

// gcas is a generation-compare-and-swap which has semantics similar to RDCSS,
// but it does not create the intermediate object except in the case of
// failures that occur due to the snapshot being taken. This ensures that the
// write occurs only if the Ctrie root generation has remained the same in
// addition to the I-node having the expected value.
func gcas(in *iNode, old, n *mainNode, ct *Ctrie) bool {
	prevPtr := (*unsafe.Pointer)(unsafe.Pointer(&n.prev))
	atomic.StorePointer(prevPtr, unsafe.Pointer(old))
	if atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&in.main)),
		unsafe.Pointer(old), unsafe.Pointer(n)) {
		gcasComplete(in, n, ct)
		return atomic.LoadPointer(prevPtr) == nil
	}
	return false
}

// gcasRead performs a GCAS-linearizable read of the I-node's main node.
func gcasRead(in *iNode, ctrie *Ctrie) *mainNode {
	m := (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&in.main))))
	prev := (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&m.prev))))
	if prev == nil {
		return m
	}
	return gcasComplete(in, m, ctrie)
}

// gcasComplete commits the GCAS operation.
func gcasComplete(i *iNode, m *mainNode, ctrie *Ctrie) *mainNode {
	for {
		if m == nil {
			return nil
		}
		prev := (*mainNode)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&m.prev))))
		root := ctrie.rdcssReadRoot(true)
		if prev == nil {
			return m
		}

		if prev.failed != nil {
			// Signals GCAS failure. Swap old value back into I-node.
			fn := prev.failed
			if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&i.main)),
				unsafe.Pointer(m), unsafe.Pointer(fn.prev)) {
				return fn.prev
			}
			m = (*mainNode)(atomic.LoadPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&i.main))))
			continue
		}

		if root.gen == i.gen && !ctrie.readOnly {
			// Commit GCAS.
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&m.prev)), unsafe.Pointer(prev), nil) {
				return m
			}
			continue
		}

		// Generations did not match. Store failed node on prev to signal
		// I-node's main node must be set back to the previous value.
		atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&m.prev)),
			unsafe.Pointer(prev),
			unsafe.Pointer(&mainNode{failed: prev}))
		m = (*mainNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&i.main))))
		return gcasComplete(i, m, ctrie)
	}
}

// rdcssDescriptor is an intermediate struct which communicates the intent to
// replace the value in an I-node and check that the root's generation has not
// changed before committing to the new value.
type rdcssDescriptor struct {
	old       *iNode
	expected  *mainNode
	nv        *iNode
	committed bool
	token     int8
}

// HACK: see note in rdcssReadRoot.
const rdcssToken int8 = -1

// readRoot performs a linearizable read of the Ctrie root. This operation is
// prioritized so that if another thread performs a GCAS on the root, a
// deadlock does not occur.
func (c *Ctrie) readRoot() *iNode {
	return c.rdcssReadRoot(false)
}

// rdcssReadRoot performs a RDCSS-linearizable read of the Ctrie root with the
// given priority.
func (c *Ctrie) rdcssReadRoot(abort bool) *iNode {
	r := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.root)))
	desc := (*rdcssDescriptor)(r)
	// This is an awful and unsafe hack. Is there a better way to determine the
	// type of a pointer?
	if desc.token == rdcssToken {
		return c.rdcssComplete(abort)
	}
	return (*iNode)(r)
}

// rdcssRoot performs a RDCSS on the Ctrie root. This is used to create a
// snapshot of the Ctrie by copying the root I-node and setting it to a new
// generation.
func (c *Ctrie) rdcssRoot(old *iNode, expected *mainNode, nv *iNode) bool {
	desc := &rdcssDescriptor{
		old:      old,
		expected: expected,
		nv:       nv,
		token:    rdcssToken,
	}
	if c.casRoot(unsafe.Pointer(old), unsafe.Pointer(desc)) {
		c.rdcssComplete(false)
		return *(*bool)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&desc.committed))))
	}
	return false
}

// rdcssComplete commits the RDCSS operation.
func (c *Ctrie) rdcssComplete(abort bool) *iNode {
	for {
		r := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.root)))
		desc := (*rdcssDescriptor)(r)
		if desc.token != rdcssToken {
			return (*iNode)(r)
		}

		var (
			ov  = desc.old
			exp = desc.expected
			nv  = desc.nv
		)

		if abort {
			if c.casRoot(unsafe.Pointer(desc), unsafe.Pointer(ov)) {
				return ov
			}
			continue
		}

		oldeMain := gcasRead(ov, c)
		if oldeMain == exp {
			// Commit the RDCSS.
			if c.casRoot(unsafe.Pointer(desc), unsafe.Pointer(nv)) {
				desc.committed = true
				return nv
			}
			continue
		}
		if c.casRoot(unsafe.Pointer(desc), unsafe.Pointer(ov)) {
			return ov
		}
		continue
	}
}

// casRoot performs a CAS on the Ctrie root. The given pointers should either
// be rdcssDescriptors or iNodes.
func (c *Ctrie) casRoot(ov, nv unsafe.Pointer) bool {
	c.assertReadWrite()
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&c.root)), ov, nv)
}
