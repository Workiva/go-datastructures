package fastinteger

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateKeys(num int) []uint64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keys := make([]uint64, 0, num)
	for i := 0; i < num; i++ {
		key := uint64(r.Int63())
		keys = append(keys, key)
	}

	return keys
}

func TestRoundUp(t *testing.T) {
	result := roundUp(21)
	assert.Equal(t, uint64(32), result)

	result = roundUp(uint64(1<<31) - 234)
	assert.Equal(t, uint64(1<<31), result)

	result = roundUp(uint64(1<<63) - 324)
	assert.Equal(t, uint64(1<<63), result)
}

func TestInsert(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)

	assert.True(t, hm.Exists(5))
	value, ok := hm.Get(5)
	assert.Equal(t, uint64(5), value)
	assert.True(t, ok)
}

func TestInsertOverwrite(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)
	hm.Set(5, 10)

	assert.True(t, hm.Exists(5))
	value, ok := hm.Get(5)
	assert.Equal(t, uint64(10), value)
	assert.True(t, ok)
}

func TestGet(t *testing.T) {
	hm := New(10)

	value, ok := hm.Get(5)
	assert.False(t, ok)
	assert.Equal(t, uint64(0), value)
}

func TestMultipleInserts(t *testing.T) {
	hm := New(10)

	hm.Set(5, 5)
	hm.Set(6, 6)

	assert.True(t, hm.Exists(6))
	value, ok := hm.Get(6)
	assert.True(t, ok)
	assert.Equal(t, uint64(6), value)
}

func TestRebuild(t *testing.T) {
	numItems := uint64(100)

	hm := New(10)

	for i := uint64(0); i < numItems; i++ {
		hm.Set(i, i)
	}

	for i := uint64(0); i < numItems; i++ {
		value, _ := hm.Get(i)
		assert.Equal(t, i, value)
	}
}

func BenchmarkInsert(b *testing.B) {
	numItems := uint64(1000)

	keys := generateKeys(int(numItems))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hm := New(numItems * 2) // so we don't rebuild
		for _, k := range keys {
			hm.Set(k, k)
		}
	}
}

func BenchmarkGoMapInsert(b *testing.B) {
	numItems := uint64(1000)

	keys := generateKeys(int(numItems))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hm := make(map[uint64]uint64, numItems*2) // so we don't rebuild
		for _, k := range keys {
			hm[k] = k
		}
	}
}

func BenchmarkExists(b *testing.B) {
	numItems := uint64(1000)

	keys := generateKeys(int(numItems))
	hm := New(numItems * 2) // so we don't rebuild
	for _, key := range keys {
		hm.Set(key, key)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			hm.Exists(key)
		}
	}
}

func BenchmarkGoMapExists(b *testing.B) {
	numItems := uint64(1000)

	keys := generateKeys(int(numItems))
	hm := make(map[uint64]uint64, numItems*2) // so we don't rebuild
	for _, key := range keys {
		hm[key] = key
	}

	b.ResetTimer()

	var ok bool
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			_, ok = hm[key] // or the compiler complains
		}
	}

	b.StopTimer()
	b.Logf(`OK: %+v`, ok)
}
