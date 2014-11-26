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
	tree.root.insert(tree, key)
}

func (tree *btree) Insert(keys ...Key) {

}

func newBTree(nodeSize uint64) *btree {
	return &btree{
		nodeSize: nodeSize,
		root:     newLeafNode(nodeSize),
	}
}
