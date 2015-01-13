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

package xfast

import "github.com/stretchr/testify/mock"

type mockEntry struct {
	mock.Mock
}

func (me *mockEntry) Key() uint64 {
	args := me.Called()
	return args.Get(0).(uint64)
}

func newMockEntry(key uint64) *mockEntry {
	me := new(mockEntry)
	me.On(`Key`).Return(key)
	return me
}
