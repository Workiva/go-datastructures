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
	//"runtime"
	"sync"
	"sync/atomic"
	//"time"

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

type ptree struct {
	root                      *node
	ary, number, cacheBuffer1 uint64
	actions                   *queue.RingBuffer
	bufferSize, cacheBuffer2  uint64
	cache                     []interface{}
	cacheBuffer3              uint64
	disposed, cacheBuffer4    uint64
	running, cacheBuffer5     uint64
}

func (ptree *ptree) checkAndRun(ifc interface{}) {
	if ptree.actions.Len() > 0 {
		ptree.actions.Put(ifc)
		if atomic.CompareAndSwapUint64(&ptree.running, 0, 1) {
			var a interface{}
			var err error
			for ptree.actions.Len() > 0 {
				a, err = ptree.actions.Get()
				if err != nil {
					return
				}
				ptree.cache = append(ptree.cache, a)
				if uint64(len(ptree.cache)) >= ptree.bufferSize {
					break
				}
			}

			go ptree.runOperations(ptree.cache)
		}
	} else if ifc != nil {
		if atomic.CompareAndSwapUint64(&ptree.running, 0, 1) {
			switch ifc.(type) {
			case *getAction:
				for ptree.read(ifc.(*getAction)) {
				}
				ptree.reset()
			default:
				ptree.runOperations([]interface{}{ifc})
			}
		} else {
			ptree.actions.Put(ifc)
			ptree.checkAndRun(nil)
		}
	}
}

func (ptree *ptree) init(bufferSize, ary uint64) {
	ptree.bufferSize = bufferSize
	ptree.ary = ary
	ptree.cache = make([]interface{}, 0, bufferSize)
	ptree.root = newNode(true, newKeys(), newNodes())
	ptree.actions = queue.NewRingBuffer(ptree.bufferSize)
}

func (ptree *ptree) read(action action) bool {
	key, i := action.getKey()
	if key == nil {
		return false
	}

	n := getParent(ptree.root, key)
	if n == nil {
		action.addResult(i, nil)
	} else {
		k, _ := n.keys.withPosition(key)
		if k == nil {
			action.addResult(i, nil)
		} else {
			action.addResult(i, k)
		}
	}
	return true
}

func (ptree *ptree) runReads(reads actions, wg *sync.WaitGroup) {
	for _, action := range reads {
		for ptree.read(action) {
		}
	}

	wg.Done()
}

func (ptree *ptree) reset() {
	for i := range ptree.cache {
		ptree.cache[i] = nil
	}

	ptree.cache = ptree.cache[:0]
	atomic.StoreUint64(&ptree.running, 0)
	ptree.checkAndRun(nil)
}

func (ptree *ptree) runOperations(xns []interface{}) {
	//println(`RUNNING OPERATIONS`)
	reads := make(actions, 0, len(xns)/2)
	writes := make(Keys, 0, len(xns)/2)

	for _, a := range xns {
		switch a.(type) {
		case *insertAction:
			writes = append(writes, a.(*insertAction).keys...)
		case *getAction:
			reads = append(reads, a.(*getAction))
		}
	}

	var wg sync.WaitGroup
	var offset int
	var inserts []*node
	wg.Add(1) // for the gets

	if len(writes) > 0 {
		inserts = make([]*node, len(writes))
		chunks := chunkKeys(writes, 8)
		wg.Add(len(chunks)) // for the inserts
		go ptree.runReads(reads, &wg)

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
		go ptree.runReads(reads, &wg)
	}

	wg.Wait()
	if len(writes) == 0 {
		ptree.reset()
		return
	}
	writeOperations := make(map[*node]Keys)
	for i, n := range inserts {
		writeOperations[n] = append(writeOperations[n], writes[i])
	}

	ptree.runAdds(writeOperations)
	for _, a := range xns {
		if ia, ok := a.(*insertAction); ok {
			ia.complete()
		}
	}

	ptree.reset()
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

	ifs := make(interfaces, 0, len(layer))
	for _, rbs := range layer {
		ifs = append(ifs, rbs)
	}

	var write sync.Mutex
	layer = make(map[*node][]*recursiveBuild, len(layer))
	dummyRoot := &node{
		keys:  newKeys(),
		nodes: newNodes(),
	}
	executeInterfacesInParallel(ifs, func(ifc interface{}) {
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
			write.Lock()
			layer[parent] = append(
				layer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			write.Unlock()
		}
	})

	ptree.recursiveAdd(layer, setRoot)
}

func (ptree *ptree) runAdds(addOperations map[*node]Keys) {
	if len(addOperations) == 0 {
		return
	}

	ifs := make(interfaces, 0, len(addOperations))
	for n := range addOperations {
		ifs = append(ifs, n)
	}

	var write sync.Mutex
	nextLayer := make(map[*node][]*recursiveBuild)
	dummyRoot := &node{} // constructed in case we need it
	var needRoot uint64
	executeInterfacesInParallel(ifs, func(ifc interface{}) {
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
			write.Lock()
			nextLayer[parent] = append(
				nextLayer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			write.Unlock()
		}
	})

	setRoot := needRoot > 0
	if setRoot {
		dummyRoot.keys = newKeys()
		dummyRoot.nodes = newNodes()
	}

	ptree.recursiveAdd(nextLayer, setRoot)
}

// Insert will add the provided keys to the tree.
func (ptree *ptree) Insert(keys ...Key) {
	ia := newInsertAction(keys)
	ptree.checkAndRun(ia)
	ia.completer.Wait()
}

// Get will retrieve a list of keys from the provided keys.
func (ptree *ptree) Get(keys ...Key) Keys {
	ga := newGetAction(keys)
	//t0 := time.Now()
	ptree.checkAndRun(ga)
	//log.Printf(`ALLOCATE TIME: %+v`, time.Since(t0).Nanoseconds())
	ga.completer.Wait()
	//log.Printf(`GET TIME: %+v`, time.Since(t0).Nanoseconds())
	return ga.result
}

// Len returns the number of items in the tree.
func (ptree *ptree) Len() uint64 {
	return atomic.LoadUint64(&ptree.number)
}

// Dispose will clean up any resources used by this tree.  This
// must be called to prevent a memory leak.
func (ptree *ptree) Dispose() {
	ptree.actions.Dispose()
	atomic.StoreUint64(&ptree.disposed, 1)
}

func (ptree *ptree) print(output *log.Logger) {
	println(`PRINTING TREE`)
	if ptree.root == nil {
		return
	}

	ptree.root.print(output)
}

func newTree(bufferSize, ary uint64) *ptree {
	ptree := &ptree{}
	ptree.init(bufferSize, ary)
	return ptree
}

// New will allocate, initialize, and return a new B-Tree based
// on PALM principles.  This type of tree is suited for in-memory
// indices in a multi-threaded environment.
func New(bufferSize, ary uint64) BTree {
	return newTree(bufferSize, ary)
}
