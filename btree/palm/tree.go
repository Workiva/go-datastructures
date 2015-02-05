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
}

func (ptree *ptree) runOperations() {
	ptree.lock.Lock()
	toPerform := ptree.pending
	ptree.pending = &pending{}
	ptree.lock.Unlock()

	q := queue.New(int64(toPerform.number))
	var key Key
	var i uint64
	for _, ab := range toPerform.bundles {
		for {
			key, i = ab.getKey()
			if key == nil {
				break
			}

			q.Put(&actionBundle{key: key, index: i, action: ab})
		}
	}

	readOperations := make(bundleMap, q.Len())
	writeOperations := make(bundleMap, q.Len())

	queue.ExecuteInParallel(q, func(ifc interface{}) {
		ab := ifc.(*actionBundle)

		node := getParent(ptree.root, ab.key)
		ab.node = node
		switch ab.action.operation() {
		case get:
			ptree.read.Lock()
			readOperations[node] = append(readOperations[node], ab)
			ptree.read.Unlock()
		case add, remove:
			ptree.write.Lock()
			writeOperations[node] = append(writeOperations[node], ab)
			ptree.write.Unlock()
		}
	})

	ptree.runReads(readOperations)
	ptree.runAdds(writeOperations)
}

func (ptree *ptree) runReads(readOperations bundleMap) {
	q := queue.New(int64(len(readOperations)))

	for _, abs := range readOperations {
		for _, ab := range abs {
			q.Put(ab)
		}
	}

	queue.ExecuteInParallel(q, func(ifc interface{}) {
		ab := ifc.(*actionBundle)
		if ab.node == nil {
			ab.action.addResult(ab.index, nil)
			return
		}

		result := ab.node.search(ab.key)
		if result == len(ab.node.keys) {
			ab.action.addResult(ab.index, nil)
			return
		}

		if ab.node.keys[result].Compare(ab.key) == 0 {
			ab.action.addResult(ab.index, ab.node.keys[result])
			return
		}

		ab.action.addResult(ab.index, nil)
	})
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

func (ptree *ptree) runAdds(addOperations bundleMap) {
	q := queue.New(int64(len(addOperations)))

	for _, abs := range addOperations {
		q.Put(abs)
	}

	nextLayer := make(map[*node][]*recursiveBuild)
	dummyRoot := &node{} // constructed in case we need it
	var needRoot uint64
	queue.ExecuteInParallel(q, func(ifc interface{}) {
		abs := ifc.(actionBundles)

		if len(abs) == 0 {
			return
		}

		n := abs[0].node
		parent := n.parent
		if parent == nil {
			parent = dummyRoot
			atomic.AddUint64(&needRoot, 1)
		}

		for _, ab := range abs {
			oldKey := n.keys.insert(ab.key)
			ab.action.addResult(ab.index, oldKey)
			if oldKey == nil {
				atomic.AddUint64(&ptree.number, 1)
			}
		}

		//log.Printf(`N BEFORE SPLIT: %+v`, n)
		if n.needsSplit(ptree.ary) {
			keys := make(Keys, 0, len(n.keys))
			nodes := make(nodes, 0, len(n.nodes))
			ptree.recursiveSplit(n, parent, nil, &nodes, &keys)
			//log.Printf(`AFTER SPLIT: %+v, NODES: %+v, KEYS: %+v`, n, nodes, keys)
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

	go ptree.runOperations()
	result := <-ia.completer
	return result
}

func (ptree *ptree) Get(keys ...Key) Keys {
	ga := newGetAction(keys)
	ptree.lock.Lock()
	ptree.pending.bundles = append(ptree.pending.bundles, ga)
	ptree.pending.number += uint64(len(keys))
	ptree.lock.Unlock()

	go ptree.runOperations()
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
	return &ptree{
		root:    newNode(true, make(Keys, 0, ary), make(nodes, 0, ary+1)),
		ary:     ary,
		pending: &pending{},
	}
}
