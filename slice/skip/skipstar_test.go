/*
Copyright 2014 Workiva, LLC

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

package skip

/*
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStarInsert(t *testing.T) {
	ssl := NewStar(uint8(0))

	e1 := newMockEntry(7)
	e2 := newMockEntry(10)

	result := ssl.Insert(e1, e2)
	assert.Equal(t, Entries{nil, nil}, result)
	assert.Equal(t, Entries{e1}, ssl.Get(7))
	assert.Equal(t, Entries{e2}, ssl.Get(10))
	assert.Equal(t, Entries{e1, e2}, ssl.Get(7, 10))
	assert.Equal(t, Entries{e1, nil}, ssl.Get(7, 13))
	assert.Equal(t, Entries{e2, nil}, ssl.Get(10, 13))
	assert.Equal(t, uint64(2), ssl.Len())
}

func TestStarOverwrite(t *testing.T) {
	ssl := NewStar(uint8(0))
	e1 := newMockEntry(7)
	e2 := newMockEntry(7)

	result := ssl.Insert(e1)
	assert.Equal(t, Entries{nil}, result)
	assert.Equal(t, uint64(1), ssl.Len())

	result = ssl.Insert(e2)
	assert.Equal(t, Entries{e1}, result)
	assert.Equal(t, uint64(1), ssl.Len())
}

func TestStarDelete(t *testing.T) {
	ssl := NewStar(uint8(0))
	e1 := newMockEntry(5)
	e2 := newMockEntry(10)
	ssl.Insert(e1, e2)

	result := ssl.Delete(e1.Key(), e2.Key())
	assert.Equal(t, Entries{e1, e2}, result)
	assert.Equal(t, uint64(0), ssl.Len())
}

func TestStarIter(t *testing.T) {
	ssl := NewStar(uint8(0))

	iter := ssl.Iter(0)
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Value())

	e1 := newMockEntry(5)
	e2 := newMockEntry(10)
	ssl.Insert(e1, e2)

	iter = ssl.Iter(0)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())

	iter = ssl.Iter(5)
	assert.Equal(t, Entries{e1, e2}, iter.exhaust())

	iter = ssl.Iter(6)
	assert.Equal(t, Entries{e2}, iter.exhaust())

	iter = ssl.Iter(10)
	assert.Equal(t, Entries{e2}, iter.exhaust())

	iter = ssl.Iter(11)
	assert.Equal(t, Entries{}, iter.exhaust())
}

func BenchmarkStarInsert(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(entries[i%numItems])
	}
}

func BenchmarkStarGet(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Get(entries[i%numItems].Key())
	}
}

func BenchmarkStarDelete(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Delete(entries[i%numItems].Key())
	}
}

func BenchmarkIterStar(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := generateMockEntries(numItems)
	sl.Insert(entries...)

	var iter Iterator
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for iter = sl.Iter(0); iter.Next(); {
		}
	}
}

func BenchmarkStarPrepend(b *testing.B) {
	numItems := b.N
	sl := NewStar(uint64(0))

	entries := make(Entries, 0, numItems)
	for i := b.N; i < b.N+numItems; i++ {
		entries = append(entries, newMockEntry(uint64(i)))
	}

	sl.Insert(entries...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sl.Insert(newMockEntry(uint64(i)))
	}
}
*/
