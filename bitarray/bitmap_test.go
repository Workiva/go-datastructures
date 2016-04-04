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

package bitarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitmap32_PopCount(t *testing.T) {
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

func TestBitmap64_PopCount(t *testing.T) {
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

func TestBitmap32_SetBit(t *testing.T) {
	m := Bitmap32(0)
	assert.Equal(t, Bitmap32(0x4), m.SetBit(2))
}

func TestBitmap32_ClearBit(t *testing.T) {
	m := Bitmap32(0x4)
	assert.Equal(t, Bitmap32(0), m.ClearBit(2))
}

func TestBitmap32_zGetBit(t *testing.T) {
	m := Bitmap32(0x55555555)
	assert.Equal(t, true, m.GetBit(2))
}

func TestBitmap64_SetBit(t *testing.T) {
	m := Bitmap64(0)
	assert.Equal(t, Bitmap64(0x4), m.SetBit(2))
}

func TestBitmap64_ClearBit(t *testing.T) {
	m := Bitmap64(0x4)
	assert.Equal(t, Bitmap64(0), m.ClearBit(2))
}

func TestBitmap64_GetBit(t *testing.T) {
	m := Bitmap64(0x55555555)
	assert.Equal(t, true, m.GetBit(2))
}

func BenchmarkBitmap32_PopCount(b *testing.B) {
	m := Bitmap32(0x33333333)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		m.PopCount()
	}
}

func BenchmarkBitmap64_PopCount(b *testing.B) {
	m := Bitmap64(0x3333333333333333)
	b.ResetTimer()
	for i := b.N; i > 0; i-- {
		m.PopCount()
	}
}
