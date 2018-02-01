/*
Copyright 2016 Workiva, LLC
Copyright 2016 Sokolov Yura aka funny_falcon

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

package futures

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrFutureCanceled signals that futures in canceled by a call to `f.Cancel()`
var ErrFutureCanceled = errors.New("future canceled")

// Selectable is a future with channel exposed for external `select`.
// Many simultaneous listeners may wait for result either with `f.Value()`
// or by selecting/fetching from `f.WaitChan()`, which is closed when future
// fulfilled.
// Selectable contains sync.Mutex, so it is not movable/copyable.
type Selectable struct {
	m      sync.Mutex
	val    interface{}
	err    error
	wait   chan struct{}
	filled uint32
}

// NewSelectable returns new selectable future.
// Note: this method is for backward compatibility.
// You may allocate it directly on stack or embedding into larger structure
func NewSelectable() *Selectable {
	return &Selectable{}
}

func (f *Selectable) wchan() <-chan struct{} {
	f.m.Lock()
	if f.wait == nil {
		f.wait = make(chan struct{})
	}
	ch := f.wait
	f.m.Unlock()
	return ch
}

// WaitChan returns channel, which is closed when future is fullfilled.
func (f *Selectable) WaitChan() <-chan struct{} {
	if atomic.LoadUint32(&f.filled) == 1 {
		return closed
	}
	return f.wchan()
}

// GetResult waits for future to be fullfilled and returns value or error,
// whatever is set first
func (f *Selectable) GetResult() (interface{}, error) {
	if atomic.LoadUint32(&f.filled) == 0 {
		<-f.wchan()
	}
	return f.val, f.err
}

// Fill sets value for future, if it were not already fullfilled
// Returns error, if it were already set to future.
func (f *Selectable) Fill(v interface{}, e error) error {
	f.m.Lock()
	if f.filled == 0 {
		f.val = v
		f.err = e
		atomic.StoreUint32(&f.filled, 1)
		w := f.wait
		f.wait = closed
		if w != nil {
			close(w)
		}
	}
	f.m.Unlock()
	return f.err
}

// SetValue is alias for Fill(v, nil)
func (f *Selectable) SetValue(v interface{}) error {
	return f.Fill(v, nil)
}

// SetError is alias for Fill(nil, e)
func (f *Selectable) SetError(e error) {
	f.Fill(nil, e)
}

// Cancel is alias for SetError(ErrFutureCanceled)
func (f *Selectable) Cancel() {
	f.SetError(ErrFutureCanceled)
}

var closed = make(chan struct{})

func init() {
	close(closed)
}
