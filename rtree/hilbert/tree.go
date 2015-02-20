package hilbert

import (
	"log"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/Workiva/go-datastructures/rtree"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

type operation int

const (
	get operation = iota
	add
	remove
	apply
)

const multiThreadAt = 1000 // number of keys before we multithread lookups

type keyBundle struct {
	key         hilbert
	left, right rtree.Rectangle
}

type tree struct {
	root            *node
	_padding0       [8]uint64
	number          uint64
	_padding1       [8]uint64
	ary, bufferSize uint64
	actions         *queue.RingBuffer
	cache           []interface{}
	buffer0         [8]uint64
	disposed        uint64
	buffer1         [8]uint64
	running         uint64
}

func (tree *tree) checkAndRun(action action) {
	if tree.actions.Len() > 0 {
		if action != nil {
			tree.actions.Put(action)
		}
		if atomic.CompareAndSwapUint64(&tree.running, 0, 1) {
			var a interface{}
			var err error
			for tree.actions.Len() > 0 {
				a, err = tree.actions.Get()
				if err != nil {
					return
				}
				tree.cache = append(tree.cache, a)
				if uint64(len(tree.cache)) >= tree.bufferSize {
					break
				}
			}

			go tree.operationRunner(tree.cache, true)
		}
	} else if action != nil {
		if atomic.CompareAndSwapUint64(&tree.running, 0, 1) {
			switch action.operation() {
			case get:
				ga := action.(*getAction)
				result := tree.search(ga.lookup)
				ga.result = result
				action.complete()
				tree.reset()
			case add, remove:
				if len(action.keys()) > multiThreadAt {
					tree.operationRunner(interfaces{action}, true)
				} else {
					tree.operationRunner(interfaces{action}, false)
				}
			}
		} else {
			tree.actions.Put(action)
			tree.checkAndRun(nil)
		}
	}
}

func (tree *tree) init(bufferSize, ary uint64) {
	tree.bufferSize = bufferSize
	tree.ary = ary
	tree.cache = make([]interface{}, 0, bufferSize)
	tree.root = newNode(true, newKeys(ary), newNodes(ary))
	tree.root.mbr = &rectangle{}
	tree.actions = queue.NewRingBuffer(tree.bufferSize)
}

func (tree *tree) operationRunner(xns interfaces, threaded bool) {
	writeOperations, deleteOperations, toComplete := tree.fetchKeys(xns, threaded)
	tree.recursiveMutate(writeOperations, deleteOperations, false, threaded)
	for _, a := range toComplete {
		a.complete()
	}

	tree.reset()
}

func (tree *tree) fetchKeys(xns interfaces, inParallel bool) (map[*node][]*keyBundle, map[*node][]*keyBundle, actions) {
	if inParallel {
		tree.fetchKeysInParallel(xns)
	} else {
		tree.fetchKeysInSerial(xns)
	}

	writeOperations := make(map[*node][]*keyBundle)
	deleteOperations := make(map[*node][]*keyBundle)
	toComplete := make(actions, 0, len(xns)/2)
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add:
			for i, n := range action.nodes() {
				writeOperations[n] = append(writeOperations[n], &keyBundle{key: action.rects()[i].hilbert, left: action.rects()[i].rect})
			}
			toComplete = append(toComplete, action)
		case remove:
			for i, n := range action.nodes() {
				deleteOperations[n] = append(deleteOperations[n], &keyBundle{key: action.rects()[i].hilbert, left: action.rects()[i].rect})
			}
			toComplete = append(toComplete, action)
		case get:
			action.complete()
		}
	}

	return writeOperations, deleteOperations, toComplete
}

func (tree *tree) fetchKeysInSerial(xns interfaces) {
	for _, ifc := range xns {
		action := ifc.(action)
		switch action.operation() {
		case add, remove:
			for i, key := range action.rects() {
				n := getParent(tree.root, key.hilbert)
				action.addNode(int64(i), n)
			}
		case get:
			ga := action.(*getAction)
			rects := tree.search(ga.lookup)
			ga.result = rects
		}
	}
}

func (tree *tree) reset() {
	for i := range tree.cache {
		tree.cache[i] = nil
	}

	tree.cache = tree.cache[:0]
	atomic.StoreUint64(&tree.running, 0)
	tree.checkAndRun(nil)
}

func (tree *tree) fetchKeysInParallel(xns []interface{}) {
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

				switch action.operation() {
				case add, remove:
					n := getParent(tree.root, action.keys()[j])
					action.addNode(j, n)
				case get:
					ga := action.(*getAction)
					result := tree.search(ga.lookup)
					ga.result = result
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (tree *tree) splitNode(n, parent *node, nodes *[]*node, keys *hilberts) {
	if !n.needsSplit(tree.ary) {
		return
	}

	length := n.keys.len()
	splitAt := tree.ary - 1

	for i := splitAt; i < length; i += splitAt {
		offset := length - i
		k, left, right := n.split(offset, tree.ary)
		left.right = right
		*keys = append(*keys, k)
		*nodes = append(*nodes, left, right)
		left.parent = parent
		right.parent = parent
	}
}

func (tree *tree) applyNode(n *node, adds, deletes []*keyBundle) {
	for _, kb := range deletes {
		if n.keys.len() == 0 {
			break
		}

		deleted := n.keys.delete(kb.key)
		if deleted != -1 {
			atomic.AddUint64(&tree.number, ^uint64(0))
		}
	}

	for _, kb := range adds {
		old := n.insert(kb)
		if n.isLeaf && old == nil {
			atomic.AddUint64(&tree.number, 1)
		}
	}
}

func (tree *tree) recursiveMutate(adds, deletes map[*node][]*keyBundle, setRoot, inParallel bool) {
	if len(adds) == 0 && len(deletes) == 0 {
		return
	}

	if setRoot && len(adds) > 1 {
		panic(`SHOULD ONLY HAVE ONE ROOT`)
	}

	ifs := make(interfaces, 0, len(adds))
	for n := range adds {
		if n.parent == nil {
			setRoot = true
		}
		ifs = append(ifs, n)
	}

	for n := range deletes {
		if n.parent == nil {
			setRoot = true
		}

		if _, ok := adds[n]; !ok {
			ifs = append(ifs, n)
		}
	}

	var dummyRoot *node
	if setRoot {
		dummyRoot = &node{
			keys:  newKeys(tree.ary),
			nodes: newNodes(tree.ary),
			mbr:   &rectangle{},
		}
	}

	var write sync.Mutex
	nextLayerWrite := make(map[*node][]*keyBundle)
	nextLayerDelete := make(map[*node][]*keyBundle)

	var mutate func(interfaces, func(interface{}))
	if inParallel {
		mutate = executeInterfacesInParallel
	} else {
		mutate = executeInterfacesInSerial
	}

	mutate(ifs, func(ifc interface{}) {
		n := ifc.(*node)
		adds := adds[n]
		deletes := deletes[n]

		if len(adds) == 0 && len(deletes) == 0 {
			return
		}

		if setRoot {
			tree.root = n
		}

		parent := n.parent
		if parent == nil {
			parent = dummyRoot
			setRoot = true
		}

		tree.applyNode(n, adds, deletes)

		if n.needsSplit(tree.ary) {
			keys := make(hilberts, 0, n.keys.len())
			nodes := make([]*node, 0, n.nodes.len())
			tree.splitNode(n, parent, &nodes, &keys)
			write.Lock()
			for i, k := range keys {
				nextLayerWrite[parent] = append(nextLayerWrite[parent], &keyBundle{key: k, left: nodes[i*2], right: nodes[i*2+1]})
			}
			write.Unlock()
		}
	})

	tree.recursiveMutate(nextLayerWrite, nextLayerDelete, setRoot, inParallel)
}

// Insert will add the provided keys to the tree.
func (tree *tree) Insert(rects ...rtree.Rectangle) {
	ia := newInsertAction(rects)
	tree.checkAndRun(ia)
	ia.completer.Wait()
}

// Delete will remove the provided keys from the tree.  If no
// matching key is found, this is a no-op.
func (tree *tree) Delete(rects ...rtree.Rectangle) {
	ra := newRemoveAction(rects)
	tree.checkAndRun(ra)
	ra.completer.Wait()
}

func (tree *tree) search(r *rectangle) rtree.Rectangles {
	if tree.root == nil {
		return rtree.Rectangles{}
	}

	result := make(rtree.Rectangles, 0, 10)
	whs := tree.root.searchRects(r)
	for len(whs) > 0 {
		wh := whs[0]
		if n, ok := wh.(*node); ok {
			whs = append(whs, n.searchRects(r)...)
		} else {
			result = append(result, wh)
		}
		whs = whs[1:]
	}

	return result
}

// Search will return a list of rectangles that intersect the provided
// rectangle.
func (tree *tree) Search(rect rtree.Rectangle) rtree.Rectangles {
	ga := newGetAction(rect)
	tree.checkAndRun(ga)
	ga.completer.Wait()
	return ga.result
}

// Len returns the number of items in the tree.
func (tree *tree) Len() uint64 {
	return atomic.LoadUint64(&tree.number)
}

// Dispose will clean up any resources used by this tree.  This
// must be called to prevent a memory leak.
func (tree *tree) Dispose() {
	tree.actions.Dispose()
	atomic.StoreUint64(&tree.disposed, 1)
}

func newTree(bufferSize, ary uint64) *tree {
	tree := &tree{}
	tree.init(bufferSize, ary)
	return tree
}
