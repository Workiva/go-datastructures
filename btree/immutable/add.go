package btree

import (
	"runtime"
	"sort"
	"sync"

	terr "github.com/Workiva/go-datastructures/threadsafe/err"
)

func (t *Tr) AddItems(its ...*Item) ([]*Item, error) {
	if len(its) == 0 {
		return nil, nil
	}

	keys := make(Keys, 0, len(its))
	for _, item := range its {
		keys = append(keys, &Key{Value: item.Value, Payload: item.Payload})
	}

	overwrittens, err := t.add(keys)
	if err != nil {
		return nil, err
	}

	return overwrittens.toItems(), nil
}

func (t *Tr) add(keys Keys) (Keys, error) {
	if t.Root == nil {
		n := t.createRoot()
		t.Root = n.ID
		t.context.addNode(n)
	}

	nodes, err := t.determinePaths(keys)
	if err != nil {
		return nil, err
	}

	var overwrittens Keys

	var wg sync.WaitGroup
	wg.Add(len(nodes))
	var treeLock sync.Mutex
	localOverwrittens := make([]Keys, len(nodes))
	tree := make(map[string]*path, runtime.NumCPU())
	lerr := terr.New()

	i := 0
	for id, bundles := range nodes {
		go func(i int, id string, bundles []*nodeBundle) {
			defer wg.Done()
			if len(bundles) == 0 {
				return
			}

			n, err := t.contextOrCachedNode(ID(id), true)
			if err != nil {
				lerr.Set(err)
				return
			}

			if !t.context.nodeExists(n.ID) {
				n = n.copy()
				t.context.addNode(n)
			}

			overwrittens, err := insertLastDimension(t, n, bundles)

			if err != nil {
				lerr.Set(err)
				return
			}

			localOverwrittens[i] = overwrittens
			path := bundles[0].path
			treeLock.Lock()
			tree[string(n.ID)] = path
			treeLock.Unlock()
		}(i, id, bundles)
		i++
	}

	wg.Wait()

	if lerr.Get() != nil {
		return nil, lerr.Get()
	}

	t.walkupInsert(tree)

	for _, chunk := range localOverwrittens {
		overwrittens = append(overwrittens, chunk...)
	}

	t.Count += len(keys) - len(overwrittens)

	return overwrittens, nil
}

func (t *Tr) determinePaths(keys Keys) (map[string][]*nodeBundle, error) {
	chunks := splitKeys(keys, runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	chunkPaths := make([]map[interface{}]*nodeBundle, len(chunks))
	lerr := terr.New()

	for i := range chunks {
		go func(i int) {
			defer wg.Done()
			keys := chunks[i]
			if len(keys) == 0 {
				return
			}
			mp := make(map[interface{}]*nodeBundle, len(keys))
			for _, key := range keys {
				path, err := t.iterativeFind(
					key.Value, t.Root,
				)

				if err != nil {
					lerr.Set(err)
					return
				}
				mp[key.Value] = &nodeBundle{path: path, k: key}
			}
			chunkPaths[i] = mp
		}(i)
	}

	wg.Wait()

	if lerr.Get() != nil {
		return nil, lerr.Get()
	}

	nodes := make(map[string][]*nodeBundle, 10)
	for _, chunk := range chunkPaths {
		for _, pb := range chunk {
			nodes[string(pb.path.peek().n.ID)] = append(nodes[string(pb.path.pop().n.ID)], pb)
		}
	}

	return nodes, nil
}

func insertByMerge(comparator Comparator, n *Node, bundles []*nodeBundle) (Keys, error) {
	positions := make(map[interface{}]int, len(n.ChildValues))
	overwrittens := make(Keys, 0, 10)

	for i, value := range n.ChildValues {
		positions[value] = i
	}

	for _, bundle := range bundles {
		if i, ok := positions[bundle.k.Value]; ok {
			overwrittens = append(overwrittens, n.ChildKeys[i])
			n.ChildKeys[i] = bundle.k
		} else {
			n.ChildValues = append(n.ChildValues, bundle.k.Value)
			n.ChildKeys = append(n.ChildKeys, bundle.k)
		}
	}

	nsw := &nodeSortWrapper{
		values:     n.ChildValues,
		keys:       n.ChildKeys,
		comparator: comparator,
	}

	sort.Sort(nsw)

	for i := 0; i < len(nsw.values); i++ {
		if nsw.values[i] != nil {
			nsw.values = nsw.values[i:]
			nsw.keys = nsw.keys[i:]
			break
		}

		nsw.keys[i] = nil
	}

	n.ChildValues = nsw.values
	n.ChildKeys = nsw.keys
	return overwrittens, nil
}

func insertLastDimension(t *Tr, n *Node, bundles []*nodeBundle) (Keys, error) {
	if n.IsLeaf && len(bundles) >= n.lenValues()/16 { // Found through empirical testing, it appears that the memmoves are more sensitive when dealing with interface{}'s.
		return insertByMerge(t.config.Comparator, n, bundles)
	}

	overwrittens := make(Keys, 0, len(bundles))
	for _, bundle := range bundles {
		overwritten := n.insert(t.config.Comparator, bundle.k)
		if overwritten != nil {
			overwrittens = append(overwrittens, overwritten)
		}
	}

	return overwrittens, nil
}

func (t *Tr) iterativeSplit(n *Node) Keys {
	keys := make(Keys, 0, 10)
	for n.needsSplit(t.config.NodeWidth) {
		leftValue, leftNode := n.splitAt(t.config.NodeWidth / 2)
		t.context.addNode(leftNode)
		keys = append(keys, &Key{UUID: leftNode.ID, Value: leftValue})
	}

	return keys
}

// walkupInsert walks up nodes during the insertion process and adds
// any new keys due to splits.  Each layer of the tree can have insertions
// performed in parallel as splits are local changes.
func (t *Tr) walkupInsert(nodes map[string]*path) error {
	mapping := make(map[string]*Node, len(nodes))

	for len(nodes) > 0 {
		splitNodes := make(map[string]Keys)
		newNodes := make(map[string]*path)
		for id, path := range nodes {
			node := t.context.getNode(ID(id))

			parentPath := path.pop()
			if parentPath == nil {
				t.Root = node.ID
				continue
			}

			parent := parentPath.n
			newNode := mapping[string(parent.ID)]
			if newNode == nil {
				if !t.context.nodeExists(parent.ID) {
					cp := parent.copy()
					if string(t.Root) == string(parent.ID) {
						t.Root = cp.ID
					}

					t.context.addNode(cp)
					mapping[string(parent.ID)] = cp
					parent = cp
				} else {
					newNode = t.context.getNode(parent.ID)
					mapping[string(parent.ID)] = newNode
					parent = newNode
				}
			} else {
				parent = newNode
			}

			i := parentPath.i

			parent.replaceKeyAt(&Key{UUID: node.ID}, i)
			splitNodes[string(parent.ID)] = append(splitNodes[string(parent.ID)], t.iterativeSplit(node)...)
			newNodes[string(parent.ID)] = path
		}

		var wg sync.WaitGroup
		wg.Add(len(splitNodes))
		lerr := terr.New()

		for id, keys := range splitNodes {
			go func(id ID, keys Keys) {
				defer wg.Done()
				n, err := t.contextOrCachedNode(id, true)
				if err != nil {
					lerr.Set(err)
					return
				}
				for _, key := range keys {
					n.insert(t.config.Comparator, key)
				}
			}(ID(id), keys)
		}

		wg.Wait()

		if lerr.Get() != nil {
			return lerr.Get()
		}

		nodes = newNodes
	}

	n := t.context.getNode(t.Root)
	for n.needsSplit(t.config.NodeWidth) {
		root := newNode()
		t.Root = root.ID
		t.context.addNode(root)
		root.appendChild(&Key{UUID: n.ID})
		keys := t.iterativeSplit(n)
		for _, key := range keys {
			root.insert(t.config.Comparator, key)
		}
		n = root
	}

	return nil
}
