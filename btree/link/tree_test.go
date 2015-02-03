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

package link

import (
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getConsoleLogger() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}

func generateRandomKeys(num int) Keys {
	keys := make(Keys, 0, num)
	for i := 0; i < num; i++ {
		keys = append(keys, mockKey(uint64(rand.Uint32())))
	}

	return keys
}

func generateKeys(num int) Keys {
	keys := make(Keys, 0, num)
	for i := 0; i < num; i++ {
		keys = append(keys, mockKey(uint64(i)))
	}

	return keys
}

func TestSimpleInsert(t *testing.T) {
	k1 := mockKey(5)

	tree := newTree(8, 1)
	result := tree.Insert(k1)
	assert.Equal(t, Keys{nil}, result)
	assert.Equal(t, uint64(1), tree.Len())
	if !assert.Equal(t, Keys{k1}, tree.Get(k1)) {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsert(t *testing.T) {
	k1 := mockKey(10)
	k2 := mockKey(5)
	tree := newTree(8, 1)

	result := tree.Insert(k1, k2)
	assert.Equal(t, Keys{nil, nil}, result)
	assert.Equal(t, uint64(2), tree.Len())
	assert.Equal(t, Keys{k1, k2}, tree.Get(k1, k2))
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryReverseOrder(t *testing.T) {
	k1, k2, k3 := mockKey(15), mockKey(10), mockKey(5)
	tree := newTree(3, 1)

	result := tree.Insert(k1, k2, k3)
	assert.Equal(t, Keys{nil, nil, nil}, result)
	assert.Equal(t, uint64(3), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3}, tree.Get(k1, k2, k3)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAry(t *testing.T) {
	k1, k2, k3 := mockKey(5), mockKey(10), mockKey(15)
	tree := newTree(3, 1)

	result := tree.Insert(k1, k2, k3)
	assert.Equal(t, Keys{nil, nil, nil}, result)
	assert.Equal(t, uint64(3), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3}, tree.Get(k1, k2, k3)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryRandomOrder(t *testing.T) {
	k1, k2, k3 := mockKey(10), mockKey(5), mockKey(15)
	tree := newTree(3, 1)

	result := tree.Insert(k1, k2, k3)
	assert.Equal(t, Keys{nil, nil, nil}, result)
	assert.Equal(t, uint64(3), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3}, tree.Get(k1, k2, k3)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryReverseOrder(t *testing.T) {
	k1, k2, k3, k4 := mockKey(20), mockKey(15), mockKey(10), mockKey(5)
	tree := newTree(4, 1)

	result := tree.Insert(k1, k2, k3, k4)
	assert.Equal(t, Keys{nil, nil, nil, nil}, result)
	assert.Equal(t, uint64(4), tree.Len())
	if !assert.Equal(t, Keys{k3}, tree.Get(k3)) {
		tree.print(getConsoleLogger())
	}

	/*
		if !assert.Equal(t, Keys{k1, k2, k3, k4}, tree.Get(k1, k2, k3, k4)) {
			tree.print(getConsoleLogger())
		}*/
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAry(t *testing.T) {
	k1, k2, k3, k4 := mockKey(5), mockKey(10), mockKey(15), mockKey(20)
	tree := newTree(4, 1)

	result := tree.Insert(k1, k2, k3, k4)
	assert.Equal(t, Keys{nil, nil, nil, nil}, result)
	assert.Equal(t, uint64(4), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3, k4}, tree.Get(k1, k2, k3, k4)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryRandomOrder(t *testing.T) {
	k1, k2, k3, k4 := mockKey(10), mockKey(15), mockKey(20), mockKey(5)
	tree := newTree(4, 1)

	result := tree.Insert(k1, k2, k3, k4)
	assert.Equal(t, Keys{nil, nil, nil, nil}, result)
	assert.Equal(t, uint64(4), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3, k4}, tree.Get(k1, k2, k3, k4)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryMultiThreaded(t *testing.T) {
	keys := generateRandomKeys(16)
	tree := newTree(16, 8)

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys))
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesCascadingSplitsOddAry(t *testing.T) {
	keys := generateRandomKeys(16)
	tree := newTree(3, 8)

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys)) // about all we can assert, random may produce duplicates

	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesCascadingSplitsOddAryReverseOrder(t *testing.T) {
	keys := generateKeys(16)
	tree := newTree(3, 1)

	reversed := keys.reverse()

	result := tree.Insert(reversed...)
	assert.Len(t, result, len(keys)) // about all we can assert, random may produce duplicates

	println(`SHIT STARTS HERE`)
	if !assert.Equal(t, mockKey(3), tree.Get(mockKey(3))) {
		tree.print(getConsoleLogger())
	}
	/*
		if !assert.Equal(t, keys, tree.Get(keys...)) {
			tree.print(getConsoleLogger())
		}*/
	//checkTree(t, tree)
}

/*
func TestMultipleInsertCausesCascadingSplitsEvenAry(t *testing.T) {
	keys := generateRandomKeys(20)
	tree := newTree(4)

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys)) // about all we can assert, random may produce duplicates
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
}

func TestOverwriteOddAry(t *testing.T) {
	keys := generateRandomKeys(15)
	tree := newTree(3)
	duplicate := mockKey(uint64(keys[0].(mockKey)))

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys))
	oldLength := tree.Len()

	result = tree.Insert(duplicate)
	assert.Equal(t, Keys{keys[0]}, result)
	assert.Equal(t, oldLength, tree.Len())
}

func TestOverwriteEvenAry(t *testing.T) {
	keys := generateRandomKeys(15)
	tree := newTree(4)
	duplicate := mockKey(uint64(keys[0].(mockKey)))

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys))
	oldLength := tree.Len()

	result = tree.Insert(duplicate)
	assert.Equal(t, Keys{keys[0]}, result)
	assert.Equal(t, oldLength, tree.Len())
}

func BenchmarkSimpleAdd(b *testing.B) {
	numItems := 1000
	keys := generateRandomKeys(numItems)
	tree := newTree(16)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i%numItems])
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := 1000
	keys := generateRandomKeys(numItems)
	tree := newTree(16)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys[i%numItems])
	}
}*/

func BenchmarkBulkAdd(b *testing.B) {
	numItems := 1000
	keys := generateRandomKeys(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newTree(16, 4)
		tree.Insert(keys...)
	}
}
