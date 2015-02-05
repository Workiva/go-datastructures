package palm

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/queue"
)

func init() {
	log.Printf(`I HATE THIS.`)
}

type operation int

const (
	get operation = iota
	add
	remove
)

type recursiveBuild struct {
	keys   Keys
	nodes  nodes
	parent *node
}

type pending struct {
	bundles actions
	number  uint64
}

type bundleMap map[*node]actionBundles

type ptree struct {
	root                             *node
	ary, number, threads, bufferSize uint64
	pending                          *pending
	lock, read, write                sync.Mutex
	waiter                           *queue.Queue
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

func (ptree *ptree) runOperations() {
	println(`RUNNING OPERATIONS`)
	ptree.lock.Lock()
	toPerform := ptree.pending
	ptree.pending = &pending{}
	ptree.lock.Unlock()

	log.Printf(`TO PERFORM: %+v`, toPerform)
	if toPerform.number == 0 {
		return
	}

	var key Key
	var i uint64

	writeOperations := make(map[*node]Keys, toPerform.number/2)

	for _, action := range toPerform.bundles {
		for {
			key, i = action.getKey()
			if key == nil {
				break
			}

			n := getParent(ptree.root, key)
			ab := &actionBundle{key: key, index: i, action: action, node: n}
			switch ab.action.operation() {
			case get:
				if n == nil {
					ab.action.addResult(i, nil)
				}
				index := n.keys.search(key)
				if index < len(n.keys) && n.keys[index].Compare(key) == 0 {
					ab.action.addResult(i, n.keys[index])
				} else {
					ab.action.addResult(i, nil)
				}
			case add, remove:
				writeOperations[n] = append(writeOperations[n], key)
			}
		}
	}

	ptree.runAdds(writeOperations)
	for _, action := range toPerform.bundles {
		if action.operation() == add || action.operation() == remove {
			action.complete()
		}
	}
}

func (ptree *ptree) recursiveSplit(n, parent, left *node, nodes *nodes, keys *Keys) {
	if !n.needsSplit(ptree.ary) {
		return
	}

	//log.Printf(`N: %+v, parent: %+v, left: %+v, nodes: %+v, keys: %+v`, n, parent, left, nodes, keys)
	key, l, r := n.split()
	if left != nil {
		left.right = l
	}
	//log.Printf(`KEY: %+v, L: %+v: R: %+v`, key, l, r)
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
	dummyRoot := &node{}
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
				//log.Printf(`LOOP N: %+v`, n)
				if len(n.keys) == 0 {
					n.keys.insert(k)
					n.nodes.push(rb.nodes[i*2])
					n.nodes.push(rb.nodes[i*2+1])
					continue
				}

				index := n.search(k)
				n.keys.insertAt(k, index)
				n.nodes[index] = rb.nodes[i*2]
				n.nodes.insertAt(rb.nodes[i*2+1], index+1)
			}
		}

		if n.needsSplit(ptree.ary) {
			keys := make(Keys, 0, len(n.keys))
			nodes := make(nodes, 0, len(n.nodes))
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
	dummyRoot := &node{} // constructed in case we need it
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
			keys := make(Keys, 0, len(n.keys))
			nodes := make(nodes, 0, len(n.nodes))
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

func (ptree *ptree) Insert(keys ...Key) Keys {
	ia := newInsertAction(keys)
	ptree.lock.Lock()
	ptree.pending.bundles = append(ptree.pending.bundles, ia)
	ptree.pending.number += uint64(len(keys))
	ptree.lock.Unlock()

	ptree.waiter.Put(true)
	result := <-ia.completer
	return result
}

func (ptree *ptree) Get(keys ...Key) Keys {
	ga := newGetAction(keys)
	ptree.lock.Lock()
	ptree.pending.bundles = append(ptree.pending.bundles, ga)
	ptree.pending.number += uint64(len(keys))
	ptree.lock.Unlock()

	ptree.waiter.Put(true)
	return <-ga.completer
}

func (ptree *ptree) Len() uint64 {
	return atomic.LoadUint64(&ptree.number)
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
		root:    newNode(true, make(Keys, 0, ary), make(nodes, 0, ary+1)),
		ary:     ary,
		pending: &pending{},
		waiter:  queue.New(10),
	}
	go ptree.listen()
	return ptree
}
