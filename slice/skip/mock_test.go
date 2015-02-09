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

package skip

import "github.com/stretchr/testify/mock"

type mockEntry uint64

func (me mockEntry) Compare(other Entry) int {
	otherU := other.(mockEntry)
	if me == otherU {
		return 0
	}

	if me > otherU {
		return 1
	}

	return -1
}

func newMockEntry(key uint64) mockEntry {
	return mockEntry(key)
}

type mockIterator struct {
	mock.Mock
}

func (mi *mockIterator) Next() bool {
	args := mi.Called()
	return args.Bool(0)
}

func (mi *mockIterator) Value() Entry {
	args := mi.Called()
	result, ok := args.Get(0).(Entry)
	if !ok {
		return nil
	}

	return result
}

func (mi *mockIterator) exhaust() Entries {
	return nil
}
