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

package palm

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/queue"
)

type actions []action

type action interface {
	operation() operation
	getKey() (Key, uint64) // returns nil if operation complete
	addResult(index uint64, result Key)
	len() uint64
	complete()
}

type getAction struct {
	result, keys Keys
	count, done  uint64
	completer    chan Keys
}

func (ga *getAction) complete() {
	ga.completer <- ga.result
	close(ga.completer)
}

func (ga *getAction) operation() operation {
	return get
}

func (ga *getAction) addResult(index uint64, result Key) {
	i := atomic.AddUint64(&ga.done, 1)
	i--
	if i >= uint64(len(ga.keys)) {
		return
	}
	ga.result[index] = result
	if i == uint64(len(ga.keys))-1 {
		ga.complete()
	}
}

func (ga *getAction) getKey() (Key, uint64) {
	index := atomic.AddUint64(&ga.count, 1)
	index-- // 0-index
	if index >= uint64(len(ga.keys)) {
		return nil, 0
	}

	return ga.keys[index], index
}

func (ga *getAction) len() uint64 {
	return uint64(len(ga.keys))
}

func newGetAction(keys Keys) *getAction {
	return &getAction{
		keys:      keys,
		completer: make(chan Keys),
		result:    make(Keys, len(keys)),
	}
}

type insertAction struct {
	keys      Keys
	completer chan bool
}

func (ia *insertAction) complete() {
	close(ia.completer)
}

func newInsertAction(keys Keys) *insertAction {
	return &insertAction{
		keys:      keys,
		completer: make(chan bool),
	}
}

func executeInParallel(q *queue.RingBuffer, fn func(interface{})) {
	if q == nil {
		return
	}

	todo, done := q.Len(), uint64(0)
	if todo == 0 {
		return
	}

	numCPU := 1
	if runtime.NumCPU() > 1 {
		numCPU = runtime.NumCPU() - 1
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < numCPU; i++ {
		go func() {
			for {
				ifc, err := q.Get()
				if err != nil {
					return
				}
				fn(ifc)
				if atomic.AddUint64(&done, 1) == todo {
					wg.Done()
					break
				}
			}
		}()
	}
	wg.Wait()
	q.Dispose()
}
