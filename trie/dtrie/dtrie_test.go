/*
Copyright (c) 2016, Theodore Butler
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package dtrie

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHasher(t *testing.T) {
	assert.Equal(t,
		defaultHasher(map[int]string{11234: "foo"}),
		defaultHasher(map[int]string{11234: "foo"}))
	assert.NotEqual(t, defaultHasher("foo"), defaultHasher("bar"))
}

func collisionHash(key interface{}) uint32 {
	return uint32(0xffffffff) // for testing collisions
}

func TestInsert(t *testing.T) {
	insertTest(t, defaultHasher, 10000)
	insertTest(t, collisionHash, 1000)
}

func insertTest(t *testing.T, hashfunc func(interface{}) uint32, count int) *node {
	n := emptyNode(0, 32)
	for i := 0; i < count; i++ {
		n = insert(n, &entry{hashfunc(i), i, i})
	}
	return n
}

func TestGet(t *testing.T) {
	getTest(t, defaultHasher, 10000)
	getTest(t, collisionHash, 1000)
}

func getTest(t *testing.T, hashfunc func(interface{}) uint32, count int) {
	n := insertTest(t, hashfunc, count)
	for i := 0; i < count; i++ {
		x := get(n, hashfunc(i), i)
		assert.Equal(t, i, x.Value())
	}
}

func TestRemove(t *testing.T) {
	removeTest(t, defaultHasher, 10000)
	removeTest(t, collisionHash, 1000)
}

func removeTest(t *testing.T, hashfunc func(interface{}) uint32, count int) {
	n := insertTest(t, hashfunc, count)
	for i := 0; i < count; i++ {
		n = remove(n, hashfunc(i), i)
	}
	for _, e := range n.entries {
		if e != nil {
			t.Fatal("final node is not empty")
		}
	}
}

func TestUpdate(t *testing.T) {
	updateTest(t, defaultHasher, 10000)
	updateTest(t, collisionHash, 1000)
}

func updateTest(t *testing.T, hashfunc func(interface{}) uint32, count int) {
	n := insertTest(t, hashfunc, count)
	for i := 0; i < count; i++ {
		n = insert(n, &entry{hashfunc(i), i, -i})
	}
}

func TestIterate(t *testing.T) {
	n := insertTest(t, defaultHasher, 10000)
	echan := iterate(n, nil)
	c := 0
	for _ = range echan {
		c++
	}
	assert.Equal(t, 10000, c)
	// test with stop chan
	c = 0
	stop := make(chan struct{})
	echan = iterate(n, stop)
	for _ = range echan {
		c++
		if c == 100 {
			close(stop)
			break
		}
	}
	assert.True(t, c == 100)
	// test with collisions
	n = insertTest(t, collisionHash, 1000)
	c = 0
	echan = iterate(n, nil)
	for _ = range echan {
		c++
	}
	assert.Equal(t, 1000, c)
}

func TestSize(t *testing.T) {
	n := insertTest(t, defaultHasher, 10000)
	d := &Dtrie{n, defaultHasher}
	assert.Equal(t, 10000, d.Size())
}

func TestDtriePersistentVersions(t *testing.T) {
	tests := []struct {
		name   string
		keys   []int
		key    int
		remove bool
		hasher func(interface{}) uint32
	}{
		{name: "insert into empty trie", key: 1},
		{name: "insert root entry", keys: []int{0, 1}, key: 2},
		{name: "replace root entry", keys: []int{0, 1}, key: 0},
		{name: "insert nested entry", keys: []int{0, 32}, key: 64},
		{name: "replace nested entry", keys: []int{0, 32}, key: 32},
		{name: "create collision node", keys: []int{0}, key: 1, hasher: collisionHash},
		{name: "append collision entry", keys: []int{0, 1, 2, 3}, key: 4, hasher: collisionHash},
		{name: "remove root entry", keys: []int{0, 1}, key: 0, remove: true},
		{name: "remove nested entry", keys: []int{0, 32, 64}, key: 32, remove: true},
		{name: "remove collision entry", keys: []int{0, 1, 2, 3}, key: 1, remove: true, hasher: collisionHash},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := New(tt.hasher)
			wantOriginal := make(map[int]int)
			wantUpdated := make(map[int]int)
			for _, key := range tt.keys {
				original = original.Insert(key, key)
				wantOriginal[key] = key
				wantUpdated[key] = key
			}

			var updated *Dtrie
			if tt.remove {
				updated = original.Remove(tt.key)
				delete(wantUpdated, tt.key)
			} else {
				updated = original.Insert(tt.key, -1)
				wantUpdated[tt.key] = -1
			}

			assertDtrieValues(t, original, wantOriginal)
			assertDtrieValues(t, updated, wantUpdated)
		})
	}
}

func TestDtriePersistentBranches(t *testing.T) {
	original := New(nil).Insert(0, 0).Insert(1, 1)
	left := original.Insert(0, 10)
	right := original.Insert(1, 20)
	left = left.Remove(1)

	assertDtrieValues(t, original, map[int]int{0: 0, 1: 1})
	assertDtrieValues(t, left, map[int]int{0: 10})
	assertDtrieValues(t, right, map[int]int{0: 0, 1: 20})
}

func TestDtrieConcurrentVersions(t *testing.T) {
	original := New(nil).Insert(0, 0).Insert(1, 1)
	const count = 16
	versions := make([]*Dtrie, count)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(count)
	for i := range versions {
		go func(i int) {
			defer wg.Done()
			<-start
			versions[i] = original.Insert(0, i).Remove(1)
		}(i)
	}
	close(start)
	wg.Wait()

	assertDtrieValues(t, original, map[int]int{0: 0, 1: 1})
	for i, version := range versions {
		assertDtrieValues(t, version, map[int]int{0: i})
	}
}

func assertDtrieValues(t *testing.T, d *Dtrie, want map[int]int) {
	t.Helper()
	got := make(map[int]int)
	for e := range d.Iterator(nil) {
		got[e.Key().(int)] = e.Value().(int)
	}
	assert.Equal(t, want, got)
	assert.Equal(t, len(want), d.Size())
	for key, value := range want {
		assert.Equal(t, value, d.Get(key))
	}
}

func BenchmarkInsert(b *testing.B) {
	b.ReportAllocs()
	n := emptyNode(0, 32)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		n = insert(n, &entry{defaultHasher(i), i, i})
	}
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	n := insertTest(nil, defaultHasher, b.N)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		get(n, defaultHasher(i), i)
	}
}

func BenchmarkRemove(b *testing.B) {
	b.ReportAllocs()
	n := insertTest(nil, defaultHasher, b.N)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		n = remove(n, defaultHasher(i), i)
	}
}

func BenchmarkUpdate(b *testing.B) {
	b.ReportAllocs()
	n := insertTest(nil, defaultHasher, b.N)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		n = insert(n, &entry{defaultHasher(i), i, -i})
	}
}
