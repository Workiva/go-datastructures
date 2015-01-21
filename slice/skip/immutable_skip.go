/*
Copyright 2014 Workiva, LLC

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

package skip

import (
	"log"

	"github.com/Workiva/go-datastructures/hashmap/fastinteger"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

type ImmutableSkipList struct {
	maxLevel, level uint8
	head            *node
	num             uint64
}

func (isl *ImmutableSkipList) copy() *ImmutableSkipList {
	return &ImmutableSkipList{
		maxLevel: isl.maxLevel,
		level:    isl.level,
		head:     isl.head.copy(),
		num:      isl.num,
	}
}

// init will initialize this skiplist.  The parameter is expected
// to be of some uint type which will set this skiplist's maximum
// level.
func (isl *ImmutableSkipList) init(ifc interface{}) {
	switch ifc.(type) {
	case uint8:
		isl.maxLevel = 8
	case uint16:
		isl.maxLevel = 16
	case uint32:
		isl.maxLevel = 32
	case uint64, uint:
		isl.maxLevel = 64
	}
	isl.head = newNode(nil, isl.maxLevel)
}

func (isl *ImmutableSkipList) search(key uint64, update []*node,
	fihm *fastinteger.FastIntegerHashMap) *node {

	if isl.num == 0 { // nothing in the list
		return nil
	}

	needsUpdate := update != nil

	var offset uint8
	n := isl.head // this has already been copied
	log.Printf(`SEARCH HEAD: %p HEAD: %+v`, n, n)
	for i := uint8(0); i <= isl.level; i++ {
		offset = isl.level - i
		for n.forward[offset] != nil && n.forward[offset].entry.Key() < key {
			// we need to make a copy now
			var f *node
			if needsUpdate {
				log.Printf(`n.forward[offset]: %p`, n.forward[offset])
				println(`HIT THIS`)
				if fihm != nil {
					if fihm.Exists(n.forward[offset].entry.Key()) {
						f = n.forward[offset]
					} else {
						f = n.forward[offset].copy()
						fihm.Set(f.entry.Key(), 0)
						for i := int(offset) - 1; i >= 0; i-- {
							if n.forward[i] == n.forward[offset] {
								n.forward[i] = f
							} else {
								break
							}
						}
					}
				}
				log.Printf(`F: %p`, f)
				n.forward[offset] = f
				n = f
			} else {
				n = n.forward[offset]
			}
		}

		log.Printf(`SEARCH N POINTER: %p`, n)
		if update != nil {
			update[offset] = n
		}
	}

	f := n.forward[0]

	if needsUpdate && f != nil {
		if fihm != nil {
			if !fihm.Exists(f.entry.Key()) {
				f = f.copy()
				n.forward[0] = f
				fihm.Set(f.entry.Key(), 0)
			}
		}
	}

	return f
}

// Get will retrieve values associated with the keys provided.  If an
// associated value could not be found, a nil is returned in its place.
// This is an O(log n) operation.
func (isl *ImmutableSkipList) Get(keys ...uint64) Entries {
	entries := make(Entries, 0, len(keys))

	var n *node
	for _, key := range keys {
		n = isl.search(key, nil, nil)
		if n != nil && n.entry.Key() == key {
			entries = append(entries, n.entry)
		} else {
			entries = append(entries, nil)
		}
	}

	return entries
}

func (isl *ImmutableSkipList) copyNodes(ns nodes, fihm *fastinteger.FastIntegerHashMap) nodes {
	cp := make(nodes, len(ns))
	for i, n := range ns {
		if n == nil {
			break
		}

		if fihm.Exists(n.entry.Key()) {
			cp[i] = n
			continue
		}

		cp[i] = n.copy()
		fihm.Set(n.entry.Key(), 0) // we only care about key existence.
	}

	return cp
}

func (isl *ImmutableSkipList) insert(entry Entry, cache nodes,
	fihm *fastinteger.FastIntegerHashMap) Entry {

	n := isl.search(entry.Key(), cache, fihm)
	if n != nil && n.key() == entry.Key() { // a simple update in this case
		oldEntry := n.entry
		n = n.copy()
		n.entry = entry
		return oldEntry
	}
	isl.num++

	nodeLevel := generateLevel(isl.maxLevel)
	if nodeLevel > isl.level {
		for i := isl.level; i <= nodeLevel; i++ {
			cache[i] = isl.head
		}
		isl.level = nodeLevel
	}

	nn := newNode(entry, isl.maxLevel)
	log.Printf(`CACHE: %+v, nodelevel: %+v, node: %p`, cache, nodeLevel, nn)
	for i := uint8(0); i < nodeLevel; i++ {
		log.Printf(`I: %+v, cache[i].forward[i]: %+v, cache[i]: %+v`, i, cache[i].forward[i], cache[i])
		nn.forward[i] = cache[i].forward[i]
		cache[i].forward[i] = nn
	}
	cache[0].forward[0] = nn

	return nil
}

func (isl *ImmutableSkipList) Insert(entries ...Entry) (*ImmutableSkipList, Entries) {
	if len(entries) == 0 {
		return isl, Entries{}
	}
	cp := isl.copy()
	fihm := fastinteger.New(uint64(len(entries)))
	cache := make(nodes, isl.maxLevel)
	overwritten := make(Entries, 0, len(entries))
	for _, e := range entries {
		overwritten = append(overwritten, cp.insert(e, cache, fihm))
		cache.reset()
	}

	return cp, overwritten
}

func (isl *ImmutableSkipList) delete(key uint64,
	cache nodes, fihm *fastinteger.FastIntegerHashMap) Entry {

	n := isl.search(key, cache, fihm)
	if n == nil || n.entry.Key() != key {
		return nil
	}

	isl.num--

	for i := uint8(0); i <= isl.level; i++ {
		if cache[i].forward[i] != n {
			break
		}

		cache[i].forward[i] = n.forward[i]
	}

	for isl.level > 0 && isl.head.forward[isl.level] == nil {
		isl.level = isl.level - 1
	}

	return n.entry
}

func (isl *ImmutableSkipList) Delete(keys ...uint64) (*ImmutableSkipList, Entries) {
	if len(keys) == 0 {
		return isl, Entries{}
	}

	cp := isl.copy()
	fihm := fastinteger.New(uint64(len(keys)))
	cache := make(nodes, isl.maxLevel)
	deleted := make(Entries, 0, len(keys))
	for _, key := range keys {
		deleted = append(deleted, cp.delete(key, cache, fihm))
		cache.reset()
	}

	return cp, deleted
}

func (isl *ImmutableSkipList) iter(key uint64) *iterator {
	n := isl.search(key, nil, nil)
	if n == nil {
		return nilIterator()
	}

	return &iterator{
		first: true,
		n:     n,
	}
}

func (isl *ImmutableSkipList) Iter(key uint64) Iterator {
	return isl.iter(key)
}

func (isl *ImmutableSkipList) Len() uint64 {
	return isl.num
}

func NewImmutable(ifc interface{}) *ImmutableSkipList {
	isl := &ImmutableSkipList{}
	isl.init(ifc)
	return isl
}
