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

/*
Package bitarray or bitmap is useful when comparing large amounts of structured
data if the data can be represented as integers.  For instance, set
intersection of {1, 3, 5} and {3, 5, 7} represented as bitarrays can
be done in a single clock cycle (not counting the time it takes to convert)
the resultant array back into integers).  When Go implements a command
to get trailing zeroes, the reconversion back into integers should be much faster.
*/
package bitarray

// BitArray represents a structure that can be used to
// quickly check for existence when using a large number
// of items in a very memory efficient way.
type BitArray interface {
	// SetBit sets the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	SetBit(k uint64) error
	// GetBit gets the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	GetBit(k uint64) (bool, error)
	// GetSetBits gets the position of bits set in the array. Will
	// return as many set bits as can fit in the provided buffer
	// starting from the specified position in the array.
	GetSetBits(from uint64, buffer []uint64) []uint64
	// ClearBit clears the bit at the given position.  This
	// function returns an error if the position is out
	// of range.  A sparse bit array never returns an error.
	ClearBit(k uint64) error
	// Reset sets all values to zero.
	Reset()
	// Blocks returns an iterator to be used to iterate
	// over the bit array.
	Blocks() Iterator
	// Equals returns a bool indicating equality between the
	// two bit arrays.
	Equals(other BitArray) bool
	// Intersects returns a bool indicating if the other bit
	// array intersects with this bit array.
	Intersects(other BitArray) bool
	// Capacity returns either the given capacity of the bit array
	// in the case of a dense bit array or the highest possible
	// seen capacity of the sparse array.
	Capacity() uint64
	// Count returns the number of set bits in this array.
	Count() int
	// Or will bitwise or the two bitarrays and return a new bitarray
	// representing the result.
	Or(other BitArray) BitArray
	// And will bitwise and the two bitarrays and return a new bitarray
	// representing the result.
	And(other BitArray) BitArray
	// Nand will bitwise nand the two bitarrays and return a new bitarray
	// representing the result.
	Nand(other BitArray) BitArray
	// ToNums converts this bit array to the list of numbers contained
	// within it.
	ToNums() []uint64
	// IsEmpty checks to see if any values are set on the bitarray
	IsEmpty() bool
}

// Iterator defines methods used to iterate over a bit array.
type Iterator interface {
	// Next moves the pointer to the next block.  Returns
	// false when no blocks remain.
	Next() bool
	// Value returns the next block and its index
	Value() (uint64, block)
}
