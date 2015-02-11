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
Package err implements a threadsafe error interface.  In my places,
I found myself needing a lock to protect writing to a common error interface
from multiple go routines (channels are great but slow).  This just makes
that process more convenient.
*/
package err

import "sync"

// Error is a struct that holds an error and allows this error
// to be set and retrieved in a threadsafe manner.
type Error struct {
	lock sync.RWMutex
	err  error
}

// Set will set the error of this structure to the provided
// value.
func (e *Error) Set(err error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.err = err
}

// Get will return any error associated with this structure.
func (e *Error) Get() error {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.err
}

// New is a constructor to generate a new error object
// that can be set and retrieved in a threadsafe manner.
func New() *Error {
	return &Error{}
}
