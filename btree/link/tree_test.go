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

func TestSimpleInsert(t *testing.T) {
	k1 := mockKey(5)

	tree := newTree(8)
	result := tree.Insert(k1)
	assert.Equal(t, Keys{nil}, result)
	assert.Equal(t, uint64(1), tree.Len())
	assert.Equal(t, Keys{k1}, tree.Get(k1))
}

func TestMultipleInsert(t *testing.T) {
	k1 := mockKey(10)
	k2 := mockKey(5)
	tree := newTree(8)

	result := tree.Insert(k1, k2)
	assert.Equal(t, Keys{nil, nil}, result)
	assert.Equal(t, uint64(2), tree.Len())
	assert.Equal(t, Keys{k1, k2}, tree.Get(k1, k2))
}

func TestMultipleInsertCausesSplitOddAry(t *testing.T) {
	k1, k2, k3 := mockKey(15), mockKey(10), mockKey(5)
	tree := newTree(3)

	result := tree.Insert(k1, k2, k3)
	assert.Equal(t, Keys{nil, nil, nil}, result)
	assert.Equal(t, uint64(3), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3}, tree.Get(k1, k2, k3)) {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertCausesSplitEvenAry(t *testing.T) {
	k1, k2, k3, k4 := mockKey(20), mockKey(15), mockKey(10), mockKey(5)
	tree := newTree(4)

	result := tree.Insert(k1, k2, k3, k4)
	assert.Equal(t, Keys{nil, nil, nil, nil}, result)
	assert.Equal(t, uint64(4), tree.Len())
	if !assert.Equal(t, Keys{k1, k2, k3, k4}, tree.Get(k1, k2, k3, k4)) {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertCausesCascadingSplitsOddAry(t *testing.T) {
	keys := generateRandomKeys(15)
	tree := newTree(3)

	result := tree.Insert(keys...)
	assert.Len(t, result, len(keys)) // about all we can assert, random may produce duplicates
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
}

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
