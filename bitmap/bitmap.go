// Package bitmap contains bitmaps of length 32 and 64 for tracking bool
// values without the need for arrays or hashing.
package bitmap

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

// PopCount returns the ammount of bits set to 1 in the Bitmap32
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

// PopCount returns the ammount of bits set to 1 in the Bitmap64
func (b Bitmap64) PopCount() int {
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	b -= (b >> 1) & 0x5555555555555555
	b = (b>>2)&0x3333333333333333 + b&0x3333333333333333
	b += b >> 4
	b &= 0x0f0f0f0f0f0f0f0f
	b *= 0x0101010101010101
	return int(byte(b >> 56))
}
