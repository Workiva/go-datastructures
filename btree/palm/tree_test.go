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

	"github.com/Workiva/go-datastructures/common"
)

func checkTree(t testing.TB, tree *ptree) bool {
	return true

	if tree.root == nil {
		return true
	}

	return checkNode(t, tree.root)
}

func checkNode(t testing.TB, n *node) bool {
	if n.keys.len() == 0 {
		assert.Equal(t, uint64(0), n.nodes.len())
		return false
	}

	if n.isLeaf {
		assert.Equal(t, uint64(0), n.nodes.len())
		return false
	}

	if !assert.Equal(t, n.keys.len()+1, n.nodes.len()) {
		return false
	}

	for i, k := range n.keys.list {
		nd := n.nodes.list[i]
		if !assert.NotNil(t, nd) {
			return false
		}

		if !assert.True(t, k.Compare(nd.key()) >= 0) {
			t.Logf(`N: %+v %p, n.keys[i]: %+v, n.nodes[i]: %+v`, n, n, k, nd)
			return false
		}
	}

	k := n.keys.last()
	nd := n.nodes.byPosition(n.nodes.len() - 1)
	if !assert.True(t, k.Compare(nd.key()) < 0) {
		t.Logf(`m: %+v, %p, n.nodes[len(n.nodes)-1].key(): %+v, n.keys.last(): %+v`, n, n, nd, k)
		return false
	}
	for _, child := range n.nodes.list {
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

func generateRandomKeys(num int) common.Comparators {
	keys := make(common.Comparators, 0, num)
	for i := 0; i < num; i++ {
		m := rand.Int()
		keys = append(keys, mockKey(m%50))
	}
	return keys
}

func generateKeys(num int) common.Comparators {
	keys := make(common.Comparators, 0, num)
	for i := 0; i < num; i++ {
		keys = append(keys, mockKey(i))
	}

	return keys
}

func TestSimpleInsert(t *testing.T) {
	tree := newTree(16, 16)
	defer tree.Dispose()
	m1 := mockKey(1)

	tree.Insert(m1)
	assert.Equal(t, common.Comparators{m1}, tree.Get(m1))
	assert.Equal(t, uint64(1), tree.Len())
	checkTree(t, tree)
}

func TestSimpleDelete(t *testing.T) {
	tree := newTree(8, 8)
	defer tree.Dispose()
	m1 := mockKey(1)
	tree.Insert(m1)

	tree.Delete(m1)
	assert.Equal(t, uint64(0), tree.Len())
	assert.Equal(t, common.Comparators{nil}, tree.Get(m1))
	checkTree(t, tree)
}

func TestMultipleAdd(t *testing.T) {
	tree := newTree(16, 16)
	defer tree.Dispose()
	m1 := mockKey(1)
	m2 := mockKey(10)

	tree.Insert(m1, m2)
	if !assert.Equal(t, common.Comparators{m1, m2}, tree.Get(m1, m2)) {
		tree.print(getConsoleLogger())
	}
	assert.Equal(t, uint64(2), tree.Len())
	checkTree(t, tree)
}

func TestMultipleDelete(t *testing.T) {
	tree := newTree(16, 16)
	defer tree.Dispose()
	m1 := mockKey(1)
	m2 := mockKey(10)
	tree.Insert(m1, m2)

	tree.Delete(m1, m2)
	assert.Equal(t, uint64(0), tree.Len())
	assert.Equal(t, common.Comparators{nil, nil}, tree.Get(m1, m2))
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryReverseOrder(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()
	keys := generateKeys(100)
	reversed := reverseKeys(keys)

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleDeleteOddAryReverseOrder(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()
	keys := generateKeys(100)
	reversed := reverseKeys(keys)
	tree.Insert(reversed...)
	assert.Equal(t, uint64(100), tree.Len())

	tree.Delete(reversed...)
	assert.Equal(t, uint64(0), tree.Len())
	for _, k := range reversed {
		assert.Equal(t, common.Comparators{nil}, tree.Get(k))
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAry(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()
	keys := generateKeys(100)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitOddAryRandomOrder(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()
	keys := generateRandomKeys(10)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleBulkInsertOddAry(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()
	keys1 := generateRandomKeys(100)
	keys2 := generateRandomKeys(100)

	tree.Insert(keys1...)

	if !assert.Equal(t, keys1, tree.Get(keys1...)) {
		tree.print(getConsoleLogger())
	}

	tree.Insert(keys2...)

	if !assert.Equal(t, keys2, tree.Get(keys2...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleBulkInsertEvenAry(t *testing.T) {
	tree := newTree(4, 4)
	defer tree.Dispose()
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
	tree := newTree(4, 4)
	defer tree.Dispose()
	keys := generateKeys(100)
	reversed := reverseKeys(keys)

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAry(t *testing.T) {
	tree := newTree(4, 4)
	defer tree.Dispose()
	keys := generateKeys(100)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestMultipleInsertCausesSplitEvenAryRandomOrder(t *testing.T) {
	tree := newTree(4, 4)
	defer tree.Dispose()
	keys := generateRandomKeys(100)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
	checkTree(t, tree)
}

func TestInsertOverwrite(t *testing.T) {
	tree := newTree(4, 4)
	defer tree.Dispose()
	keys := generateKeys(10)
	duplicate := mockKey(0)
	tree.Insert(keys...)

	tree.Insert(duplicate)
	assert.Equal(t, common.Comparators{duplicate}, tree.Get(duplicate))
	checkTree(t, tree)
}

func TestSimultaneousReadsAndWrites(t *testing.T) {
	numLoops := 3
	keys := make([]common.Comparators, 0, numLoops)
	for i := 0; i < numLoops; i++ {
		keys = append(keys, generateRandomKeys(10))
	}

	tree := newTree(16, 16)
	defer tree.Dispose()
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

func TestInsertAndDelete(t *testing.T) {
	tree := newTree(1024, 1024)
	defer tree.Dispose()

	keys := generateKeys(100)
	keys1 := keys[:50]
	keys2 := keys[50:]
	tree.Insert(keys1...)
	assert.Equal(t, uint64(len(keys1)), tree.Len())
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		tree.Insert(keys2...)
		wg.Done()
	}()

	go func() {
		tree.Delete(keys1...)
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, uint64(len(keys2)), tree.Len())
	assert.Equal(t, keys2, tree.Get(keys2...))
}

func TestInsertAndDeletesWithSplits(t *testing.T) {
	tree := newTree(3, 3)
	defer tree.Dispose()

	keys := generateKeys(100)
	keys1 := keys[:50]
	keys2 := keys[50:]
	tree.Insert(keys1...)
	assert.Equal(t, uint64(len(keys1)), tree.Len())
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		tree.Insert(keys2...)
		wg.Done()
	}()

	go func() {
		tree.Delete(keys1...)
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, uint64(len(keys2)), tree.Len())
	assert.Equal(t, keys2, tree.Get(keys2...))
}

func BenchmarkReadAndWrites(b *testing.B) {
	numItems := 1000
	keys := make([]common.Comparators, 0, b.N)
	for i := 0; i < b.N; i++ {
		keys = append(keys, generateRandomKeys(numItems))
	}

	tree := newTree(16, 8)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i]...)
		tree.Get(keys[i]...)
	}
}

func BenchmarkSimultaneousReadsAndWrites(b *testing.B) {
	numItems := 1000
	numRoutines := 8
	keys := make([]common.Comparators, 0, numRoutines)
	for i := 0; i < numRoutines; i++ {
		keys = append(keys, generateRandomKeys(numItems))
	}

	tree := newTree(16, 8)
	var wg sync.WaitGroup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(numRoutines)
		for j := 0; j < numRoutines; j++ {
			go func(j int) {
				tree.Insert(keys[j]...)
				tree.Get(keys[j]...)
				wg.Done()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkBulkAdd(b *testing.B) {
	numItems := 10000
	keys := generateRandomKeys(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newTree(8, 8)
		tree.Insert(keys...)
	}
}

func BenchmarkAdd(b *testing.B) {
	numItems := 500
	keys := generateRandomKeys(numItems)
	tree := newTree(32, 8)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keys[i%numItems])
	}
}

func BenchmarkBulkAddToExisting(b *testing.B) {
	numItems := 100000
	keySet := make([]common.Comparators, 0, b.N)
	for i := 0; i < b.N; i++ {
		keySet = append(keySet, generateRandomKeys(numItems))
	}

	tree := newTree(8, 8)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Insert(keySet[i]...)
	}
}

func BenchmarkGet(b *testing.B) {
	numItems := 10000
	keys := generateRandomKeys(numItems)
	tree := newTree(32, 8)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys[i%numItems])
	}
}

func BenchmarkBulkGet(b *testing.B) {
	numItems := 1000
	keys := generateRandomKeys(numItems)
	tree := newTree(16, 8)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Get(keys...)
	}
}

func BenchmarkDelete(b *testing.B) {
	numItems := b.N
	keys := generateRandomKeys(numItems)
	tree := newTree(8, 8)
	tree.Insert(keys...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree.Delete(keys[i%numItems])
	}
}

func BenchmarkBulkDelete(b *testing.B) {
	numItems := 10000
	keys := generateRandomKeys(numItems)
	trees := make([]*ptree, 0, b.N)
	for i := 0; i < b.N; i++ {
		tree := newTree(8, 8)
		tree.Insert(keys...)
		trees = append(trees, tree)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		trees[i].Delete(keys...)
	}
}
