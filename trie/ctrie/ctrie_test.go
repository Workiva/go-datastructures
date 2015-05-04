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
