package plus

func keySearch(keys keys, key Key) int {
	low, high := 0, len(keys)-1
	var mid int
	for low <= high {
		mid = (high + low) / 2
		switch keys[mid].Compare(key) {
		case 1:
			low = mid + 1
		case -1:
			high = mid - 1
		case 0:
			return mid
		}
	}
	return low
}

type btree struct {
	root             node
	nodeSize, number uint64
}

func (tree *btree) insert(key Key) {
	if tree.root == nil {
		n := newLeafNode(tree.nodeSize)
		n.insert(tree, key)
		tree.number = 1
		return
	}

	result := tree.root.insert(tree, key)
	if result {
		tree.number++
	}

	if tree.root.needsSplit(tree.nodeSize) {
		println(`calling split here`)
		tree.root = split(tree, nil, tree.root)
	}
}

func (tree *btree) Insert(keys ...Key) {
	for _, key := range keys {
		tree.insert(key)
	}
}

func (tree *btree) Iterate(key Key) *iterator {
	if tree.root == nil {
		return nilIterator()
	}

	return tree.root.find(key)
}

func newBTree(nodeSize uint64) *btree {
	return &btree{
		nodeSize: nodeSize,
		root:     newLeafNode(nodeSize),
	}
}
