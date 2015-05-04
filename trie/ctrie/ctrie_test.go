package ctrie

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertAndLookup(t *testing.T) {
	assert := assert.New(t)
	ctrie := New()

	_, ok := ctrie.Lookup([]byte("foo"))
	assert.False(ok)

	ctrie.Insert([]byte("foo"), "bar")
	val, ok := ctrie.Lookup([]byte("foo"))
	assert.True(ok)
	assert.Equal("bar", val)

	ctrie.Insert([]byte("fooooo"), "baz")
	val, ok = ctrie.Lookup([]byte("foo"))
	assert.True(ok)
	assert.Equal("bar", val)
	val, ok = ctrie.Lookup([]byte("fooooo"))
	assert.True(ok)
	assert.Equal("baz", val)

	for i := 0; i < 100; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), "blah")
	}
	for i := 0; i < 100; i++ {
		val, ok = ctrie.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal("blah", val)
	}

	val, ok = ctrie.Lookup([]byte("foo"))
	assert.True(ok)
	assert.Equal("bar", val)
	ctrie.Insert([]byte("foo"), "qux")
	val, ok = ctrie.Lookup([]byte("foo"))
	assert.True(ok)
	assert.Equal("qux", val)
}

func BenchmarkInsert(b *testing.B) {
	ctrie := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrie.Insert([]byte("foo"), 0)
	}
}

func BenchmarkLookup(b *testing.B) {
	numItems := 1000
	ctrie := New()
	for i := 0; i < numItems; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	key := []byte(strconv.Itoa(numItems / 2))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctrie.Lookup(key)
	}
}
