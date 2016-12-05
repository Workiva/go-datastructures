package cache

import (
	"container/list"
	"testing"
)

func TestEvictionPolicy(t *testing.T) {
	c := &cache{keyList: list.New()}
	EvictionPolicy(LeastRecentlyUsed)(c)
	if accessed, added := c.recordAccess("foo"), c.recordAdd("foo"); accessed == nil || added != nil {
		t.Errorf("EvictionPolicy failed to set LRU policy")
	}

	c = &cache{keyList: list.New()}
	EvictionPolicy(LeastRecentlyAdded)(c)
	if accessed, added := c.recordAccess("foo"), c.recordAdd("foo"); accessed != nil || added == nil {
		t.Errorf("EvictionPolicy failed to set LRU policy")
	}
}

func TestNew(t *testing.T) {
	optionApplied := false
	option := func(*cache) {
		optionApplied = true
	}

	c := New(314159, option).(*cache)

	if c.cap != 314159 {
		t.Errorf("Expected cache capacity of %d", 314159)
	}
	if c.size != 0 {
		t.Errorf("Expected initial size of zero")
	}
	if c.items == nil {
		t.Errorf("Expected items to be initialized")
	}
	if c.keyList == nil {
		t.Errorf("Expected keyList to be initialized")
	}
	if !optionApplied {
		t.Errorf("New did not apply its provided option")
	}
	if accessed, added := c.recordAccess("foo"), c.recordAdd("foo"); accessed == nil || added != nil {
		t.Errorf("Expected default LRU policy")
	}
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
		if testCase.cache.Size() != testCase.expectedSize {
			t.Errorf("Expected size of %d, got %d", testCase.expectedSize, testCase.cache.Size())
		}
		actual := testCase.cache.Get(keys...)
		if len(actual) != len(testCase.expectedItems) {
			t.Errorf("Expected to get %d items, got %d", len(testCase.expectedItems), len(actual))
		} else {
			for i, expectedItem := range testCase.expectedItems {
				if actual[i] != expectedItem {
					t.Errorf("Expected Get to return %v in position %d, got %v", expectedItem, i, actual[i])
				}
			}
		}
	}
}
