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

package link

import (
	"log"
	"sync"
	"sync/atomic"
)

type blink struct {
	root        *node
	lock        sync.RWMutex
	number, ary uint64
}

func (blink *blink) insert(key Key) Key {
	var parent *node
	blink.lock.Lock()
	if blink.root == nil {
		blink.root = newNode(true)
		blink.root.keys = make(Keys, 0, blink.ary)
		blink.root.isLeaf = true
	}
	parent = blink.root
	blink.lock.Unlock()

	result := insert(blink, parent, make(nodes, 0, blink.ary), key)
	if result == nil {
		atomic.AddUint64(&blink.number, 1)
		return nil
	}

	return result
}

func (blink *blink) Insert(keys ...Key) Keys {
	overwritten := make(Keys, 0, len(keys))
	for _, k := range keys {
		overwritten = append(overwritten, blink.insert(k))
	}

	return overwritten
}

func (blink *blink) Len() uint64 {
	return atomic.LoadUint64(&blink.number)
}

func (blink *blink) get(key Key) Key {
	n, index := search(blink.root, key)
	if index == len(n.keys) {
		return nil
	}

	if n.keys[index].Compare(key) == 0 {
		return n.keys[index]
	}

	return nil
}

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

func newTree(ary uint64) *blink {
	return &blink{ary: ary}
}
