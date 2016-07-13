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

package bitarray

import "sort"

// uintSlice is an alias for a slice of ints.  Len, Swap, and Less
// are exported to fulfill an interface needed for the search
// function in the sort library.
type uintSlice []uint64

// Len returns the length of the slice.
func (u uintSlice) Len() int64 {
	return int64(len(u))
}

// Swap swaps values in this slice at the positions given.
func (u uintSlice) Swap(i, j int64) {
	u[i], u[j] = u[j], u[i]
}

// Less returns a bool indicating if the value at position i is
// less than position j.
func (u uintSlice) Less(i, j int64) bool {
	return u[i] < u[j]
}

func (u uintSlice) search(x uint64) int64 {
	return int64(sort.Search(len(u), func(i int) bool { return uint64(u[i]) >= x }))
}

func (u *uintSlice) insert(x uint64) (int64, bool) {
	i := u.search(x)

	if i == int64(len(*u)) {
		*u = append(*u, x)
		return i, true
	}

	if (*u)[i] == x {
		return i, false
	}

	*u = append(*u, 0)
	copy((*u)[i+1:], (*u)[i:])
	(*u)[i] = x
	return i, true
}

func (u *uintSlice) deleteAtIndex(i int64) {
	copy((*u)[i:], (*u)[i+1:])
	(*u)[len(*u)-1] = 0
	*u = (*u)[:len(*u)-1]
}

func (u uintSlice) get(x uint64) int64 {
	i := u.search(x)
	if i == int64(len(u)) {
		return -1
	}

	if u[i] == x {
		return i
	}

	return -1
}

type blocks []block

func (b *blocks) insert(index int64) {
	if index == int64(len(*b)) {
		*b = append(*b, block(0))
		return
	}

	*b = append(*b, block(0))
	copy((*b)[index+1:], (*b)[index:])
	(*b)[index] = block(0)
}

func (b *blocks) deleteAtIndex(i int64) {
	copy((*b)[i:], (*b)[i+1:])
	(*b)[len(*b)-1] = block(0)
	*b = (*b)[:len(*b)-1]
}

type sparseBitArray struct {
	blocks  blocks
	indices uintSlice
}

// SetBit sets the bit at the given position.
func (sba *sparseBitArray) SetBit(k uint64) error {
	index, position := getIndexAndRemainder(k)
	i, inserted := sba.indices.insert(index)
	if inserted {
		sba.blocks.insert(i)
	}
	sba.blocks[i] = sba.blocks[i].insert(position)
	return nil
}

// GetBit gets the bit at the given position.
func (sba *sparseBitArray) GetBit(k uint64) (bool, error) {
	index, position := getIndexAndRemainder(k)
	i := sba.indices.get(index)
	if i == -1 {
		return false, nil
	}

	return sba.blocks[i].get(position), nil
}

// ToNums converts this sparse bitarray to a list of numbers contained
// within it.
func (sba *sparseBitArray) ToNums() []uint64 {
	if len(sba.indices) == 0 {
		return nil
	}

	diff := uint64(len(sba.indices)) * s
	nums := make([]uint64, 0, diff/4)

	for i, offset := range sba.indices {
		sba.blocks[i].toNums(offset*s, &nums)
	}

	return nums
}

// ClearBit clears the bit at the given position.
func (sba *sparseBitArray) ClearBit(k uint64) error {
	index, position := getIndexAndRemainder(k)
	i := sba.indices.get(index)
	if i == -1 {
		return nil
	}

	sba.blocks[i] = sba.blocks[i].remove(position)
	if sba.blocks[i] == 0 {
		sba.blocks.deleteAtIndex(i)
		sba.indices.deleteAtIndex(i)
	}

	return nil
}

// Reset erases all values from this bitarray.
func (sba *sparseBitArray) Reset() {
	sba.blocks = sba.blocks[:0]
	sba.indices = sba.indices[:0]
}

// Blocks returns an iterator to iterator of this bitarray's blocks.
func (sba *sparseBitArray) Blocks() Iterator {
	return newCompressedBitArrayIterator(sba)
}

// Capacity returns the value of the highest possible *seen* value
// in this sparse bitarray.
func (sba *sparseBitArray) Capacity() uint64 {
	if len(sba.indices) == 0 {
		return 0
	}

	return (sba.indices[len(sba.indices)-1] + 1) * s
}

// Equals returns a bool indicating if the provided bit array
// equals this bitarray.
func (sba *sparseBitArray) Equals(other BitArray) bool {
	if other.Capacity() == 0 && sba.Capacity() > 0 {
		return false
	}

	var selfIndex uint64
	for iter := other.Blocks(); iter.Next(); {
		otherIndex, otherBlock := iter.Value()
		if len(sba.indices) == 0 {
			if otherBlock > 0 {
				return false
			}

			continue
		}

		if selfIndex >= uint64(len(sba.indices)) {
			return false
		}

		if otherIndex < sba.indices[selfIndex] {
			if otherBlock > 0 {
				return false
			}
			continue
		}

		if otherIndex > sba.indices[selfIndex] {
			return false
		}

		if !sba.blocks[selfIndex].equals(otherBlock) {
			return false
		}

		selfIndex++
	}

	return true
}

// Or will perform a bitwise or operation with the provided bitarray and
// return a new result bitarray.
func (sba *sparseBitArray) Or(other BitArray) BitArray {
	if ba, ok := other.(*sparseBitArray); ok {
		return orSparseWithSparseBitArray(sba, ba)
	}

	return orSparseWithDenseBitArray(sba, other.(*bitArray))
}

// And will perform a bitwise and operation with the provided bitarray and
// return a new result bitarray.
func (sba *sparseBitArray) And(other BitArray) BitArray {
	if ba, ok := other.(*sparseBitArray); ok {
		return andSparseWithSparseBitArray(sba, ba)
	}

	return andSparseWithDenseBitArray(sba, other.(*bitArray))
}

func (sba *sparseBitArray) copy() *sparseBitArray {
	blocks := make(blocks, len(sba.blocks))
	copy(blocks, sba.blocks)
	indices := make(uintSlice, len(sba.indices))
	copy(indices, sba.indices)
	return &sparseBitArray{
		blocks:  blocks,
		indices: indices,
	}
}

// Intersects returns a bool indicating if the provided bit array
// intersects with this bitarray.
func (sba *sparseBitArray) Intersects(other BitArray) bool {
	if other.Capacity() == 0 {
		return true
	}

	var selfIndex int64
	for iter := other.Blocks(); iter.Next(); {
		otherI, otherBlock := iter.Value()
		if len(sba.indices) == 0 {
			if otherBlock > 0 {
				return false
			}
			continue
		}
		// here we grab where the block should live in ourselves
		i := uintSlice(sba.indices[selfIndex:]).search(otherI)
		// this is a block we don't have, doesn't intersect
		if i == int64(len(sba.indices)) {
			return false
		}

		if sba.indices[i] != otherI {
			return false
		}

		if !sba.blocks[i].intersects(otherBlock) {
			return false
		}

		selfIndex = i
	}

	return true
}

func (sba *sparseBitArray) IntersectsBetween(other BitArray, start, stop uint64) bool {
	return true
}

func newSparseBitArray() *sparseBitArray {
	return &sparseBitArray{}
}

// NewSparseBitArray will create a bit array that consumes a great
// deal less memory at the expense of longer sets and gets.
func NewSparseBitArray() BitArray {
	return newSparseBitArray()
}
