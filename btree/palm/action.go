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
	//"log"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Workiva/go-datastructures/queue"
)

type gets Keys

type adds Keys

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
	completer    *sync.WaitGroup
}

func (ga *getAction) complete() {
	ga.completer.Done()
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
	ga := &getAction{
		keys:      keys,
		completer: new(sync.WaitGroup),
		result:    make(Keys, len(keys)),
	}
	ga.completer.Add(1)
	return ga
}

type insertAction struct {
	keys      Keys
	completer *sync.WaitGroup
}

func (ia *insertAction) complete() {
	ia.completer.Done()
}

func newInsertAction(keys Keys) *insertAction {
	ia := &insertAction{
		keys:      keys,
		completer: new(sync.WaitGroup),
	}
	ia.completer.Add(1)
	return ia
}

func executeInParallel(q *queue.RingBuffer, fn func(interface{})) {
	if q == nil {
		return
	}

	todo, done := q.Len(), uint64(0)
	if todo == 0 {
		return
	}

	goRoutines := minUint64(todo, uint64(runtime.NumCPU()-1))

	var wg sync.WaitGroup
	wg.Add(1)

	for i := uint64(0); i < goRoutines; i++ {
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
