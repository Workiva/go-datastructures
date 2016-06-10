package btree

import "bytes"

func (t *Tr) DeleteItems(values ...interface{}) ([]*Item, error) {
	if len(values) == 0 {
		return nil, nil
	}

	keys := make(Keys, 0, len(values))
	err := t.Apply(func(item *Item) {
		keys = append(keys, &Key{Value: item.Value, Payload: item.Payload})
	}, values...)

	err = t.delete(keys)
	if err != nil {
		return nil, err
	}

	t.Count -= len(keys)

	return keys.toItems(), nil
}

func (t *Tr) delete(keys Keys) error {
	if len(keys) == 0 {
		return nil
	}

	toDelete := make([]*Key, 0, len(keys))

	for i := 0; i < len(keys); {
		key := keys[i]
		mapping := make(map[string]*Node, 10)
		path, err := t.iterativeFind(key.Value, t.Root)
		if err != nil {
			return err
		}

		pb := path.peek()
		node := pb.n
		isRoot := bytes.Compare(node.ID, t.Root) == 0
		if !t.context.nodeExists(node.ID) {
			cp := node.copy()
			t.context.addNode(cp)
			mapping[string(node.ID)] = cp
			node = cp
		}
		base := node

		toDelete = append(toDelete, key)
		for j := i + 1; j <= len(keys); j++ {
			i = j
			if j == len(keys) {
				break
			}
			neighbor := keys[j]
			if t.config.Comparator(neighbor.Value, node.lastValue()) <= 0 {
				toDelete = append(toDelete, neighbor)
			} else {
				break
			}
		}

		if len(toDelete) > len(node.ChildValues)/4 {
			node.multiDelete(t.config.Comparator, toDelete...)
		} else {
			for _, k := range toDelete {
				node.delete(t.config.Comparator, k)
			}
		}

		toDelete = toDelete[:0]
		if isRoot {
			t.Root = node.ID
			continue
		}

		for pb.prev != nil {
			parentBundle := pb.prev
			parent := parentBundle.n
			isRoot := bytes.Compare(parent.ID, t.Root) == 0
			if !t.context.nodeExists(parent.ID) {
				cp := parent.copy()
				t.context.addNode(cp)
				mapping[string(parent.ID)] = cp
				parent = cp
			} else {
				mapping[string(parent.ID)] = parent
			}

			if isRoot {
				t.Root = parent.ID
			}

			i := pb.prev.i
			parent.replaceKeyAt(&Key{UUID: node.ID}, i)
			node = parent
			pb = pb.prev
		}

		path.pop()
		err = t.walkupDelete(key, base, path, mapping)
		if err != nil {
			return err
		}
	}

	n := t.context.getNode(t.Root)
	if n.lenValues() == 0 {
		t.Root = nil
	}

	return nil
}

// walkupDelete is similar to walkupInsert but is only done one at a time.
// This is because deletes can cause borrowing or merging with neighbors which makes
// the changes non-local.
// TODO: come up with a good way to parallelize this.
func (t *Tr) walkupDelete(key *Key, node *Node, path *path, mapping map[string]*Node) error {
	needsMerged := t.config.NodeWidth / 2
	if needsMerged < 1 {
		needsMerged = 1
	}
	if node.lenValues() >= needsMerged {
		return nil
	}

	if string(node.ID) == string(t.Root) {
		if node.lenKeys() == 1 {
			id := node.keyAt(0)
			t.Root = id.UUID
		}

		return nil
	}

	var getSibling = func(parent *Node, i int) (*Node, error) {
		key := parent.keyAt(i)
		n, err := t.contextOrCachedNode(key.UUID, true)
		if err != nil {
			return nil, err
		}

		if !t.context.nodeExists(n.ID) {
			cp := t.copyNode(n)
			mapping[string(n.ID)] = cp
			parent.replaceKeyAt(&Key{UUID: cp.ID}, i)
			n = cp
		}

		return n, nil
	}

	parentBundle := path.pop()
	parent := mapping[string(parentBundle.n.ID)]

	_, i := parent.searchKey(t.config.Comparator, key.Value)
	siblingPosition := i
	if i == parent.lenValues() {
		siblingPosition--
	} else {
		siblingPosition++
	}

	sibling, err := getSibling(parent, siblingPosition)
	if err != nil {
		return err
	}

	prepend := false
	// thing are just easier if we make this swap so we can grok
	// left to right always assuming node is on the left and sibling
	// is on the right
	if siblingPosition < i {
		node, sibling = sibling, node
		prepend = true
	}

	// first case, can we just borrow?  if so, simply shift values from one node
	// to the other.  Once done, replace the parent value with the middle value
	// shifted and return.
	if (sibling.lenValues()+node.lenValues())/2 >= needsMerged {
		if i == parent.lenValues() {
			i--
		}

		var key *Key
		var value interface{}
		for node.lenValues() < needsMerged || sibling.lenValues() < needsMerged {
			if prepend {
				correctedValue, key := node.popValue(), node.popKey()
				if node.IsLeaf {
					sibling.prependValue(correctedValue)
					sibling.prependKey(key)
					parent.replaceValueAt(i, node.lastValue())
				} else {
					parentValue := parent.valueAt(i)
					sibling.prependKey(key)
					sibling.prependValue(parentValue)
					parent.replaceValueAt(i, correctedValue)
				}
			} else {
				value, key = sibling.popFirstValue(), sibling.popFirstKey()
				correctedValue := value
				if !node.IsLeaf {
					correctedValue = parent.valueAt(i)
				}
				node.appendValue(correctedValue)
				node.appendChild(key)
				parent.replaceValueAt(i, value)
			}
		}

		return nil
	}

	// the harder case, we need to merge with sibling, pull a value down
	// from the parent, and recurse on this function

	// easier case, merge the nodes and delete value and child from parent
	if node.IsLeaf {
		node.append(sibling)
		if prepend {
			parent.deleteKeyAt(i)
		} else {
			parent.deleteKeyAt(i + 1)
		}

		if i == parent.lenValues() {
			i--
		}

		parent.deleteValueAt(i)
		return t.walkupDelete(key, parent, path, mapping)
	}

	// harder case, need to pull a value down from the parent, insert
	// value into the left node, append the nodes, and then delete
	// the value from the parent

	valueIndex := i
	if i == parent.lenValues() {
		valueIndex--
	}

	parentValue := parent.valueAt(valueIndex)
	node.appendValue(parentValue)
	node.append(sibling)
	parent.deleteKeyAt(i)
	parent.deleteValueAt(valueIndex)
	parent.replaceKeyAt(&Key{UUID: node.ID}, valueIndex)
	return t.walkupDelete(key, parent, path, mapping)
}
