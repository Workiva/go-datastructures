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

package rangetree

type mockEntry struct {
	id         uint64
	dimensions []int64
}

func (me *mockEntry) ID() uint64 {
	return me.id
}

func (me *mockEntry) ValueAtDimension(dimension uint64) int64 {
	return me.dimensions[dimension-1]
}

func constructMockEntry(id uint64, values ...int64) *mockEntry {
	return &mockEntry{
		id:         id,
		dimensions: values,
	}
}

type dimension struct {
	low, high int64
}

type mockInterval struct {
	dimensions []dimension
}

func (mi *mockInterval) LowAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].low
}

func (mi *mockInterval) HighAtDimension(dimension uint64) int64 {
	return mi.dimensions[dimension-1].high
}

func constructMockInterval(dimensions ...dimension) *mockInterval {
	return &mockInterval{
		dimensions: dimensions,
	}
}
