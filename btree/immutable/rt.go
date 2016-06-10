//go:generate msgp -tests=false -io=false

package btree

import (
	"sync"

	"github.com/satori/go.uuid"
)

// context is used to keep track of the nodes in this mutable
// that have been created.  This is basically any node that had
// to be touched to perform mutations.  Further mutations will visit
// this context first so we don't have to constantly copy if
// we don't need to.
type context struct {
	lock      sync.RWMutex
	seenNodes map[string]*Node
}

func (c *context) nodeExists(id ID) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.seenNodes[string(id)]
	return ok
}

func (c *context) addNode(n *Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.seenNodes[string(n.ID)] = n
}

func (c *context) getNode(id ID) *Node {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.seenNodes[string(id)]
}

func newContext() *context {
	return &context{
		seenNodes: make(map[string]*Node, 10),
	}
}

// Tr itself is exported so that the code generated for serialization/deserialization
// works on Tr.  Exported fields on Tr are those fields that need to be
// serialized.
type Tr struct {
	UUID      ID  `msg:"u"`
	Count     int `msg:"c"`
	config    Config
	Root      ID `msg:"r"`
	cacher    *cacher
	context   *context
	NodeWidth int `msg:"nw"`
	mutable   bool
}

func (t *Tr) createRoot() *Node {
	n := newNode()
	n.IsLeaf = true
	return n
}

// contextOrCachedNode is a convenience function for either fetching
// a node from the context or persistence.
func (t *Tr) contextOrCachedNode(id ID, cache bool) (*Node, error) {
	if t.context != nil {
		n := t.context.getNode(id)
		if n != nil {
			return n, nil
		}
	}

	return t.cacher.getNode(t, id, cache)
}

func (t *Tr) ID() ID {
	return t.UUID
}

// toBytes encodes this tree into a byte array.  Panics if unable
// as this error has to be fixed in code.
func (t *Tr) toBytes() []byte {
	buf, err := t.MarshalMsg(nil)
	if err != nil {
		panic(`unable to encode tree`)
	}

	return buf
}

// reset is called on a tree to empty the context and clear the cache.
func (t *Tr) reset() {
	t.cacher.clear()
	t.context = nil
}

// commit will gather up all created nodes and serialize them into
// items that can be persisted.
func (t *Tr) commit() []*Payload {
	items := make([]*Payload, 0, len(t.context.seenNodes))
	for _, n := range t.context.seenNodes {
		n.ChildValues, n.ChildKeys = n.flatten()
		buf, err := n.MarshalMsg(nil)
		if err != nil {
			panic(`unable to encode node`)
		}

		n.ChildValues, n.ChildKeys = nil, nil
		item := &Payload{n.ID, buf}
		items = append(items, item)
	}

	return items
}

func (t *Tr) copyNode(n *Node) *Node {
	if t.context.nodeExists(n.ID) {
		return n
	}

	cp := n.copy()
	t.context.addNode(cp)
	return cp
}

func (t *Tr) Len() int {
	return t.Count
}

func (t *Tr) AsMutable() MutableTree {
	return &Tr{
		Count:     t.Count,
		UUID:      uuid.NewV4().Bytes(),
		Root:      t.Root,
		config:    t.config,
		cacher:    t.cacher,
		context:   newContext(),
		NodeWidth: t.NodeWidth,
		mutable:   true,
	}
}

func (t *Tr) Commit() (ReadableTree, error) {
	t.NodeWidth = t.config.NodeWidth
	items := make([]*Payload, 0, len(t.context.seenNodes))
	items = append(items, t.commit()...)

	// save self
	items = append(items, &Payload{t.ID(), t.toBytes()})

	err := t.config.Persister.Save(items...)
	if err != nil {
		return nil, err
	}

	t.reset()
	t.context = nil
	return t, nil
}

func treeFromBytes(p Persister, data []byte, comparator Comparator) (*Tr, error) {
	t := &Tr{}
	_, err := t.UnmarshalMsg(data)
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig(p, comparator)
	if t.NodeWidth > 0 {
		cfg.NodeWidth = t.NodeWidth
	}
	t.config = cfg
	t.cacher = newCacher(cfg.Persister)

	return t, nil
}

func newTree(cfg Config) *Tr {
	return &Tr{
		config: cfg,
		UUID:   uuid.NewV4().Bytes(),
		cacher: newCacher(cfg.Persister),
	}
}

// New creates a new ReadableTree using the provided config.
func New(cfg Config) ReadableTree {
	return newTree(cfg)
}

// Load returns a ReadableTree from persistence.  The provided
// config should contain a persister that can be used for this purpose.
// An error is returned if the tree could not be found or an error
// occurred in the persistence layer.
func Load(p Persister, id []byte, comparator Comparator) (ReadableTree, error) {
	items, err := p.Load(id)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 || items[0] == nil {
		return nil, ErrTreeNotFound
	}

	rt, err := treeFromBytes(p, items[0].Payload, comparator)
	if err != nil {
		return nil, err
	}

	return rt, nil
}
