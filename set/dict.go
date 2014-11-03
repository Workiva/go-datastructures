package set

import "sync"

var pool = sync.Pool{}

// Set is an implementation of ISet using the builtin map type. Set is threadsafe.
type Set struct {
	items     map[interface{}]bool
	lock      sync.RWMutex
	flattened []interface{}
}

// Add will add the provided items to the set.
func (set *Set) Add(items ...interface{}) {
	set.lock.Lock()
	defer set.lock.Unlock()

	set.flattened = nil
	for _, item := range items {
		set.items[item] = true
	}
}

// Remove will remove the given items from the set.
func (set *Set) Remove(items ...interface{}) {
	set.lock.Lock()
	defer set.lock.Unlock()

	set.flattened = nil
	for _, item := range items {
		delete(set.items, item)
	}
}

// Exists returns a bool indicating if the given item exists in the set.
func (set *Set) Exists(item interface{}) bool {
	set.lock.RLock()
	defer set.lock.RUnlock()

	_, ok := set.items[item]
	return ok
}

// Flatten will return a list of the items in the set.
func (set *Set) Flatten() []interface{} {
	set.lock.Lock()
	defer set.lock.Unlock()

	if set.flattened != nil {
		return set.flattened
	}

	set.flattened = make([]interface{}, 0, len(set.items))
	for item := range set.items {
		set.flattened = append(set.flattened, item)
	}
	return set.flattened
}

// Len returns the number of items in the set.
func (set *Set) Len() int64 {
	set.lock.RLock()
	defer set.lock.RUnlock()

	return int64(len(set.items))
}

// Clear will remove all items from the set.
func (set *Set) Clear() {
	set.lock.Lock()
	defer set.lock.Unlock()

	set.items = map[interface{}]bool{}
}

// All returns a bool indicating if all of the supplied items exist in the set.
func (set *Set) All(items ...interface{}) bool {
	set.lock.RLock()
	defer set.lock.RUnlock()

	for _, item := range items {
		if _, ok := set.items[item]; !ok {
			return false
		}
	}

	return true
}

// Dispose will add this set back into the pool.
func (set *Set) Dispose() {
	set.lock.Lock()
	defer set.lock.Unlock()

	for k := range set.items {
		delete(set.items, k)
	}

	//this is so we don't hang onto any references
	for i := 0; i < len(set.flattened); i++ {
		set.flattened[i] = nil
	}

	set.flattened = set.flattened[:0]
	pool.Put(set)
}

// New is the constructor for sets.  It will pull from a reuseable memory pool if it can.
// Takes a list of items to initialize the set with.
func New(items ...interface{}) *Set {
	set := pool.Get().(*Set)
	for _, item := range items {
		set.items[item] = true
	}

	return set
}

func init() {
	pool.New = func() interface{} {
		return &Set{
			items: make(map[interface{}]bool, 10),
		}
	}
}
