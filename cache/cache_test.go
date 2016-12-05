package cache

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvictionPolicy(t *testing.T) {
	c := &cache{keyList: list.New()}
	EvictionPolicy(LeastRecentlyUsed)(c)
	accessed, added := c.recordAccess("foo"), c.recordAdd("foo")
	assert.NotNil(t, accessed)
	assert.Nil(t, added)

	c = &cache{keyList: list.New()}
	EvictionPolicy(LeastRecentlyAdded)(c)
	accessed, added = c.recordAccess("foo"), c.recordAdd("foo")
	assert.Nil(t, accessed)
	assert.NotNil(t, added)
}

func TestNew(t *testing.T) {
	optionApplied := false
	option := func(*cache) {
		optionApplied = true
	}

	c := New(314159, option).(*cache)

	assert.Equal(t, uint64(314159), c.cap)
	assert.Equal(t, uint64(0), c.size)
	assert.NotNil(t, c.items)
	assert.NotNil(t, c.keyList)
	assert.True(t, optionApplied)

	accessed, added := c.recordAccess("foo"), c.recordAdd("foo")
	assert.NotNil(t, accessed)
	assert.Nil(t, added)
}

type testItem uint64

func (ti testItem) Size() uint64 {
	return uint64(ti)
}

func TestPutGetRemoveSize(t *testing.T) {
	keys := []string{"foo", "bar", "baz"}
	testCases := []struct {
		label         string
		cache         Cache
		useCache      func(c Cache)
		expectedSize  uint64
		expectedItems []Item
	}{{
		label: "Items added, key doesn't exist",
		cache: New(10000),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
		},
		expectedSize:  1,
		expectedItems: []Item{testItem(1), nil, nil},
	}, {
		label: "Items added, key exists",
		cache: New(10000),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
			c.Put("foo", testItem(10))
		},
		expectedSize:  10,
		expectedItems: []Item{testItem(10), nil, nil},
	}, {
		label: "Items added, LRA eviction",
		cache: New(2, EvictionPolicy(LeastRecentlyAdded)),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
			c.Put("bar", testItem(1))
			c.Get("foo")
			c.Put("baz", testItem(1))
		},
		expectedSize:  2,
		expectedItems: []Item{nil, testItem(1), testItem(1)},
	}, {
		label: "Items added, LRU eviction",
		cache: New(2, EvictionPolicy(LeastRecentlyUsed)),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
			c.Put("bar", testItem(1))
			c.Get("foo")
			c.Put("baz", testItem(1))
		},
		expectedSize:  2,
		expectedItems: []Item{testItem(1), nil, testItem(1)},
	}, {
		label: "Items removed, key doesn't exist",
		cache: New(10000),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
			c.Remove("baz")
		},
		expectedSize:  1,
		expectedItems: []Item{testItem(1), nil, nil},
	}, {
		label: "Items removed, key exists",
		cache: New(10000),
		useCache: func(c Cache) {
			c.Put("foo", testItem(1))
			c.Remove("foo")
		},
		expectedSize:  0,
		expectedItems: []Item{nil, nil, nil},
	}}

	for _, testCase := range testCases {
		t.Log(testCase.label)
		testCase.useCache(testCase.cache)
		assert.Equal(t, testCase.expectedSize, testCase.cache.Size())
		assert.Equal(t, testCase.expectedItems, testCase.cache.Get(keys...))
	}
}
