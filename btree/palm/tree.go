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
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/common"
	"github.com/Workiva/go-datastructures/queue"
)

type operation int

const (
	get operation = iota
	add
	remove
)

const multiThreadAt = 1000 // number of keys before we multithread lookups

type recursiveBuild struct {
	keys   common.Comparators
	nodes  []*node
	parent *node
}

type ptree struct {
	root                    *node
	ary, number, bufferSize uint64
	actions                 *queue.RingBuffer
	cache                   []interface{}
	buffer0                 [8]uint64
	disposed                uint64
	buffer1                 [8]uint64
	running                 uint64
}

func (ptree *ptree) checkAndRun(action action) {
	if ptree.actions.Len() > 0 {
		if action != nil {
			ptree.actions.Put(action)
		}
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

			go ptree.operationRunner(ptree.cache, true)
		}
	} else if action != nil {
		if atomic.CompareAndSwapUint64(&ptree.running, 0, 1) {
			switch action.operation() {
			case get:
				ptree.read(action)
				action.complete()
				ptree.reset()
			case add:
				if len(action.keys()) > multiThreadAt {
					ptree.operationRunner(interfaces{action}, true)
				} else {
					ptree.operationRunner(interfaces{action}, false)
				}
			}
		} else {
			ptree.actions.Put(action)
			ptree.checkAndRun(nil)
		}
	}
}

func (ptree *ptree) init(bufferSize, ary uint64) {
	ptree.bufferSize = bufferSize
	ptree.ary = ary
	ptree.cache = make([]interface{}, 0, bufferSize)
	ptree.root = newNode(true, newKeys(ary), newNodes(ary))
	ptree.actions = queue.NewRingBuffer(ptree.bufferSize)
}

func (ptree *ptree) operationRunner(xns interfaces, threaded bool) {
	var writeOperations map[*node]common.Comparators
	var toComplete actions

	if threaded {
		writeOperations, toComplete = ptree.fetchKeys(xns)
	} else {
		writeOperations, toComplete = ptree.singleThreadedFetchKeys(xns)
	}

	ptree.runAdds(writeOperations)
	for _, a := range toComplete {
		a.complete()
	}

	ptree.reset()
}

func (ptree *ptree) read(action action) {
	for i, k := range action.keys() {
		n := getParent(ptree.root, k)
		if n == nil {
			action.keys()[i] = nil
		} else {
			key, _ := n.keys.withPosition(k)
			if key == nil {
				action.keys()[i] = nil
			} else {
				action.keys()[i] = key
			}
		}
	}
}

func (ptree *ptree) singleThreadedFetchKeys(xns interfaces) (map[*node]common.Comparators, actions) {
	for _, ifc := range xns {
		action := ifc.(action)
		for i, key := range action.keys() {
			n := getParent(ptree.root, key)
			switch action.operation() {
			case add:
				action.addNode(int64(i), n)
			case get:
				if n == nil {
					action.keys()[i] = nil
				} else {
					k, _ := n.keys.withPosition(key)
					if k == nil {
						action.keys()[i] = nil
					} else {
						action.keys()[i] = k
					}
				}
			}
		}
	}

	writeOperations := make(map[*node]common.Comparators, len(xns)/2)
	toComplete := make(actions, 0, len(xns)/2)
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add:
			for i, n := range action.nodes() {
				writeOperations[n] = append(writeOperations[n], action.keys()[i])
			}
			toComplete = append(toComplete, action)
		case get:
			action.complete()
		}
	}

	return writeOperations, toComplete
}

func (ptree *ptree) reset() {
	for i := range ptree.cache {
		ptree.cache[i] = nil
	}

	ptree.cache = ptree.cache[:0]
	atomic.StoreUint64(&ptree.running, 0)
	ptree.checkAndRun(nil)
}

func (ptree *ptree) fetchKeys(xns []interface{}) (map[*node]common.Comparators, actions) {
	var forCache struct {
		i      int64
		buffer [8]uint64 // different cache lines
		js     []int64
	}

	for j := 0; j < len(xns); j++ {
		forCache.js = append(forCache.js, -1)
	}
	numCPU := runtime.NumCPU()
	if numCPU > 1 {
		numCPU--
	}
	var wg sync.WaitGroup
	wg.Add(numCPU)

	for k := 0; k < numCPU; k++ {
		go func() {
			for {
				index := atomic.LoadInt64(&forCache.i)
				if index >= int64(len(xns)) {
					break
				}
				action := xns[index].(action)

				j := atomic.AddInt64(&forCache.js[index], 1)
				if j > int64(len(action.keys())) { // someone else is updating i
					continue
				} else if j == int64(len(action.keys())) {
					atomic.StoreInt64(&forCache.i, index+1)
					continue
				}

				n := getParent(ptree.root, action.keys()[j])
				switch action.operation() {
				case add:
					action.addNode(j, n)
				case get:
					if n == nil {
						action.keys()[j] = nil
					} else {
						k, _ := n.keys.withPosition(action.keys()[j])
						if k == nil {
							action.keys()[j] = nil
						} else {
							action.keys()[j] = k
						}
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	writeOperations := make(map[*node]common.Comparators, len(xns)/2)
	toComplete := make(actions, 0, len(xns)/2)
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add:
			for i, n := range action.nodes() {
				writeOperations[n] = append(writeOperations[n], action.keys()[i])
			}
			toComplete = append(toComplete, action)
		case get:
			action.complete()
		}
	}

	return writeOperations, toComplete
}

func (ptree *ptree) splitNode(n, parent *node, nodes *[]*node, keys *common.Comparators) {
	if !n.needsSplit(ptree.ary) {
		return
	}

	length := n.keys.len()
	splitAt := ptree.ary - 1

	for i := splitAt; i < length; i += splitAt {
		offset := length - i
		k, left, right := n.split(offset, ptree.ary)
		left.right = right
		*keys = append(*keys, k)
		*nodes = append(*nodes, left, right)
		left.parent = parent
		right.parent = parent
	}
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
		if rbs[0].parent.parent == nil {
			setRoot = true
		}
		ifs = append(ifs, rbs)
	}

	var dummyRoot *node
	if setRoot {
		dummyRoot = &node{
			keys:  newKeys(ptree.ary),
			nodes: newNodes(ptree.ary),
		}
	}

	var write sync.Mutex
	layer = make(map[*node][]*recursiveBuild, len(layer))

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
			keys := make(common.Comparators, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			ptree.splitNode(n, parent, &nodes, &keys)
			write.Lock()
			layer[parent] = append(
				layer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			write.Unlock()
		}
	})

	ptree.recursiveAdd(layer, setRoot)
}

func (ptree *ptree) runAdds(addOperations map[*node]common.Comparators) {
	if len(addOperations) == 0 {
		return
	}

	var needRoot bool
	ifs := make(interfaces, 0, len(addOperations))
	for n := range addOperations {
		if n.parent == nil {
			needRoot = true
		}
		ifs = append(ifs, n)
	}

	var dummyRoot *node
	if needRoot {
		dummyRoot = &node{
			keys:  newKeys(ptree.ary),
			nodes: newNodes(ptree.ary),
		}
	}

	var write sync.Mutex
	nextLayer := make(map[*node][]*recursiveBuild)
	executeInterfacesInParallel(ifs, func(ifc interface{}) {
		n := ifc.(*node)
		keys := addOperations[n]

		if len(keys) == 0 {
			return
		}

		parent := n.parent
		if parent == nil {
			parent = dummyRoot
		}

		for _, key := range keys {
			oldKey := n.keys.insert(key)
			if oldKey == nil {
				atomic.AddUint64(&ptree.number, 1)
			}
		}

		if n.needsSplit(ptree.ary) {
			keys := make(common.Comparators, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			ptree.splitNode(n, parent, &nodes, &keys)
			write.Lock()
			nextLayer[parent] = append(
				nextLayer[parent], &recursiveBuild{keys: keys, nodes: nodes, parent: parent},
			)
			write.Unlock()
		}
	})

	ptree.recursiveAdd(nextLayer, needRoot)
}

// Insert will add the provided keys to the tree.
func (ptree *ptree) Insert(keys ...common.Comparator) {
	ia := newInsertAction(keys)
	ptree.checkAndRun(ia)
	ia.completer.Wait()
}

// Get will retrieve a list of keys from the provided keys.
func (ptree *ptree) Get(keys ...common.Comparator) common.Comparators {
	ga := newGetAction(keys)
	ptree.checkAndRun(ga)
	ga.completer.Wait()
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
