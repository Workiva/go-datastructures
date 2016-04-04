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

// Bitmap32 tracks 32 bool values within a uint32
type Bitmap32 uint32

// SetBit returns a Bitmap32 with the bit at the given position set to 1
func (b Bitmap32) SetBit(pos uint) Bitmap32 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap32 with the bit at the given position set to 0
func (b Bitmap32) ClearBit(pos uint) Bitmap32 {
	return b & ^(1 << pos)
}

// HasBit returns true if the bit at the given position in the Bitmap32 is 1
func (b Bitmap32) HasBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap32
func (b Bitmap32) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x55555555
	b = (b>>2)&0x33333333 + b&0x33333333
	b += b >> 4
	b &= 0x0f0f0f0f
	b *= 0x01010101
	return int(byte(b >> 24))
}

// Bitmap64 tracks 64 bool values within a uint64
type Bitmap64 uint64

// SetBit returns a Bitmap64 with the bit at the given position set to 1
func (b Bitmap64) SetBit(pos uint) Bitmap64 {
	return b | (1 << pos)
}

// ClearBit returns a Bitmap64 with the bit at the given position set to 0
func (b Bitmap64) ClearBit(pos uint) Bitmap64 {
	return b & ^(1 << pos)
}

// HasBit returns true if the bit at the given position in the Bitmap64 is 1
func (b Bitmap64) HasBit(pos uint) bool {
	return (b & (1 << pos)) != 0
}

// PopCount returns the amount of bits set to 1 in the Bitmap64
func (b Bitmap64) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x5555555555555555
	b = (b>>2)&0x3333333333333333 + b&0x3333333333333333
	b += b >> 4
	b &= 0x0f0f0f0f0f0f0f0f
	b *= 0x0101010101010101
	return int(byte(b >> 56))
}
