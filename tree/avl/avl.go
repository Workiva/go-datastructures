package avl

import (
	"log"
	"math"
)

func init() {
	log.Println(`I HATE THIS.`)
}

type Immutable struct {
	root   *node
	number uint64
	dummy  node
	cache  nodes
}

func (immutable *Immutable) copy() *Immutable {
	var root *node
	if immutable.root != nil {
		root = immutable.root.copy()
	}
	cp := &Immutable{
		root:   root,
		number: immutable.number,
		dummy:  *newNode(nil),
	}
	return cp
}

func (immutable *Immutable) resetCache() {
	immutable.cache.reset()
}

func (immutable *Immutable) resetDummy() {
	immutable.dummy.children[0], immutable.dummy.children[1] = nil, nil
	immutable.dummy.balance = 0
}

func (immutable *Immutable) init() {
	immutable.dummy = node{
		children: [2]*node{},
	}
	immutable.cache = make(nodes, 64) // this should cover every number in the 64 bit universe
}

func (immutable *Immutable) get(entry Entry) Entry {
	n := immutable.root
	var result int
	for n != nil {
		switch result = n.entry.Compare(entry); {
		case result == 0:
			return n.entry
		case result > 0:
			n = n.children[0]
		case result < 0:
			n = n.children[1]
		}
	}

	return nil
}

func (immutable *Immutable) Get(entries ...Entry) Entries {
	result := make(Entries, 0, len(entries))
	for _, e := range entries {
		result = append(result, immutable.get(e))
	}

	return result
}

// Len returns the number of items in this immutable.
func (immutable *Immutable) Len() uint64 {
	return immutable.number
}

func (immutable *Immutable) insert(entry Entry) Entry {
	if immutable.root == nil {
		immutable.root = newNode(entry)
		immutable.number++
		return nil
	}

	immutable.resetDummy()
	var (
		dummy           = immutable.dummy
		p, s, q         *node
		dir, normalized int
		helper          = &dummy
	)

	// set this AFTER clearing dummy
	helper.children[1] = immutable.root
	for s, p = helper.children[1], helper.children[1]; ; {
		dir = p.entry.Compare(entry)

		normalized = normalizeComparison(dir)
		if dir > 0 { // go left
			if p.children[0] != nil {
				q = p.children[0].copy()
				p.children[0] = q
			} else {
				q = nil
			}
		} else if dir < 0 { // go right
			if p.children[1] != nil {
				q = p.children[1].copy()
				p.children[1] = q
			} else {
				q = nil
			}
		} else { // equality
			oldEntry := p.entry
			p.entry = entry
			return oldEntry
		}
		if q == nil {
			break
		}

		if q.balance != 0 {
			helper = p
			s = q
		}
		p = q
	}

	immutable.number++
	q = newNode(entry)
	p.children[normalized] = q

	immutable.root = dummy.children[1]
	for p = s; p != q; p = p.children[normalized] {
		normalized = normalizeComparison(p.entry.Compare(entry))
		if normalized == 0 {
			p.balance += -1
		} else {
			p.balance += 1
		}
	}

	q = s

	if math.Abs(float64(s.balance)) > 1 {
		normalized = normalizeComparison(s.entry.Compare(entry))
		s = insertBalance(s, normalized)
	}

	if q == dummy.children[1] {
		immutable.root = s
	} else {
		helper.children[intFromBool(helper.children[1] == q)] = s
	}
	return nil
}

func (immutable *Immutable) Insert(entries ...Entry) (*Immutable, Entries) {
	if len(entries) == 0 {
		return immutable, Entries{}
	}

	overwritten := make(Entries, 0, len(entries))
	cp := immutable.copy()
	for _, e := range entries {
		overwritten = append(overwritten, cp.insert(e))
	}

	return cp, overwritten
}

func (immutable *Immutable) remove(entry Entry) Entry {
	return nil
}

func insertBalance(root *node, dir int) *node {
	n := root.children[dir]
	var bal int8
	if dir == 0 {
		bal = -1
	} else {
		bal = 1
	}

	if n.balance == bal {
		root.balance, n.balance = 0, 0
		root = rotate(root, takeOpposite(dir))
	} else { /* n->balance == -bal */
		adjustBalance(root, dir, int(bal))
		root = doubleRotate(root, takeOpposite(dir))
	}

	return root
}

func removeBalance(root *node, dir int, done *int) *node {
	n := root.children[takeOpposite(dir)]
	var bal int8
	if dir == 0 {
		bal = -1
	} else {
		bal = 1
	}

	if n.balance == -bal {
		root.balance, n.balance = 0, 0
		root = rotate(root, dir)
	} else if n.balance == bal {
		adjustBalance(root, takeOpposite(dir), int(-bal))
		root = doubleRotate(root, dir)
	} else {
		root.balance = -bal
		n.balance = bal
		root = rotate(root, dir)
		*done = 1
	}

	return root
}

func intFromBool(value bool) int {
	if value {
		return 1
	}

	return 0
}

func takeOpposite(value int) int {
	return 1 - value
}

func adjustBalance(root *node, dir, bal int) {
	n := root.children[dir]
	nn := n.children[takeOpposite(dir)]

	if nn.balance == 0 {
		root.balance, n.balance = 0, 0
	} else if int(nn.balance) == bal {
		root.balance = int8(-bal)
		n.balance = 0
	} else { /* nn->balance == -bal */
		root.balance = 0
		n.balance = int8(bal)
	}
	nn.balance = 0
}

func rotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	child := parent.children[otherDir]
	parent.children[otherDir] = child.children[dir]
	child.children[dir] = parent

	return child
}

func doubleRotate(parent *node, dir int) *node {
	otherDir := takeOpposite(dir)

	parent.children[otherDir] = rotate(parent.children[otherDir], otherDir)
	return rotate(parent, dir)
}

func normalizeComparison(i int) int {
	if i < 0 {
		return 1
	}

	if i > 0 {
		return 0
	}

	return -1
}

func NewImmutable() *Immutable {
	immutable := &Immutable{}
	immutable.init()
	return immutable
}
