package fastinteger

import (
	"encoding/binary"
	"hash/fnv"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	key := uint64(5)
	h := hash(key)

	assert.NotEqual(t, key, h)
}

func BenchmarkHash(b *testing.B) {
	numItems := 1000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keys := make([]uint64, 0, numItems)
	for i := 0; i < numItems; i++ {
		key := uint64(r.Int63())
		keys = append(keys, key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			hash(key)
		}
	}
}

func BenchmarkFnvHash(b *testing.B) {
	numItems := 1000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keys := make([]uint64, 0, numItems)
	for i := 0; i < numItems; i++ {
		key := uint64(r.Int63())
		keys = append(keys, key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, key := range keys {
			hasher := fnv.New64()
			binary.Write(hasher, binary.LittleEndian, key)
			hasher.Sum64()
		}
	}
}
