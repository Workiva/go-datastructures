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
Package futures is useful for broadcasting an identical message to a multitude
of listeners as opposed to channels which will choose a listener at random
if multiple listeners are listening to the same channel.  The future will
also cache the result so any future interest will be immediately returned
to the consumer.
*/
package futures

import (
	"fmt"
	"sync"
	"time"
)

// Completer is a channel that the future expects to receive
// a result on.  The future only receives on this channel.
type Completer <-chan interface{}

// Future represents an object that can be used to perform asynchronous
// tasks.  The constructor of the future will complete it, and listeners
// will block on getresult until a result is received.  This is different
// from a channel in that the future is only completed once, and anyone
// listening on the future will get the result, regardless of the number
// of listeners.
type Future struct {
	triggered bool // because item can technically be nil and still be valid
	item      interface{}
	err       error
	lock      sync.Mutex
	wg        sync.WaitGroup
}

// GetResult will immediately fetch the result if it exists
// or wait on the result until it is ready.
func (f *Future) GetResult() (interface{}, error) {
	f.lock.Lock()
	if f.triggered {
		f.lock.Unlock()
		return f.item, f.err
	}
	f.lock.Unlock()

	f.wg.Wait()
	return f.item, f.err
}

func (f *Future) setItem(item interface{}, err error) {
	f.lock.Lock()
	f.triggered = true
	f.item = item
	f.err = err
	f.lock.Unlock()
	f.wg.Done()
}

func listenForResult(f *Future, ch Completer, timeout time.Duration, wg *sync.WaitGroup) {
	wg.Done()
	select {
	case item := <-ch:
		f.setItem(item, nil)
	case <-time.After(timeout):
		f.setItem(nil, fmt.Errorf(`Timeout after %f seconds.`, timeout.Seconds()))
	}
}

// New is the constructor to generate a new future.  Pass the completed
// item to the toComplete channel and any listeners will get
// notified.  If timeout is hit before toComplete is called,
// any listeners will get passed an error.
func New(completer Completer, timeout time.Duration) *Future {
	f := &Future{}
	f.wg.Add(1)
	var wg sync.WaitGroup
	wg.Add(1)
	go listenForResult(f, completer, timeout, &wg)
	wg.Wait()
	return f
}
