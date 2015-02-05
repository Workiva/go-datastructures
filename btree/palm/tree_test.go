package palm

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
}

func TestMultipleInsertCausesSplitOddAryReverseOrder(t *testing.T) {
	tree := newTree(3)
	keys := generateKeys(1000)
	reversed := keys.reverse()

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertCausesSplitOddAry(t *testing.T) {
	tree := newTree(3)
	keys := generateRandomKeys(1000)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
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
}

func TestMultipleInsertCausesSplitEvenAryReverseOrder(t *testing.T) {
	tree := newTree(4)
	keys := generateKeys(1000)
	reversed := keys.reverse()

	tree.Insert(reversed...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
}

func TestMultipleInsertCausesSplitEvenAry(t *testing.T) {
	tree := newTree(4)
	keys := generateRandomKeys(1000)

	tree.Insert(keys...)
	if !assert.Equal(t, keys, tree.Get(keys...)) {
		tree.print(getConsoleLogger())
	}
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
