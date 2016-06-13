package btree

import (
	"sync"
	"time"

	"github.com/Workiva/go-datastructures/futures"
)

// cacher provides a convenient construct for retrieving,
// storing, and caching nodes; basically wrapper persister with a caching layer.
// This ensures that we don't have to constantly
// run to the persister to fetch nodes we are using over and over again.
// TODO: this should probably evict items from the cache if the cache gets
// too full.
type cacher struct {
	lock      sync.Mutex
	cache     map[string]*futures.Future
	persister Persister
}

func (c *cacher) asyncLoadNode(t *Tr, key ID, completer chan interface{}) {
	n, err := c.loadNode(t, key)
	if err != nil {
		completer <- err
		return
	}

	if n == nil {
		completer <- ErrNodeNotFound
		return
	}

	completer <- n
}

// clear deletes all items from the cache.
func (c *cacher) clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache = make(map[string]*futures.Future, 10)
}

// deleteFromCache will remove the provided ID from the cache.  This
// is a threadsafe operation.
func (c *cacher) deleteFromCache(id ID) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.cache, string(id))
}

func (c *cacher) loadNode(t *Tr, key ID) (*Node, error) {
	items, err := c.persister.Load(key)
	if err != nil {
		return nil, err
	}

	n, err := nodeFromBytes(t, items[0].Payload)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// getNode will return a Node matching the provided id.  An error is returned
// if the cacher could not go to the persister or the node could not be found.
// All found nodes are cached so subsequent calls should be faster than
// the initial.  This blocks until the node is loaded, but is also threadsafe.
func (c *cacher) getNode(t *Tr, key ID, useCache bool) (*Node, error) {
	if !cache {
		return c.loadNode(t, key)
	}

	c.lock.Lock()
	future, ok := c.cache[string(key)]
	if ok {
		c.lock.Unlock()
		ifc, err := future.GetResult()
		if err != nil {
			return nil, err
		}

		return ifc.(*Node), nil
	}

	completer := make(chan interface{}, 1)
	future = futures.New(completer, 30*time.Second)
	c.cache[string(key)] = future
	c.lock.Unlock()

	go c.asyncLoadNode(t, key, completer)

	ifc, err := future.GetResult()
	if err != nil {
		c.deleteFromCache(key)
		return nil, err
	}

	if err, ok := ifc.(error); ok {
		c.deleteFromCache(key)
		return nil, err
	}

	return ifc.(*Node), nil
}

// newCacher is the constructor for a cacher that caches nodes for
// an indefinite period of time.
func newCacher(persister Persister) *cacher {
	return &cacher{
		persister: persister,
		cache:     make(map[string]*futures.Future, 10),
	}
}
