/*
Copyright 2015 Workiva, LLC

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

package ctrie

import (
	"hash"
	"hash/fnv"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCtrie(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)

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

	val, ok = ctrie.Remove([]byte("foo"))
	assert.True(ok)
	assert.Equal("qux", val)

	_, ok = ctrie.Remove([]byte("foo"))
	assert.False(ok)

	val, ok = ctrie.Remove([]byte("fooooo"))
	assert.True(ok)
	assert.Equal("baz", val)

	for i := 0; i < 100; i++ {
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}
}

type mockHash32 struct {
	hash.Hash32
}

func (m *mockHash32) Sum32() uint32 {
	return 0
}

func mockHashFactory() hash.Hash32 {
	return &mockHash32{fnv.New32a()}
}

func TestInsertLNode(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(mockHashFactory)

	for i := 0; i < 10; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	for i := 0; i < 10; i++ {
		val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}
	_, ok := ctrie.Lookup([]byte("11"))
	assert.False(ok)

	for i := 0; i < 10; i++ {
		val, ok := ctrie.Remove([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}
}

func TestInsertTNode(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)

	for i := 0; i < 100000; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	for i := 0; i < 50000; i++ {
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}

	for i := 0; i < 100000; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	for i := 0; i < 100000; i++ {
		val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}
}

func TestConcurrency(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 100000; i++ {
			ctrie.Insert([]byte(strconv.Itoa(i)), i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
			if ok {
				assert.Equal(i, val)
			}
		}
		wg.Done()
	}()

	for i := 0; i < 100000; i++ {
		time.Sleep(5)
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}

	wg.Wait()
}

func BenchmarkInsert(b *testing.B) {
	ctrie := New(nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctrie.Insert([]byte("foo"), 0)
	}
}

func BenchmarkLookup(b *testing.B) {
	numItems := 1000
	ctrie := New(nil)
	for i := 0; i < numItems; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	key := []byte(strconv.Itoa(numItems / 2))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctrie.Lookup(key)
	}
}

func BenchmarkRemove(b *testing.B) {
	numItems := 1000
	ctrie := New(nil)
	for i := 0; i < numItems; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	key := []byte(strconv.Itoa(numItems / 2))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctrie.Remove(key)
	}
}
