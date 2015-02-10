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

package palm

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Workiva/go-datastructures/futures"
	"github.com/Workiva/go-datastructures/queue"
)

type operation int

const (
	get operation = iota
	add
	remove
)

type recursiveBuild struct {
	keys   Keys
	nodes  []*node
	parent *node
}

type pending struct {
	reads    actions
	writes   Keys
	number   uint64
	signal   *futures.Future
	signaler chan interface{}
}

type ptree struct {
	root        *node
	ary, number uint64
	pending     *pending
	lock, write sync.Mutex
	waiter      *queue.Queue
}

func (ptree *ptree) initPending() {
	signaler := make(chan interface{})
	future := futures.New(signaler, 30*time.Minute)
	ptree.pending = &pending{
		signal:   future,
		signaler: signaler,
	}
}

func (ptree *ptree) listen() {
	for {
		_, err := ptree.waiter.Get(10)
		if err != nil {
			return
		}

		ptree.runOperations()
	}
}

func (ptree *ptree) runReads(reads actions, wg *sync.WaitGroup) {
	var key Key
	var i uint64
	for _, action := range reads {
		for {
			key, i = action.getKey()
			if key == nil {
				break
			}

			n := getParent(ptree.root, key)
			switch action.operation() {
			case get:
				if n == nil {
					action.addResult(i, nil)
				}
				index := n.keys.search(key)
				k := n.keys.byPosition(index)
				if index < n.keys.len() && k.Compare(key) == 0 {
					action.addResult(i, k)
				} else {
					action.addResult(i, nil)
				}
			}
		}
	}

	wg.Done()
}

func (ptree *ptree) runOperations() {
	ptree.lock.Lock()
	toPerform := ptree.pending
	ptree.initPending()
	ptree.lock.Unlock()

	if toPerform.number == 0 {
		return
	}

	var wg sync.WaitGroup
	var offset int
	var inserts []*node
	wg.Add(1) // for the gets

	if len(toPerform.writes) > 0 {
		inserts = make([]*node, len(toPerform.writes))
		chunks := chunkKeys(toPerform.writes, 8)
		wg.Add(len(chunks)) // for the inserts
		go ptree.runReads(toPerform.reads, &wg)

		for _, chunk := range chunks {
			go func(offset int, chunk Keys) {
				for i, k := range chunk {
					n := getParent(ptree.root, k)
					inserts[offset+i] = n
				}

				wg.Done()
			}(offset, chunk)
			offset += len(chunk)
		}
	} else {
		go ptree.runReads(toPerform.reads, &wg)
	}

	wg.Wait()
	if len(toPerform.writes) == 0 {
		return
	}
	writeOperations := make(map[*node]Keys)
	for i, n := range inserts {
		writeOperations[n] = append(writeOperations[n], toPerform.writes[i])
	}

	toPerform.signaler <- true
	ptree.runAdds(writeOperations)
}

func (ptree *ptree) recursiveSplit(n, parent, left *node, nodes *[]*node, keys *Keys) {
	if !n.needsSplit(ptree.ary) {
		return
	}

	key, l, r := n.split()
	if left != nil {
		left.right = l
	}
	l.parent = parent
	r.parent = parent
	*keys = append(*keys, key)
	*nodes = append(*nodes, l, r)
	ptree.recursiveSplit(l, parent, left, nodes, keys)
	ptree.recursiveSplit(r, parent, l, nodes, keys)
}

func (ptree *ptree) recursiveAdd(layer map[*node][]*recursiveBuild, setRoot bool) {
	if len(layer) == 0 {
		return
	}

	if setRoot && len(layer) > 1 {
		panic(`SHOULD ONLY HAVE ONE ROOT`)
	}

	q := queue.New(int64(len(layer)))
	for _, rbs := range layer {
		q.Put(rbs)
	}

	layer = make(map[*node][]*recursiveBuild, len(layer))
	dummyRoot := &node{
		keys:  newKeys(),
		nodes: newNodes(),
	}
	queue.ExecuteInParallel(q, func(ifc interface{}) {
		rbs := ifc.([]*recursiveBuild)

		if len(rbs) == 0 {
			return
		}

		n := rbs[0].parent
		if setRoot {
			ptree.root = n
		}

		parent := n.parent
		if parent == nil {
			parent = dummyRoot
			setRoot = true
		}

		for _, rb := range rbs {
			for i, k := range rb.keys {
				if n.keys.len() == 0 {
					n.keys.insert(k)
					n.nodes.push(rb.nodes[i*2])
					n.nodes.push(rb.nodes[i*2+1])
					continue
				}

				n.keys.insert(k)
				index := n.search(k)
				n.nodes.replaceAt(index, rb.nodes[i*2])
				n.nodes.insertAt(index+1, rb.nodes[i*2+1])
			}
		}

		if n.needsSplit(ptree.ary) {
			keys := make(Keys, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			ptree.recursiveSplit(n, parent, nil, &nodes, &keys)
			ptree.write.Lock()
			layer[parent] = append(
				layer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			ptree.write.Unlock()
		}
	})

	ptree.recursiveAdd(layer, setRoot)
}

func (ptree *ptree) runAdds(addOperations map[*node]Keys) {
	if len(addOperations) == 0 {
		return
	}

	q := queue.New(int64(len(addOperations)))

	for n := range addOperations {
		q.Put(n)
	}

	nextLayer := make(map[*node][]*recursiveBuild)
	dummyRoot := &node{
		keys:  newKeys(),
		nodes: newNodes(),
	} // constructed in case we need it
	var needRoot uint64
	queue.ExecuteInParallel(q, func(ifc interface{}) {
		n := ifc.(*node)
		keys := addOperations[n]

		if len(keys) == 0 {
			return
		}

		parent := n.parent
		if parent == nil {
			parent = dummyRoot
			atomic.AddUint64(&needRoot, 1)
		}

		for _, key := range keys {
			oldKey := n.keys.insert(key)
			if oldKey == nil {
				atomic.AddUint64(&ptree.number, 1)
			}
		}

		if n.needsSplit(ptree.ary) {
			keys := make(Keys, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			ptree.recursiveSplit(n, parent, nil, &nodes, &keys)
			ptree.write.Lock()
			nextLayer[parent] = append(
				nextLayer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			ptree.write.Unlock()
		}
	})

	setRoot := needRoot > 0

	ptree.recursiveAdd(nextLayer, setRoot)
}

// Insert will add the provided keys to the tree.
func (ptree *ptree) Insert(keys ...Key) {
	var signaler *futures.Future
	ptree.lock.Lock()
	ptree.pending.writes = append(ptree.pending.writes, keys...)
	ptree.pending.number += uint64(len(keys))
	signaler = ptree.pending.signal
	ptree.lock.Unlock()

	ptree.waiter.Put(true)
	signaler.GetResult()
}

// Get will retrieve a list of keys from the provided keys.
func (ptree *ptree) Get(keys ...Key) Keys {
	ga := newGetAction(keys)
	ptree.lock.Lock()
	ptree.pending.reads = append(ptree.pending.reads, ga)
	ptree.pending.number += uint64(len(keys))
	ptree.lock.Unlock()

	ptree.waiter.Put(true)
	return <-ga.completer
}

// Len returns the number of items in the tree.
func (ptree *ptree) Len() uint64 {
	return atomic.LoadUint64(&ptree.number)
}

// Dispose will clean up any resources used by this tree.  This
// must be called to prevent a memory leak.
func (ptree *ptree) Dispose() {
	ptree.waiter.Dispose()
}

func (ptree *ptree) print(output *log.Logger) {
	println(`PRINTING TREE`)
	if ptree.root == nil {
		return
	}

	ptree.root.print(output)
}

func newTree(ary uint64) *ptree {
	ptree := &ptree{
		root:    newNode(true, newKeys(), newNodes()),
		ary:     ary,
		pending: &pending{},
		waiter:  queue.New(10),
	}
	ptree.initPending()
	go ptree.listen()
	return ptree
}

// New will allocate, initialize, and return a new B-Tree based
// on PALM principles.  This type of tree is suited for in-memory
// indices in a multi-threaded environment.
func New(ary uint64) BTree {
	return newTree(ary)
}
