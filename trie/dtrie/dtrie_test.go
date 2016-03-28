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
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopCount(t *testing.T) {
	b := []uint32{
		uint32(0x55555555), // 0x55555555 = 01010101 01010101 01010101 01010101
		uint32(0x33333333), // 0x33333333 = 00110011 00110011 00110011 00110011
		uint32(0x0F0F0F0F), // 0x0F0F0F0F = 00001111 00001111 00001111 00001111
		uint32(0x00FF00FF), // 0x00FF00FF = 00000000 11111111 00000000 11111111
		uint32(0x0000FFFF), // 0x0000FFFF = 00000000 00000000 11111111 11111111
	}
	for _, x := range b {
		assert.Equal(t, 16, popCount(x))
	}
}

func TestDefaultHasher(t *testing.T) {
	assert.Equal(t,
		defaultHasher(map[int]string{11234: "foo"}),
		defaultHasher(map[int]string{11234: "foo"}))
	assert.NotEqual(t, defaultHasher("foo"), defaultHasher("bar"))
}

type testEntry struct {
	hash  uint32
	key   int
	value int
}

func (e *testEntry) KeyHash() uint32 {
	return e.hash
}

func (e *testEntry) Key() interface{} {
	return e.key
}

func (e *testEntry) Value() interface{} {
	return e.value
}

func (e *testEntry) String() string {
	return fmt.Sprint(e.value)
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
		n = insert(n, &testEntry{hashfunc(i), i, i})
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
		n = insert(n, &testEntry{hashfunc(i), i, -i})
	}
}

func TestIterate(t *testing.T) {
	n := insertTest(t, defaultHasher, 10000)
	echan := iterate(n, nil)
	var c int64
	for _ = range echan {
		c++
	}
	assert.Equal(t, int64(10000), c)
	// test with stop chan
	c = 0
	stop := make(chan struct{})
	echan = iterate(n, stop)
	go func() {
		for range echan {
			atomic.AddInt64(&c, 1)
		}
	}()
	for atomic.LoadInt64(&c) < 100 {
	}
	close(stop)
	assert.True(t, c > 99 && c < 1000)
	// test with collisions
	n = insertTest(t, collisionHash, 1000)
	c = 0
	echan = iterate(n, nil)
	for _ = range echan {
		c++
	}
	assert.Equal(t, int64(1000), c)
}

func TestSize(t *testing.T) {
	n := insertTest(t, defaultHasher, 10000)
	d := &Dtrie{n, defaultHasher}
	assert.Equal(t, 10000, d.Size())
}

func BenchmarkInsert(b *testing.B) {
	b.ReportAllocs()
	n := emptyNode(0, 32)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		n = insert(n, &testEntry{defaultHasher(i), i, i})
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
		n = insert(n, &testEntry{defaultHasher(i), i, -i})
	}
}
