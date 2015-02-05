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

package palm

import (
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkTree(t testing.TB, tree *ptree) bool {
	if tree.root == nil {
		return true
	}

	return checkNode(t, tree.root)
}

func checkNode(t testing.TB, n *node) bool {
	if len(n.keys) == 0 {
		assert.Len(t, n.nodes, 0)
		return false
	}

	if n.isLeaf {
		assert.Len(t, n.nodes, 0)
		return false
	}

	if !assert.Len(t, n.nodes, len(n.keys)+1) {
		return false
	}

	for i := 0; i < len(n.keys); i++ {
		if !assert.True(t, n.keys[i].Compare(n.nodes[i].keys[len(n.nodes[i].keys)-1]) >= 0) {
			t.Logf(`N: %+v %p, n.keys[i]: %+v, n.nodes[i]: %+v`, n, n, n.keys[i], n.nodes[i])
			return false
		}
	}

	if !assert.True(t, n.nodes[len(n.nodes)-1].key().Compare(n.keys[len(n.keys)-1]) > 0) {
		t.Logf(`m: %+v, %p, n.nodes[len(n.nodes)-1].key(): %+v, n.keys.last(): %+v`, n, n, n.nodes[len(n.nodes)-1].key(), n.keys[len(n.keys)-1])
		return false
	}
	for _, child := range n.nodes {
		if !assert.NotNil(t, child) {
			return false
		}
		if !checkNode(t, child) {
			return false
		}
	}

	return true
}

func getConsoleLogger() *log.Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}

func generateRandomKeys(num int) Keys {
	keys := make(Keys, 0, num)
	for i := 0; i < num; i++ {
		m := rand.Int()
		keys = append(keys, mockKey(m))
	}
	return keys
}

func generateKeys(num int) Keys {
	keys := make(Keys, 0, num)
	for i := 0; i < num; i++ {
		keys = append(keys, mockKey(i))
	}

	return keys
}

func TestSimpleInsert(t *testing.T) {
	tree := newTree(16)
	m1 := mockKey(1)

	tree.Insert(m1)
	assert.Equal(t, Keys{m1}, tree.Get(m1))
	assert.Equal(t, uint64(1), tree.Len())
	checkTree(t, tree)
}

func TestMultipleAdd(t *testing.T) {
	tree := newTree(16)
	m1 := mockKey(1)
	m2 := mockKey(10)

	tree.Insert(m1, m2)
	if !assert.Equal(t, Keys{m1, m2}, tree.Get(m1, m2)) {
		tree.print(getConsoleLogger())
	}
	assert.Equal(t, uint64(2), tree.Len())
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryReverseOrder(t *testing.T) {
	tree := newTree(3)
	keys := generateKeys(100)
	reversed := keys.reverse()

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAry(t *testing.T) {
	tree := newTree(3)
	keys := generateKeys(1000)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryRandomOrder(t *testing.T) {
	tree := newTree(3)
	keys := generateRandomKeys(100)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleBulkInsertOddAry(t *testing.T) {
	tree := newTree(3)
	keys1 := generateRandomKeys(100)
	keys2 := generateRandomKeys(100)

	tree.Insert(keys1...)
	tree.Insert(keys2...)

	if !assert.Equal(t, keys1, tree.Get(keys1...)) {
		tree.print(getConsoleLogger())
	}

	if !assert.Equal(t, keys2, tree.Get(keys2...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleBulkInsertEvenAry(t *testing.T) {
	tree := newTree(4)
	keys1 := generateRandomKeys(100)
	keys2 := generateRandomKeys(100)

	tree.Insert(keys1...)
	tree.Insert(keys2...)

	if !assert.Equal(t, keys1, tree.Get(keys1...)) {
		tree.print(getConsoleLogger())
	}

	if !assert.Equal(t, keys2, tree.Get(keys2...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryReverseOrder(t *testing.T) {
	tree := newTree(4)
	keys := generateKeys(1000)
	reversed := keys.reverse()

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAry(t *testing.T) {
	tree := newTree(4)
	keys := generateKeys(1000)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryRandomOrder(t *testing.T) {
	tree := newTree(4)
	keys := generateRandomKeys(1000)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestInsertOverwrite(t *testing.T) {
	tree := newTree(4)
	keys := generateKeys(10)
	duplicate := mockKey(0)
	tree.Insert(keys...)

	tree.Insert(duplicate)
	assert.Equal(t, Keys{duplicate}, tree.Get(duplicate))
	checkTree(t, tree)
}

func TestSimultaneousReadsAndWrites(t *testing.T) {
	numLoops := 3
	keys := make([]Keys, 0, numLoops)
	for i := 0; i < numLoops; i++ {
		keys = append(keys, generateRandomKeys(1000))
	}

	tree := newTree(16)
	var wg sync.WaitGroup
	wg.Add(numLoops)
	for i := 0; i < numLoops; i++ {
		go func(i int) {
			tree.Insert(keys[i]...)
			tree.Get(keys[i]...)
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < numLoops; i++ {
		assert.Equal(t, keys[i], tree.Get(keys[i]...))
	}
	checkTree(t, tree)
}

func BenchmarkReadAndWrites(b *testing.B) {
	numItems := 1000
	keys := make([]Keys, 0, b.N)
	for i := 0; i < b.N; i++ {
		keys = append(keys, generateRandomKeys(numItems))
	}

	tree := newTree(16)
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i]...)
		tree.Get(keys[i]...)
		wg.Done()
	}

	wg.Wait()

}

func BenchmarkBulkAdd(b *testing.B) {
	numItems := 10000
	keys := generateKeys(numItems)
	keySet := make([]Keys, 0, b.N)
	for i := 0; i < b.N; i++ {
		cp := make(Keys, len(keys))
		copy(cp, keys)
		keySet = append(keySet, cp)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newTree(1024)
		tree.Insert(keySet[i]...)
	}
}

func BenchmarkBulkAddToExisting(b *testing.B) {
	numItems := 10000
	keySet := make([]Keys, 0, b.N)
	for i := 0; i < b.N; i++ {
		keySet = append(keySet, generateRandomKeys(numItems))
	}

	tree := newTree(1024)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keySet[i]...)
	}
}
