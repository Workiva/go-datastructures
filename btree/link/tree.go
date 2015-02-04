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

/*
This is a b-link tree in progress from the following paper:
http://www.csd.uoc.gr/~hy460/pdf/p650-lehman.pdf

This is still a work in progress and the CRUD methods on the tree
need to be parallelized.  Until this is complete, there is no
constructor method for this package.

Time complexities:
Space: O(n)
Search: O(log n)
Insert: O(log n)
Delete: O(log n)

Current benchmarks with 16 ary:
BenchmarkSimpleAdd-8	 	1000000	      1455 ns/op
BenchmarkGet-8	 			2000000	       704 ns/op

B-link was chosen after examining this paper:
http://www.vldb.org/journal/VLDBJ2/P361.pdf
*/

package link

import (
	"log"
	"sync"
	"sync/atomic"
)

// numberOfItemsBeforeMultithread defines the number of items that have
// to be called with a method before we multithread.
const numberOfItemsBeforeMultithread = 10

type blink struct {
	root                     *node
	lock                     sync.RWMutex
	number, ary, numRoutines uint64
}

func (blink *blink) insert(key Key, stack *nodes) Key {
	var parent *node
	blink.lock.Lock()
	if blink.root == nil {
		blink.root = newNode(
			true, make(Keys, 0, blink.ary), make(nodes, 0, blink.ary+1),
		)
		blink.root.keys = make(Keys, 0, blink.ary)
		blink.root.isLeaf = true
	}
	parent = blink.root
	blink.lock.Unlock()

	result := insert(blink, parent, stack, key)
	if result == nil {
		atomic.AddUint64(&blink.number, 1)
		return nil
	}

	return result
}

func (blink *blink) multithreadedInsert(keys Keys) Keys {
	chunks := chunkKeys(keys, int64(blink.numRoutines))
	overwritten := make(Keys, len(keys))
	var offset uint64
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(chunk Keys, offset uint64) {
			defer wg.Done()
			stack := make(nodes, 0, blink.ary)

			for i := 0; i < len(chunk); i++ {
				result := blink.insert(chunk[i], &stack)
				stack.reset()
				overwritten[offset+uint64(i)] = result
			}
		}(chunk, offset)
		offset += uint64(len(chunk))
	}

	wg.Wait()

	return overwritten
}

// Insert will insert the provided keys into the b-tree and return
// a list of keys overwritten, if any.  Each insert is an O(log n)
// operation.
func (blink *blink) Insert(keys ...Key) Keys {
	if len(keys) > numberOfItemsBeforeMultithread {
		return blink.multithreadedInsert(keys)
	}
	overwritten := make(Keys, 0, len(keys))
	stack := make(nodes, 0, blink.ary)
	for _, k := range keys {
		overwritten = append(overwritten, blink.insert(k, &stack))
		stack.reset()
	}

	return overwritten
}

// Len returns the number of items in this b-link tree.
func (blink *blink) Len() uint64 {
	return atomic.LoadUint64(&blink.number)
}

func (blink *blink) get(key Key) Key {
	var parent *node
	blink.lock.RLock()
	parent = blink.root
	blink.lock.RUnlock()
	k := search(parent, key)
	if k == nil {
		return nil
	}

	if k.Compare(key) == 0 {
		return k
	}

	return nil
}

// Get will retrieve the keys if they exist in this tree.  If not,
// a nil is returned in the proper place in the list of keys.  Each
// lookup is O(log n) time complexity.
func (blink *blink) Get(keys ...Key) Keys {
	found := make(Keys, 0, len(keys))
	for _, k := range keys {
		found = append(found, blink.get(k))
	}

	return found
}

func (blink *blink) print(output *log.Logger) {
	output.Println(`PRINTING B-LINK`)
	if blink.root == nil {
		return
	}

	blink.root.print(output)
}

func newTree(ary, numRoutines uint64) *blink {
	return &blink{ary: ary, numRoutines: numRoutines}
}
