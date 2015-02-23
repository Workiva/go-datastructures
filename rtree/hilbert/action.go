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

package hilbert

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/rtree"
)

type actions []action

type action interface {
	operation() operation
	keys() hilberts
	rects() []*hilbertBundle
	complete()
	addNode(int64, *node)
	nodes() []*node
}

type getAction struct {
	result    rtree.Rectangles
	completer *sync.WaitGroup
	lookup    *rectangle
}

func (ga *getAction) complete() {
	ga.completer.Done()
}

func (ga *getAction) operation() operation {
	return get
}

func (ga *getAction) keys() hilberts {
	return nil
}

func (ga *getAction) addNode(i int64, n *node) {
	return // not necessary for gets
}

func (ga *getAction) nodes() []*node {
	return nil
}

func (ga *getAction) rects() []*hilbertBundle {
	return nil
}

func newGetAction(rect rtree.Rectangle) *getAction {
	r := newRectangeFromRect(rect)
	ga := &getAction{
		completer: new(sync.WaitGroup),
		lookup:    r,
	}
	ga.completer.Add(1)
	return ga
}

type insertAction struct {
	rs        []*hilbertBundle
	completer *sync.WaitGroup
	ns        []*node
}

func (ia *insertAction) complete() {
	ia.completer.Done()
}

func (ia *insertAction) operation() operation {
	return add
}

func (ia *insertAction) keys() hilberts {
	return nil
}

func (ia *insertAction) addNode(i int64, n *node) {
	ia.ns[i] = n
}

func (ia *insertAction) nodes() []*node {
	return ia.ns
}

func (ia *insertAction) rects() []*hilbertBundle {
	return ia.rs
}

func newInsertAction(rects rtree.Rectangles) *insertAction {
	ia := &insertAction{
		rs:        bundlesFromRects(rects...),
		completer: new(sync.WaitGroup),
		ns:        make([]*node, len(rects)),
	}
	ia.completer.Add(1)
	return ia
}

type removeAction struct {
	*insertAction
}

func (ra *removeAction) operation() operation {
	return remove
}

func newRemoveAction(rects rtree.Rectangles) *removeAction {
	return &removeAction{
		newInsertAction(rects),
	}
}

func minUint64(choices ...uint64) uint64 {
	min := choices[0]
	for i := 1; i < len(choices); i++ {
		if choices[i] < min {
			min = choices[i]
		}
	}

	return min
}

type interfaces []interface{}

func executeInterfacesInParallel(ifs interfaces, fn func(interface{})) {
	if len(ifs) == 0 {
		return
	}

	done := int64(-1)
	numCPU := uint64(runtime.NumCPU())
	if numCPU > 1 {
		numCPU--
	}

	numCPU = minUint64(numCPU, uint64(len(ifs)))

	var wg sync.WaitGroup
	wg.Add(int(numCPU))

	for i := uint64(0); i < numCPU; i++ {
		go func() {
			defer wg.Done()

			for {
				i := atomic.AddInt64(&done, 1)
				if i >= int64(len(ifs)) {
					return
				}

				fn(ifs[i])
			}
		}()
	}

	wg.Wait()
}

func executeInterfacesInSerial(ifs interfaces, fn func(interface{})) {
	if len(ifs) == 0 {
		return
	}

	for _, ifc := range ifs {
		fn(ifc)
	}
}
