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

//go:generate msgp -tests=false -io=false

package btree

import (
	"sort"

	"github.com/satori/go.uuid"
)

// ID exists because i'm tired of writing []byte
type ID []byte

// Key a convenience struct that holds both an id and a value.  Internally,
// this is how we reference items in nodes but consumers interface with
// the tree using row/col/id.
type Key struct {
	UUID    ID          `msg:"u"`
	Value   interface{} `msg:"v"`
	Payload []byte      `msg:"p"`
}

// ID returns the  unique identifier.
func (k Key) ID() []byte {
	return k.UUID[:16] // to maintain backwards compatability
}

func (k Key) ToItem() *Item {
	return &Item{
		Value:   k.Value,
		Payload: k.Payload,
	}
}

type Keys []*Key

func (k Keys) toItems() items {
	items := make(items, 0, len(k))
	for _, key := range k {
		items = append(items, key.ToItem())
	}

	return items
}

func (k Keys) sort(comparator Comparator) Keys {
	return (&keySortWrapper{comparator, k}).sort()
}

type keySortWrapper struct {
	comparator Comparator
	keys       Keys
}

func (sw *keySortWrapper) Len() int {
	return len(sw.keys)
}

func (sw *keySortWrapper) Swap(i, j int) {
	sw.keys[i], sw.keys[j] = sw.keys[j], sw.keys[i]
}

func (sw *keySortWrapper) Less(i, j int) bool {
	return sw.comparator(sw.keys[i].Value, sw.keys[j].Value) < 0
}

func (sw *keySortWrapper) sort() Keys {
	sort.Sort(sw)
	return sw.keys
}

func splitKeys(keys Keys, numParts int) []Keys {
	parts := make([]Keys, numParts)
	for i := int64(0); i < int64(numParts); i++ {
		parts[i] = keys[i*int64(len(keys))/int64(numParts) : (i+1)*int64(len(keys))/int64(numParts)]
	}
	return parts
}

// Node represents either a leaf node or an internal node.  These
// are the value containers.  This is exported because code generation
// requires it.  Only exported fields are required to be persisted.  We
// use msgpack for optimal performance.
type Node struct {
	// ID is the unique UUID that addresses this singular node.
	ID ID `msg:"id"`
	// IsLeaf is a bool indicating if this is a leaf node as opposed
	// to an internal node.  The primary difference between these nodes
	// is that leaf nodes have an equal number of values and IDs while
	// internal nodes have n+1 ids.
	IsLeaf bool `msg:"il"`
	// ChildValues is only a temporary field that is used to house all
	// values for serialization purposes.
	ChildValues []interface{} `msg:"cv"`
	// ChildKeys is similar to child values but holds the IDs of children.
	ChildKeys Keys `msg:"ck"`
}

// copy makes a deep copy of this node.  Required before any mutation.
func (n *Node) copy() *Node {
	cpValues := make([]interface{}, len(n.ChildValues))
	copy(cpValues, n.ChildValues)
	cpKeys := make(Keys, len(n.ChildKeys))
	copy(cpKeys, n.ChildKeys)

	return &Node{
		ID:          uuid.NewV4().Bytes(),
		IsLeaf:      n.IsLeaf,
		ChildValues: cpValues,
		ChildKeys:   cpKeys,
	}
}

// searchKey returns the key associated with the provided value.  If the
// provided value is greater than the highest value in this node and this
// node is an internal node, this method returns the last ID and an index
// equal to lenValues.
func (n *Node) searchKey(comparator Comparator, value interface{}) (*Key, int) {
	i := n.search(comparator, value)

	if n.IsLeaf && i == len(n.ChildValues) { // not found
		return nil, i
	}

	if n.IsLeaf { // equal number of ids and values
		return n.ChildKeys[i], i
	}

	if i == len(n.ChildValues) { // we need to go to the farthest node to the write
		return n.ChildKeys[len(n.ChildKeys)-1], i
	}

	return n.ChildKeys[i], i
}

// insert adds the provided key to this node and returns any ID that has
// been overwritten.  This method should only be called on leaf nodes.
func (n *Node) insert(comparator Comparator, key *Key) *Key {
	var overwrittenKey *Key
	i := n.search(comparator, key.Value)
	if i == len(n.ChildValues) {
		n.ChildValues = append(n.ChildValues, key.Value)
	} else {
		if n.ChildValues[i] == key.Value {
			overwrittenKey = n.ChildKeys[i]
			n.ChildKeys[i] = key
			return overwrittenKey
		} else {
			n.ChildValues = append(n.ChildValues, 0)
			copy(n.ChildValues[i+1:], n.ChildValues[i:])
			n.ChildValues[i] = key.Value
		}
	}

	if n.IsLeaf && i == len(n.ChildKeys) {
		n.ChildKeys = append(n.ChildKeys, key)
	} else {
		n.ChildKeys = append(n.ChildKeys, nil)
		copy(n.ChildKeys[i+1:], n.ChildKeys[i:])
		n.ChildKeys[i] = key
	}

	return overwrittenKey
}

// delete removes the provided key from the node and returns any key that
// was deleted.  Returns nil of the key could not be found.
func (n *Node) delete(comparator Comparator, key *Key) *Key {
	i := n.search(comparator, key.Value)
	if i == len(n.ChildValues) {
		return nil
	}

	n.deleteValueAt(i)
	n.deleteKeyAt(i)

	return key
}

func (n *Node) multiDelete(comparator Comparator, keys ...*Key) {
	indices := make([]int, 0, len(keys))
	for _, k := range keys {
		i := n.search(comparator, k.Value)
		if i < len(n.ChildValues) {
			indices = append(indices, i)
		}
	}

	for _, i := range indices {
		n.ChildValues[i] = nil
		n.ChildKeys[i] = nil
	}

	if len(indices) == len(n.ChildValues) {
		n.ChildKeys = n.ChildKeys[:0]
		n.ChildValues = n.ChildValues[:0]
		return
	}

	// get the indices in the correct order for the next stage
	// which is removing the nils
	sort.Ints(indices)

	// iterate through the list moving all values up to overwrite the
	// nils and place all nils at the "back"
	for i, j := range indices {
		index := j - i // correct for previous copies
		copy(n.ChildValues[index:], n.ChildValues[index+1:])
		copy(n.ChildKeys[index:], n.ChildKeys[index+1:])
	}

	n.ChildValues = n.ChildValues[:len(n.ChildValues)-len(indices)]
	n.ChildKeys = n.ChildKeys[:len(n.ChildKeys)-len(indices)]
}

// replaceKeyAt replaces the key at index i with the provided id.  This does
// not do any bounds checking.
func (n *Node) replaceKeyAt(key *Key, i int) {
	n.ChildKeys[i] = key
}

// flatten returns a flattened list of values and IDs.  Useful for serialization.
func (n *Node) flatten() ([]interface{}, Keys) {
	return n.ChildValues, n.ChildKeys
}

// iter returns an iterator that will iterate through the provided Morton
// numbers as they exist in this node.
func (n *Node) iter(comparator Comparator, start, stop interface{}) iterator {
	pointer := n.search(comparator, start)
	pointer--
	return &sliceIterator{
		stop:       stop,
		n:          n,
		pointer:    pointer,
		comparator: comparator,
	}
}

func (n *Node) valueAt(i int) interface{} {
	return n.ChildValues[i]
}

func (n *Node) keyAt(i int) *Key {
	return n.ChildKeys[i]
}

func (n *Node) needsSplit(max int) bool {
	return n.lenValues() > max
}

func (n *Node) lastValue() interface{} {
	return n.ChildValues[len(n.ChildValues)-1]
}

func (n *Node) firstValue() interface{} {
	return n.ChildValues[0]
}

func (n *Node) append(other *Node) {
	n.ChildValues = append(n.ChildValues, other.ChildValues...)
	n.ChildKeys = append(n.ChildKeys, other.ChildKeys...)
}

func (n *Node) replaceValueAt(i int, value interface{}) {
	n.ChildValues[i] = value
}

func (n *Node) deleteValueAt(i int) {
	copy(n.ChildValues[i:], n.ChildValues[i+1:])
	n.ChildValues[len(n.ChildValues)-1] = 0 // or the zero value of T
	n.ChildValues = n.ChildValues[:len(n.ChildValues)-1]
}

func (n *Node) deleteKeyAt(i int) {
	copy(n.ChildKeys[i:], n.ChildKeys[i+1:])
	n.ChildKeys[len(n.ChildKeys)-1] = nil // or the zero value of T
	n.ChildKeys = n.ChildKeys[:len(n.ChildKeys)-1]
}

func (n *Node) splitLeafAt(i int) (interface{}, *Node) {
	left := newNode()
	left.IsLeaf = n.IsLeaf
	left.ID = uuid.NewV4().Bytes()

	value := n.ChildValues[i]
	leftValues := make([]interface{}, i+1)
	copy(leftValues, n.ChildValues[:i+1])
	n.ChildValues = n.ChildValues[i+1:]
	leftKeys := make(Keys, i+1)
	copy(leftKeys, n.ChildKeys[:i+1])
	for j := 0; j <= i; j++ {
		n.ChildKeys[j] = nil
	}
	n.ChildKeys = n.ChildKeys[i+1:]
	left.ChildValues = leftValues
	left.ChildKeys = leftKeys
	return value, left
}

// splitInternalAt is a method that generates a new set of children
// for an internal node and returns the new set and the value that
// separates them.
func (n *Node) splitInternalAt(i int) (interface{}, *Node) {
	left := newNode()
	left.IsLeaf = n.IsLeaf
	left.ID = uuid.NewV4().Bytes()
	value := n.ChildValues[i]
	leftValues := make([]interface{}, i)
	copy(leftValues, n.ChildValues[:i])
	n.ChildValues = n.ChildValues[i+1:]
	leftKeys := make(Keys, i+1)
	copy(leftKeys, n.ChildKeys[:i+1])
	for j := 0; j <= i; j++ {
		n.ChildKeys[j] = nil
	}
	n.ChildKeys = n.ChildKeys[i+1:]
	left.ChildKeys = leftKeys
	left.ChildValues = leftValues
	return value, left
}

// splitAt breaks this node into two parts and conceptually
// returns the left part
func (n *Node) splitAt(i int) (interface{}, *Node) {
	if n.IsLeaf {
		return n.splitLeafAt(i)
	}

	return n.splitInternalAt(i)
}

func (n *Node) lenKeys() int {
	return len(n.ChildKeys)
}

func (n *Node) lenValues() int {
	return len(n.ChildValues)
}

func (n *Node) appendChild(key *Key) {
	n.ChildKeys = append(n.ChildKeys, key)
}

func (n *Node) appendValue(value interface{}) {
	n.ChildValues = append(n.ChildValues, value)
}

func (n *Node) popFirstKey() *Key {
	key := n.ChildKeys[0]
	n.deleteKeyAt(0)
	return key
}

func (n *Node) popFirstValue() interface{} {
	value := n.ChildValues[0]
	n.deleteValueAt(0)
	return value
}

func (n *Node) popKey() *Key {
	key := n.ChildKeys[len(n.ChildKeys)-1]
	n.deleteKeyAt(len(n.ChildKeys) - 1)
	return key
}

func (n *Node) popValue() interface{} {
	value := n.ChildValues[len(n.ChildValues)-1]
	n.deleteValueAt(len(n.ChildValues) - 1)
	return value
}

func (n *Node) prependKey(key *Key) {
	n.ChildKeys = append(n.ChildKeys, nil)
	copy(n.ChildKeys[1:], n.ChildKeys)
	n.ChildKeys[0] = key
}

func (n *Node) prependValue(value interface{}) {
	n.ChildValues = append(n.ChildValues, nil)
	copy(n.ChildValues[1:], n.ChildValues)
	n.ChildValues[0] = value
}

func (n *Node) search(comparator Comparator, value interface{}) int {
	return sort.Search(len(n.ChildValues), func(i int) bool {
		return comparator(n.ChildValues[i], value) >= 0
	})
}

// nodeFromBytes returns a new node struct deserialized from the provided
// bytes.  An error is returned for any deserialization errors.
func nodeFromBytes(t *Tr, data []byte) (*Node, error) {
	n := &Node{}
	_, err := n.UnmarshalMsg(data)
	if err != nil {
		panic(err)
		return nil, err
	}

	return n, nil
}

// newNode returns a node with a random id and empty values and children.
// IsLeaf is false by default.
func newNode() *Node {
	return &Node{
		ID: uuid.NewV4().Bytes(),
	}
}

type sliceIterator struct {
	stop       interface{}
	n          *Node
	pointer    int
	comparator Comparator
}

func (s *sliceIterator) next() bool {
	s.pointer++
	if s.n.IsLeaf {
		return s.pointer < len(s.n.ChildValues) && s.comparator(s.stop, s.n.ChildValues[s.pointer]) >= 0
	} else {
		if s.pointer >= len(s.n.ChildKeys) {
			return false
		}
		if s.pointer == len(s.n.ChildValues) {
			return true
		}

		if s.comparator(s.stop, s.n.ChildValues[s.pointer]) < 0 {
			return false
		}
	}

	return true
}

func (s *sliceIterator) value() (*Key, int) {
	return s.n.ChildKeys[s.pointer], s.pointer
}

type iterator interface {
	next() bool
	value() (*Key, int)
}

type nodeBundle struct {
	path *path
	k    *Key
}

type nodeSortWrapper struct {
	values     []interface{}
	keys       Keys
	comparator Comparator
}

func (n *nodeSortWrapper) Len() int {
	return len(n.values)
}

func (n *nodeSortWrapper) Swap(i, j int) {
	n.values[i], n.values[j] = n.values[j], n.values[i]
	n.keys[i], n.keys[j] = n.keys[j], n.keys[i]
}

func (n *nodeSortWrapper) Less(i, j int) bool {
	return n.comparator(n.values[i], n.values[j]) < 0
}

func splitValues(values []interface{}, numParts int) [][]interface{} {
	parts := make([][]interface{}, numParts)
	for i := int64(0); i < int64(numParts); i++ {
		parts[i] = values[i*int64(len(values))/int64(numParts) : (i+1)*int64(len(values))/int64(numParts)]
	}
	return parts
}
