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

package skiplist

type mockEntry struct {
	values []int64
}

func (me *mockEntry) ValueAtDimension(dimension uint64) int64 {
	return me.values[dimension]
}

func newMockEntry(values ...int64) *mockEntry {
	return &mockEntry{values: values}
}

type mockInterval struct {
	lows, highs []int64
}

func (mi *mockInterval) LowAtDimension(dimension uint64) int64 {
	return mi.lows[dimension]
}

func (mi *mockInterval) HighAtDimension(dimension uint64) int64 {
	return mi.highs[dimension]
}

func newMockInterval(lows, highs []int64) *mockInterval {
	return &mockInterval{
		lows:  lows,
		highs: highs,
	}
}
