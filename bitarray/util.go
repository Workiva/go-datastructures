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

package bitarray

// maxInt64 returns the highest integer in the provided list of int32s
func maxInt64(ints ...int64) int64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}

// maxUint64 returns the highest integer in the provided list of int32s
func maxUint64(ints ...uint64) uint64 {
	maxInt := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] > maxInt {
			maxInt = ints[i]
		}
	}

	return maxInt
}
