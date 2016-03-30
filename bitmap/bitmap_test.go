package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopCount32(t *testing.T) {
	b := []uint32{
		uint32(0x55555555), // 0x55555555 = 01010101 01010101 01010101 01010101
		uint32(0x33333333), // 0x33333333 = 00110011 00110011 00110011 00110011
		uint32(0x0F0F0F0F), // 0x0F0F0F0F = 00001111 00001111 00001111 00001111
		uint32(0x00FF00FF), // 0x00FF00FF = 00000000 11111111 00000000 11111111
		uint32(0x0000FFFF), // 0x0000FFFF = 00000000 00000000 11111111 11111111
	}
	for _, x := range b {
		assert.Equal(t, 16, Bitmap32(x).PopCount())
	}
}

func TestPopCount64(t *testing.T) {
	b := []uint64{
		uint64(0x5555555555555555),
		uint64(0x3333333333333333),
		uint64(0x0F0F0F0F0F0F0F0F),
		uint64(0x00FF00FF00FF00FF),
		uint64(0x0000FFFF0000FFFF),
	}
	for _, x := range b {
		assert.Equal(t, 32, Bitmap64(x).PopCount())
	}
}

func TestSetBit(t *testing.T) {
	m := Bitmap32(0)
	assert.Equal(t, Bitmap32(0x4), m.SetBit(2))
}

func TestClearBit(t *testing.T) {
	m := Bitmap32(0x4)
	assert.Equal(t, Bitmap32(0), m.ClearBit(2))
}

func TestHasBit(t *testing.T) {
	m := Bitmap32(0x55555555)
	assert.Equal(t, true, m.HasBit(2))
}

func BenchmarkPopCount32(b *testing.B) {
	m := Bitmap32(0x33333333)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		m.PopCount()
	}
}

func BenchmarkPopCount64(b *testing.B) {
	m := Bitmap64(0x3333333333333333)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		m.PopCount()
	}
}
