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

package plus

func chunkKeys(ks keys, numParts int64) []keys {
	parts := make([]keys, numParts)
	for i := int64(0); i < numParts; i++ {
		parts[i] = ks[i*int64(len(ks))/numParts : (i+1)*int64(len(ks))/numParts]
	}
	return parts
}

type mockKey struct {
	value int
}

func (mk *mockKey) Compare(other Key) int {
	key := other.(*mockKey)
	if key.value == mk.value {
		return 0
	}
	if key.value > mk.value {
		return 1
	}

	return -1
}

func newMockKey(value int) *mockKey {
	return &mockKey{value}
}
