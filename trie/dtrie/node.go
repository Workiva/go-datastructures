/*
Copyright (c) 2016, Theodore Butler
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package dtrie

import (
	"fmt"
	"sync"

	"github.com/Workiva/go-datastructures/bitarray"
)

type node struct {
	entries []Entry
	nodeMap bitarray.Bitmap32
	dataMap bitarray.Bitmap32
	level   uint8 // level starts at 0
}

func (n *node) KeyHash() uint32    { return 0 }
func (n *node) Key() interface{}   { return nil }
func (n *node) Value() interface{} { return nil }

func (n *node) String() string {
	return fmt.Sprint(n.entries)
}

type collisionNode struct {
	entries []Entry
}

func (n *collisionNode) KeyHash() uint32    { return 0 }
func (n *collisionNode) Key() interface{}   { return nil }
func (n *collisionNode) Value() interface{} { return nil }

func (n *collisionNode) String() string {
	return fmt.Sprintf("<COLLISIONS %v>%v", len(n.entries), n.entries)
}

// Entry defines anything held within the data structure
type Entry interface {
	KeyHash() uint32
	Key() interface{}
	Value() interface{}
}

func emptyNode(level uint8, capacity int) *node {
	return &node{entries: make([]Entry, capacity), level: level}
}

func insert(n *node, entry Entry) *node {
	index := uint(mask(entry.KeyHash(), n.level))
	newNode := n
	if newNode.level == 6 { // handle hash collisions on 6th level
		if newNode.entries[index] == nil {
			newNode.entries[index] = entry
			newNode.dataMap = newNode.dataMap.SetBit(index)
			return newNode
		}
		if newNode.dataMap.HasBit(index) {
			if newNode.entries[index].Key() == entry.Key() {
				newNode.entries[index] = entry
				return newNode
			}
			cNode := &collisionNode{entries: make([]Entry, 2)}
			cNode.entries[0] = newNode.entries[index]
			cNode.entries[1] = entry
			newNode.entries[index] = cNode
			newNode.dataMap = newNode.dataMap.ClearBit(index)
			return newNode
		}
		cNode := newNode.entries[index].(*collisionNode)
		cNode.entries = append(cNode.entries, entry)
		return newNode
	}
	if !newNode.dataMap.HasBit(index) && !newNode.nodeMap.HasBit(index) { // insert directly
		newNode.entries[index] = entry
		newNode.dataMap = newNode.dataMap.SetBit(index)
		return newNode
	}
	if newNode.nodeMap.HasBit(index) { // insert into sub-node
		newNode.entries[index] = insert(newNode.entries[index].(*node), entry)
		return newNode
	}
	if newNode.entries[index].Key() == entry.Key() {
		newNode.entries[index] = entry
		return newNode
	}
	// create new node with the new and existing entries
	var subNode *node
	if newNode.level == 5 { // only 2 bits left at level 6 (4 possible indices)
		subNode = emptyNode(newNode.level+1, 4)
	} else {
		subNode = emptyNode(newNode.level+1, 32)
	}
	subNode = insert(subNode, newNode.entries[index])
	subNode = insert(subNode, entry)
	newNode.dataMap = newNode.dataMap.ClearBit(index)
	newNode.nodeMap = newNode.nodeMap.SetBit(index)
	newNode.entries[index] = subNode
	return newNode
}

// returns nil if not found
func get(n *node, keyHash uint32, key interface{}) Entry {
	index := uint(mask(keyHash, n.level))
	if n.dataMap.HasBit(index) {
		return n.entries[index]
	}
	if n.nodeMap.HasBit(index) {
		return get(n.entries[index].(*node), keyHash, key)
	}
	if n.level == 6 { // get from collisionNode
		if n.entries[index] == nil {
			return nil
		}
		cNode := n.entries[index].(*collisionNode)
		for _, e := range cNode.entries {
			if e.Key() == key {
				return e
			}
		}
	}
	return nil
}

func remove(n *node, keyHash uint32, key interface{}) *node {
	index := uint(mask(keyHash, n.level))
	newNode := n
	if n.dataMap.HasBit(index) {
		newNode.entries[index] = nil
		newNode.dataMap = newNode.dataMap.ClearBit(index)
		return newNode
	}
	if n.nodeMap.HasBit(index) {
		subNode := newNode.entries[index].(*node)
		subNode = remove(subNode, keyHash, key)
		// compress if only 1 entry exists in sub-node
		if subNode.nodeMap.PopCount() == 0 && subNode.dataMap.PopCount() == 1 {
			var e Entry
			for i := uint(0); i < 32; i++ {
				if subNode.dataMap.HasBit(i) {
					e = subNode.entries[i]
					break
				}
			}
			newNode.entries[index] = e
			newNode.nodeMap = newNode.nodeMap.ClearBit(index)
			newNode.dataMap = newNode.dataMap.SetBit(index)
		}
		newNode.entries[index] = subNode
		return newNode
	}
	if n.level == 6 { // delete from collisionNode
		cNode := newNode.entries[index].(*collisionNode)
		for i, e := range cNode.entries {
			if e.Key() == key {
				cNode.entries = append(cNode.entries[:i], cNode.entries[i+1:]...)
				break
			}
		}
		// compress if only 1 entry exists in collisionNode
		if len(cNode.entries) == 1 {
			newNode.entries[index] = cNode.entries[0]
			newNode.dataMap = newNode.dataMap.SetBit(index)
		}
		return newNode
	}
	return n
}

func iterate(n *node, stop <-chan struct{}) <-chan Entry {
	out := make(chan Entry)
	go func() {
		defer close(out)
		pushEntries(n, stop, out)
	}()
	return out
}

func pushEntries(n *node, stop <-chan struct{}, out chan Entry) {
	var wg sync.WaitGroup
	for i, e := range n.entries {
		select {
		case <-stop:
			return
		default:
			index := uint(i)
			switch {
			case n.dataMap.HasBit(index):
				out <- e
			case n.nodeMap.HasBit(index):
				wg.Add(1)
				go func() {
					defer wg.Done()
					pushEntries(e.(*node), stop, out)
				}()
				wg.Wait()
			case n.level == 6 && e != nil:
				for _, ce := range n.entries[index].(*collisionNode).entries {
					select {
					case <-stop:
						return
					default:
						out <- ce
					}
				}
			}
		}
	}
}
