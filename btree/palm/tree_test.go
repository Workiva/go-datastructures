package palm

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"testing"
	"time"

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

	result := tree.Insert(m1)
	assert.Equal(t, Keys{nil}, result)
	assert.Equal(t, Keys{m1}, tree.Get(m1))
	assert.Equal(t, uint64(1), tree.Len())
}

func TestMultipleAdd(t *testing.T) {
	tree := newTree(16)
	m1 := mockKey(1)
	m2 := mockKey(10)

	result := tree.Insert(m1, m2)
	assert.Equal(t, Keys{nil, nil}, result)
	if !assert.Equal(t, Keys{m1, m2}, tree.Get(m1, m2)) {
		tree.print(getConsoleLogger())
	}
	assert.Equal(t, uint64(2), tree.Len())
}

func TestMultipleInsertCausesSplitOddAryReverseOrder(t *testing.T) {
	tree := newTree(8)
	keys := generateRandomKeys(16)
	reversed := keys.reverse()

	result := tree.Insert(reversed...)
	assert.Len(t, result, len(keys))
	time.Sleep(100 * time.Millisecond)
	//if !assert.Equal(t, keys, tree.Get(keys...)) {
	//	tree.print(getConsoleLogger())
	//}
	time.Sleep(10 * time.Millisecond)
	//tree.print(getConsoleLogger())
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

	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree := newTree(2056)
		tree.Insert(keySet[i]...)
	}
}
