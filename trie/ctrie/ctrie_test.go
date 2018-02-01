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

	for i := 0; i < 10000; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	for i := 0; i < 5000; i++ {
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}

	for i := 0; i < 10000; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	for i := 0; i < 10000; i++ {
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
		for i := 0; i < 10000; i++ {
			ctrie.Insert([]byte(strconv.Itoa(i)), i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
			if ok {
				assert.Equal(i, val)
			}
		}
		wg.Done()
	}()

	for i := 0; i < 10000; i++ {
		time.Sleep(5)
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}

	wg.Wait()
}

func TestConcurrency2(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		for i := 0; i < 10000; i++ {
			ctrie.Insert([]byte(strconv.Itoa(i)), i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
			if ok {
				assert.Equal(i, val)
			}
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			ctrie.Snapshot()
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			ctrie.ReadOnlySnapshot()
		}
		wg.Done()
	}()

	wg.Wait()
	assert.Equal(uint(10000), ctrie.Size())
}

func TestSnapshot(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)
	for i := 0; i < 100; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}

	snapshot := ctrie.Snapshot()

	// Ensure snapshot contains expected keys.
	for i := 0; i < 100; i++ {
		val, ok := snapshot.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}

	for i := 0; i < 100; i++ {
		ctrie.Remove([]byte(strconv.Itoa(i)))
	}

	// Ensure snapshot was unaffected by removals.
	for i := 0; i < 100; i++ {
		val, ok := snapshot.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}

	ctrie = New(nil)
	for i := 0; i < 100; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	snapshot = ctrie.Snapshot()

	// Ensure snapshot is mutable.
	for i := 0; i < 100; i++ {
		snapshot.Remove([]byte(strconv.Itoa(i)))
	}
	snapshot.Insert([]byte("bat"), "man")

	for i := 0; i < 100; i++ {
		_, ok := snapshot.Lookup([]byte(strconv.Itoa(i)))
		assert.False(ok)
	}
	val, ok := snapshot.Lookup([]byte("bat"))
	assert.True(ok)
	assert.Equal("man", val)

	// Ensure original Ctrie was unaffected.
	for i := 0; i < 100; i++ {
		val, ok := ctrie.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}
	_, ok = ctrie.Lookup([]byte("bat"))
	assert.False(ok)

	snapshot = ctrie.ReadOnlySnapshot()
	for i := 0; i < 100; i++ {
		val, ok := snapshot.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}

	// Ensure read-only snapshots panic on writes.
	defer func() {
		assert.NotNil(recover())
	}()
	snapshot.Remove([]byte("blah"))

	// Ensure snapshots-of-snapshots work as expected.
	snapshot2 := snapshot.Snapshot()
	for i := 0; i < 100; i++ {
		val, ok := snapshot2.Lookup([]byte(strconv.Itoa(i)))
		assert.True(ok)
		assert.Equal(i, val)
	}
	snapshot2.Remove([]byte("0"))
	_, ok = snapshot2.Lookup([]byte("0"))
	assert.False(ok)
	val, ok = snapshot.Lookup([]byte("0"))
	assert.True(ok)
	assert.Equal(0, val)
}

func TestIterator(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)
	for i := 0; i < 10; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	expected := map[string]int{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"9": 9,
	}

	count := 0
	for entry := range ctrie.Iterator(nil) {
		exp, ok := expected[string(entry.Key)]
		if assert.True(ok) {
			assert.Equal(exp, entry.Value)
		}
		count++
	}
	assert.Equal(len(expected), count)

	// Closing cancel channel should close iterator channel.
	cancel := make(chan struct{})
	iter := ctrie.Iterator(cancel)
	entry := <-iter
	exp, ok := expected[string(entry.Key)]
	if assert.True(ok) {
		assert.Equal(exp, entry.Value)
	}
	close(cancel)
	// Drain anything already put on the channel. Since select chooses a
	// pseudo-random case, we must attempt to drain for every item.
	for _ = range expected {
		<-iter
	}
	_, ok = <-iter
	assert.False(ok)
}

// TestIteratorCoversTNodes reproduces the scenario of a bug where tNodes weren't being traversed.
func TestIteratorCoversTNodes(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(mockHashFactory)
	// Add a pair of keys that collide (because we're using the mock hash).
	ctrie.Insert([]byte("a"), true)
	ctrie.Insert([]byte("b"), true)
	// Delete one key, leaving exactly one sNode in the cNode.  This will
	// trigger creation of a tNode.
	ctrie.Remove([]byte("b"))
	seenKeys := map[string]bool{}
	for entry := range ctrie.Iterator(nil) {
		seenKeys[string(entry.Key)] = true
	}
	assert.Contains(seenKeys, "a", "Iterator did not return 'a'.")
	assert.Len(seenKeys, 1)
}

func TestSize(t *testing.T) {
	ctrie := New(nil)
	for i := 0; i < 10; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	assert.Equal(t, uint(10), ctrie.Size())
}

func TestClear(t *testing.T) {
	assert := assert.New(t)
	ctrie := New(nil)
	for i := 0; i < 10; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	assert.Equal(uint(10), ctrie.Size())
	snapshot := ctrie.Snapshot()

	ctrie.Clear()

	assert.Equal(uint(0), ctrie.Size())
	assert.Equal(uint(10), snapshot.Size())
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

func BenchmarkSnapshot(b *testing.B) {
	numItems := 1000
	ctrie := New(nil)
	for i := 0; i < numItems; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctrie.Snapshot()
	}
}

func BenchmarkReadOnlySnapshot(b *testing.B) {
	numItems := 1000
	ctrie := New(nil)
	for i := 0; i < numItems; i++ {
		ctrie.Insert([]byte(strconv.Itoa(i)), i)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctrie.ReadOnlySnapshot()
	}
}
